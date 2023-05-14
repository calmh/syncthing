// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build !solaris && !windows
// +build !solaris,!windows

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/syncthing/syncthing/lib/protocol"
)

func recordPerfStats(ctx context.Context) {
	savePerfStats(ctx, fmt.Sprintf("perfstats-%d.csv", syscall.Getpid()))
}

func savePerfStats(ctx context.Context, file string) {
	fd, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	bw := bufio.NewWriter(fd)
	defer bw.Flush()

	t0 := time.Now()
	prevTime := t0.UnixNano()
	var rusage, prevRusage syscall.Rusage
	var memstats runtime.MemStats
	var prevIn, prevOut int64
	tc := time.NewTicker(250 * time.Millisecond)

	fmt.Fprintf(bw, "TIME\tREL\tCPU\tALLOC\tINUSE\tNETIN\tNETOUT\tBLOCKSIN\tBLOCKSOUT\n")
	report := func(t time.Time) {
		if err := syscall.Getrusage(syscall.RUSAGE_SELF, &rusage); err != nil {
			return
		}
		runtime.ReadMemStats(&memstats)

		relTime := t.Sub(t0).Seconds()
		curTime := t.UnixNano()
		timeDiff := curTime - prevTime
		usageNanos := rusage.Utime.Nano() - prevRusage.Utime.Nano() + rusage.Stime.Nano() - prevRusage.Stime.Nano()
		cpuUsagePercent := 100 * float64(usageNanos) / float64(timeDiff)
		blocksIn := rusage.Inblock - prevRusage.Inblock
		blocksOut := rusage.Oublock - prevRusage.Oublock
		prevTime = curTime
		prevRusage = rusage

		in, out := protocol.TotalInOut()
		netIn := in - prevIn
		netOut := out - prevOut
		prevIn, prevOut = in, out

		fmt.Fprintf(bw, "%.06f\t%.06f\t%f\t%d\t%d\t%d\t%d\t%d\t%d\n", float64(curTime)/float64(time.Second), relTime, cpuUsagePercent, memstats.Alloc, memstats.Sys-memstats.HeapReleased, netIn, netOut, blocksIn, blocksOut)
	}

	for {
		select {
		case t := <-tc.C:
			report(t)
		case <-ctx.Done():
			report(time.Now())
			bw.Write([]byte("---\n"))
			return
		}
	}
}
