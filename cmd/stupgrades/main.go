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
	Listen  string   `default:":8080" help:"Listen address"`
	URL     string   `short:"u" default:"https://api.github.com/repos/syncthing/syncthing/releases?per_page=25" help:"GitHub releases url"`
	Forward []string `short:"f" help:"Forwarded pages, format: /path->https://example/com/url"`
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
	http.HandleFunc("/meta.json", (&githubReleases{url: params.URL}).serve)

	for _, fwd := range params.Forward {
		path, url, ok := strings.Cut(fwd, "->")
		if !ok {
			return fmt.Errorf("invalid forward: %q", fwd)
		}
		http.HandleFunc(path, (&proxy{url: url}).serve)
	}

	return http.ListenAndServe(params.Listen, nil)
}

type githubReleases struct {
	url string

	mut  sync.Mutex
	data []byte
	when time.Time
}

func (p *githubReleases) serve(w http.ResponseWriter, req *http.Request) {
	p.mut.Lock()
	defer p.mut.Unlock()

	if time.Since(p.when) > 5*time.Minute {
		rels := upgrade.FetchLatestReleases(p.url, "")
		if rels == nil {
			http.Error(w, "no releases", http.StatusInternalServerError)
			return
		}

		sort.Sort(upgrade.SortByRelease(rels))
		rels = filterForLatest(rels)

		buf := new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(rels)
		p.data = buf.Bytes()
		p.when = time.Now()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=900")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Write(p.data)
}

type proxy struct {
	url string
}

func (p *proxy) serve(w http.ResponseWriter, req *http.Request) {
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
