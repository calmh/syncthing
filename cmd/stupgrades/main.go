// Copyright (C) 2019 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/kong"
	"github.com/syncthing/syncthing/lib/upgrade"
)

type cli struct {
	Listen     string        `default:":8080" help:"Listen address"`
	URL        string        `short:"u" default:"https://api.github.com/repos/syncthing/syncthing/releases?per_page=25" help:"GitHub releases url"`
	Forward    []string      `short:"f" help:"Forwarded pages, format: /path->https://example/com/url"`
	CacheBytes int           `default:"10000000" help:"Cache size"`
	CacheTime  time.Duration `default:"15m" help:"Cache time"`
}

func main() {
	var params cli
	kong.Parse(&params)
	if err := server(&params); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func server(params *cli) error {
	http.Handle("/meta.json", &simpleCache{next: &githubReleases{url: params.URL}})

	for _, fwd := range params.Forward {
		path, url, ok := strings.Cut(fwd, "->")
		if !ok {
			return fmt.Errorf("invalid forward: %q", fwd)
		}
		http.Handle(path, &simpleCache{next: &proxy{url: url}})
	}

	return http.ListenAndServe(params.Listen, nil)
}

type githubReleases struct {
	url string
}

func (p *githubReleases) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("Fetching", p.url)
	rels := upgrade.FetchLatestReleases(p.url, "")
	if rels == nil {
		http.Error(w, "no releases", http.StatusInternalServerError)
		return
	}

	sort.Sort(upgrade.SortByRelease(rels))
	rels = filterForLatest(rels)

	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(rels)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Write(buf.Bytes())
}

type proxy struct {
	url string
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("Fetching", p.url)
	req, err := http.NewRequestWithContext(req.Context(), http.MethodGet, p.url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	w.Header().Set("Content-Type", ct)
	if resp.StatusCode == http.StatusOK {
		w.Header().Set("Cache-Control", "public, max-age=900")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
	}
	w.WriteHeader(resp.StatusCode)
	if strings.HasPrefix(ct, "application/json") {
		// Special JSON handling; clean it up a bit.
		var v interface{}
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(v)
	} else {
		_, _ = io.Copy(w, resp.Body)
	}
}

// filterForLatest returns the latest stable and prerelease only. If the
// stable version is newer (comes first in the list) there is no need to go
// looking for a prerelease at all.
func filterForLatest(rels []upgrade.Release) []upgrade.Release {
	var filtered []upgrade.Release
	var havePre bool
	for _, rel := range rels {
		if !rel.Prerelease {
			// We found a stable version, we're good now.
			filtered = append(filtered, rel)
			break
		}
		if rel.Prerelease && !havePre {
			// We remember the first prerelease we find.
			filtered = append(filtered, rel)
			havePre = true
		}
	}
	return filtered
}

type simpleCache struct {
	next http.Handler

	mut  sync.Mutex
	when time.Time
	resp *recordedResponse
}

type recordedResponse struct {
	status int
	header http.Header
	data   []byte
}

type responseRecorder struct {
	resp *recordedResponse
	http.ResponseWriter
}

func (r *responseRecorder) WriteHeader(status int) {
	r.resp.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	r.resp.data = append(r.resp.data, data...)
	return r.ResponseWriter.Write(data)
}

func (r *responseRecorder) Header() http.Header {
	return r.resp.header
}

func (s *simpleCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mut.Lock()
	defer s.mut.Unlock()

	if time.Since(s.when) < 15*time.Minute {
		for k, v := range s.resp.header {
			w.Header()[k] = v
		}
		w.WriteHeader(s.resp.status)
		_, _ = w.Write(s.resp.data)
		return
	}

	rec := &recordedResponse{status: http.StatusOK, header: make(http.Header)}
	w = &responseRecorder{ResponseWriter: w, resp: rec}
	s.next.ServeHTTP(w, r)
	s.resp = rec
	s.when = time.Now()
}
