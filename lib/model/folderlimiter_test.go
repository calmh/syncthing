// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package model

import (
	"context"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/syncthing/syncthing/lib/config"
)

func TestFolderRateLimiterSetLimits(t *testing.T) {
	l := newFolderRateLimiter()

	l.setLimits([]config.FolderConfiguration{
		{ID: "limited", MaxSendKbps: 100, MaxRecvKbps: 200},
		{ID: "unlimited"},
	})

	// A configured limit is converted from KiB/s to bytes/s.
	if got := l.limiter(true, "limited").Limit(); got != rate.Limit(100*1024) {
		t.Errorf("send limit = %v, want %v", got, rate.Limit(100*1024))
	}
	if got := l.limiter(false, "limited").Limit(); got != rate.Limit(200*1024) {
		t.Errorf("recv limit = %v, want %v", got, rate.Limit(200*1024))
	}

	// A folder without limits is unlimited in both directions.
	if got := l.limiter(true, "unlimited").Limit(); got != rate.Inf {
		t.Errorf("unlimited send limit = %v, want Inf", got)
	}

	// Updating an existing folder changes its limit without replacing the
	// limiter, and dropping a folder forgets it entirely.
	l.setLimits([]config.FolderConfiguration{
		{ID: "limited", MaxSendKbps: 50},
	})
	if got := l.limiter(true, "limited").Limit(); got != rate.Limit(50*1024) {
		t.Errorf("updated send limit = %v, want %v", got, rate.Limit(50*1024))
	}
	// MaxRecvKbps went back to 0, so the recv side must now be unlimited.
	if got := l.limiter(false, "limited").Limit(); got != rate.Inf {
		t.Errorf("updated recv limit = %v, want Inf", got)
	}
	if l.limiter(true, "unlimited") != nil {
		t.Error("removed folder limiter should be forgotten")
	}
}

func TestFolderRateLimiterUnlimitedIsImmediate(t *testing.T) {
	l := newFolderRateLimiter()
	l.setLimits([]config.FolderConfiguration{{ID: "folder"}})

	// Waiting on an unlimited (or unknown) folder must not block, even for
	// amounts far larger than the burst size.
	ctx := context.Background()
	if err := l.waitSend(ctx, "folder", 100*folderLimiterBurst); err != nil {
		t.Fatal(err)
	}
	if err := l.waitRecv(ctx, "unknown", 100*folderLimiterBurst); err != nil {
		t.Fatal(err)
	}
}

func TestFolderRateLimiterThrottles(t *testing.T) {
	l := newFolderRateLimiter()
	// 512 KiB/s, so consuming two bursts (after the initial full bucket) must
	// take a measurable amount of time.
	l.setLimits([]config.FolderConfiguration{{ID: "folder", MaxSendKbps: 512}})

	ctx := context.Background()
	// Drain the initial burst allowance.
	if err := l.waitSend(ctx, "folder", folderLimiterBurst); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	if err := l.waitSend(ctx, "folder", folderLimiterBurst); err != nil {
		t.Fatal(err)
	}
	// A full burst at 512 KiB/s is 1 second; allow generous slack for slow CI.
	if elapsed := time.Since(start); elapsed < 500*time.Millisecond {
		t.Errorf("expected throttling to take at least 500ms, took %v", elapsed)
	}
}

func TestFolderRateLimiterContextCancellation(t *testing.T) {
	l := newFolderRateLimiter()
	l.setLimits([]config.FolderConfiguration{{ID: "folder", MaxSendKbps: 1}})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// With a cancelled context and a rate that can't immediately satisfy the
	// request, WaitN returns the context error rather than blocking.
	if err := l.waitSend(ctx, "folder", folderLimiterBurst); err == nil {
		t.Error("expected error from cancelled context")
	}
}
