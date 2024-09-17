// Copyright (C) 2018 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package servemetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/puzpuzpuz/xsync/v3"

	"github.com/syncthing/syncthing/lib/geoip"
	"github.com/syncthing/syncthing/lib/ur/contract"
)

type CLI struct {
	Listen          string `env:"UR_LISTEN" default:"0.0.0.0:8080"`
	GeoIPLicenseKey string `env:"UR_GEOIP_LICENSE_KEY"`
	GeoIPAccountID  int    `env:"UR_GEOIP_ACCOUNT_ID"`
}

var (
	compilerRe         = regexp.MustCompile(`\(([A-Za-z0-9()., -]+) \w+-\w+(?:| android| default)\) ([\w@.-]+)`)
	knownDistributions = []distributionMatch{
		// Maps well known builders to the official distribution method that
		// they represent

		{regexp.MustCompile(`\steamcity@build\.syncthing\.net`), "GitHub"},
		{regexp.MustCompile(`\sjenkins@build\.syncthing\.net`), "GitHub"},
		{regexp.MustCompile(`\sbuilder@github\.syncthing\.net`), "GitHub"},

		{regexp.MustCompile(`\sdeb@build\.syncthing\.net`), "APT"},
		{regexp.MustCompile(`\sdebian@github\.syncthing\.net`), "APT"},

		{regexp.MustCompile(`\sdocker@syncthing\.net`), "Docker Hub"},
		{regexp.MustCompile(`\sdocker@build.syncthing\.net`), "Docker Hub"},
		{regexp.MustCompile(`\sdocker@github.syncthing\.net`), "Docker Hub"},

		{regexp.MustCompile(`\sandroid-builder@github\.syncthing\.net`), "Google Play"},
		{regexp.MustCompile(`\sandroid-.*teamcity@build\.syncthing\.net`), "Google Play"},
		{regexp.MustCompile(`\sandroid-.*vagrant@basebox-stretch64`), "F-Droid"},
		{regexp.MustCompile(`\svagrant@bullseye`), "F-Droid"},
		{regexp.MustCompile(`\sbuilduser@(archlinux|svetlemodry)`), "Arch (3rd party)"},
		{regexp.MustCompile(`\ssyncthing@archlinux`), "Arch (3rd party)"},
		{regexp.MustCompile(`@debian`), "Debian (3rd party)"},
		{regexp.MustCompile(`@fedora`), "Fedora (3rd party)"},
		{regexp.MustCompile(`\sbrew@`), "Homebrew (3rd party)"},
		{regexp.MustCompile(`\sroot@buildkitsandbox`), "LinuxServer.io (3rd party)"},
		{regexp.MustCompile(`\sports@freebsd`), "FreeBSD (3rd party)"},
		{regexp.MustCompile(`.`), "Others"},
	}
)

type distributionMatch struct {
	matcher      *regexp.Regexp
	distribution string
}

func (cli *CLI) Run() error {
	// Listening

	listener, err := net.Listen("tcp", cli.Listen)
	if err != nil {
		log.Fatalln("listen:", err)
	}

	geo, err := geoip.NewGeoLite2CityProvider(context.Background(), cli.GeoIPAccountID, cli.GeoIPLicenseKey, os.TempDir())
	if err != nil {
		log.Fatalln("geoip:", err)
	}
	go geo.Serve(context.TODO())

	// server

	srv := &server{
		geo:     geo,
		reports: xsync.NewMapOf[string, *contract.Report](),
	}

	// New metrics endpoint

	ms := newMetricsSet(srv)
	if err := prometheus.Register(ms); err != nil {
		log.Fatalln("prometheus:", err)
	}

	http.HandleFunc("/newdata", srv.newDataHandler)
	http.Handle("/metrics", promhttp.Handler())

	httpSrv := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	return httpSrv.Serve(listener)
}

type server struct {
	geo     *geoip.Provider
	reports *xsync.MapOf[string, *contract.Report]
}

func (s *server) newDataHandler(w http.ResponseWriter, r *http.Request) {
	version := "fail"
	defer func() {
		// Version is "fail", "duplicate", "v2", "v3", ...
		metricReportsTotal.WithLabelValues(version).Inc()
	}()

	defer r.Body.Close()

	addr := r.Header.Get("X-Forwarded-For")
	if addr != "" {
		addr = strings.Split(addr, ", ")[0]
	} else {
		addr = r.RemoteAddr
	}

	if host, _, err := net.SplitHostPort(addr); err == nil {
		addr = host
	}

	if net.ParseIP(addr) == nil {
		addr = ""
	}

	var rep contract.Report
	rep.Date = time.Now().UTC().Format("20060102")
	rep.Address = addr

	lr := &io.LimitedReader{R: r.Body, N: 40 * 1024}
	bs, _ := io.ReadAll(lr)
	if err := json.Unmarshal(bs, &rep); err != nil {
		log.Println("decode:", err)
		http.Error(w, "JSON Decode Error", http.StatusInternalServerError)
		return
	}

	if err := rep.Validate(); err != nil {
		log.Println("validate:", err)
		http.Error(w, "Validation Error", http.StatusInternalServerError)
		return
	}

	if ip, err := net.LookupIP(rep.Address); err == nil && len(ip) > 0 {
		if city, err := s.geo.City(ip[0]); err == nil {
			rep.Country = city.Country.Names["en"]
			rep.City = city.City.Names["en"]
		}
	}
	if rep.Country == "" {
		rep.Country = "Unknown"
	}
	if rep.City == "" {
		rep.City = "Unknown"
	}

	rep.Version = transformVersion(rep.Version)
	rep.OS, rep.Arch, _ = strings.Cut(rep.Platform, "-")

	if m := compilerRe.FindStringSubmatch(rep.LongVersion); len(m) == 3 {
		rep.Compiler = m[1]
		rep.Builder = m[2]
	}
	for _, d := range knownDistributions {
		if d.matcher.MatchString(rep.LongVersion) {
			rep.Distribution = d.distribution
			break
		}
	}

	version = fmt.Sprintf("v%d", rep.URVersion)
	s.reports.Store(rep.UniqueID, &rep)
}

var (
	plusRe  = regexp.MustCompile(`(\+.*|\.dev\..*)$`)
	plusStr = "(+dev)"
)

// transformVersion returns a version number formatted correctly, with all
// development versions aggregated into one.
func transformVersion(v string) string {
	if v == "unknown-dev" {
		return v
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	v = plusRe.ReplaceAllString(v, " "+plusStr)

	return v
}
