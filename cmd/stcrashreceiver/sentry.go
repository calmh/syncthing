// Copyright (C) 2019 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/maruel/panicparse/stack"
)

func sendReport(path, dsn string) error {
	pkt, err := loadReport(path)
	if err != nil {
		return err
	}

	cli, err := raven.New(dsn)
	if err != nil {
		return err
	}

	// The client sets release and such on the packet before sending, in the
	// misguided idea that it knows this better than than the packet we give
	// it. So we copy the values from the packet to the client first...
	cli.SetRelease(pkt.Release)
	cli.SetEnvironment(pkt.Environment)

	_, errC := cli.Capture(pkt, nil)
	return <-errC
}

func loadReport(path string) (*raven.Packet, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseReport(path, bs)
}

func parseReport(path string, report []byte) (*raven.Packet, error) {
	parts := bytes.SplitN(report, []byte("\n"), 2)
	if len(parts) != 2 {
		return nil, errors.New("no first line")
	}

	version, err := parseVersion(string(parts[0]))
	if err != nil {
		return nil, err
	}
	report = parts[1]

	foundPanic := false
	var subjectLine []byte
	for {
		parts = bytes.SplitN(report, []byte("\n"), 2)
		if len(parts) != 2 {
			return nil, errors.New("no panic line found")
		}

		line := parts[0]
		report = parts[1]

		if foundPanic {
			// The previous line was our "Panic at ..." header. We are now
			// at the beginning of the real panic trace and this is our
			// subject line.
			subjectLine = line
			break
		} else if bytes.HasPrefix(line, []byte("Panic at")) {
			foundPanic = true
		}
	}

	r := bytes.NewReader(report)
	ctx, err := stack.ParseDump(r, ioutil.Discard, false)
	if err != nil {
		return nil, err
	}

	var trace raven.Stacktrace
	for _, gr := range ctx.Goroutines {
		if gr.First {
			trace.Frames = make([]*raven.StacktraceFrame, len(gr.Stack.Calls))
			for i, sc := range gr.Stack.Calls {
				trace.Frames[len(trace.Frames)-1-i] = &raven.StacktraceFrame{
					Function: sc.Func.Name(),
					Module:   sc.Func.PkgName(),
					Filename: sc.SrcPath,
					Lineno:   sc.Line,
				}
			}
			break
		}
	}

	pkt := &raven.Packet{
		Message:  string(subjectLine),
		Platform: "go",
		Release:  version.tag,
		Tags: raven.Tags{
			raven.Tag{Key: "version", Value: version.version},
			raven.Tag{Key: "tag", Value: version.tag},
			raven.Tag{Key: "commit", Value: version.commit},
			raven.Tag{Key: "codename", Value: version.codename},
			raven.Tag{Key: "runtime", Value: version.runtime},
			raven.Tag{Key: "goos", Value: version.goos},
			raven.Tag{Key: "goarch", Value: version.goarch},
			raven.Tag{Key: "builder", Value: version.builder},
		},
		Extra: raven.Extra{
			"path": path,
		},
		Interfaces: []raven.Interface{&trace},
	}

	return pkt, nil
}

// syncthing v1.1.4-rc.1+30-g6aaae618-dirty-crashrep "Erbium Earthworm" (go1.12.5 darwin-amd64) jb@kvin.kastelo.net 2019-05-23 16:08:14 UTC
var longVersionRE = regexp.MustCompile(`syncthing\s+(v[^\s]+)\s+"([^"]+)"\s\(([^\s]+)\s+([^-]+)-([^)]+)\)\s+([^\s]+)`)

type version struct {
	version  string // "v1.1.4-rc.1+30-g6aaae618-dirty-crashrep"
	tag      string // "v1.1.4-rc.1"
	commit   string // "6aaae618", blank when absent
	codename string // "Erbium Earthworm"
	runtime  string // "go1.12.5"
	goos     string // "darwin"
	goarch   string // "amd64"
	builder  string // "jb@kvin.kastelo.net"
}

func parseVersion(line string) (version, error) {
	m := longVersionRE.FindStringSubmatch(line)
	if len(m) == 0 {
		return version{}, errors.New("unintelligeble version string")
	}

	v := version{
		version:  m[1],
		codename: m[2],
		runtime:  m[3],
		goos:     m[4],
		goarch:   m[5],
		builder:  m[6],
	}
	parts := strings.Split(v.version, "+")
	v.tag = parts[0]
	if len(parts) > 1 {
		fields := strings.Split(parts[1], "-")
		if len(fields) >= 2 && strings.HasPrefix(fields[1], "g") {
			v.commit = fields[1][1:]
		}
	}

	return v, nil
}
