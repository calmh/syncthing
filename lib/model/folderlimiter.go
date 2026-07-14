// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package model

import (
	"context"
	"sync"

	"golang.org/x/time/rate"

	"github.com/syncthing/syncthing/lib/config"
)

// folderLimiterBurst caps how many tokens we consume in a single WaitN call.
// A WaitN larger than the burst never completes, so larger transfers are
// split into chunks of this size. It matches the burst used by the connection
// level (global and per-device) limiters.
const folderLimiterBurst = 4 * 128 << 10

// folderRateLimiter enforces the per-folder send and receive rate limits.
// These apply in addition to, and independently of, the connection level
// global and per-device limits: a transfer must satisfy all applicable limits.
// The zero value is not usable; use newFolderRateLimiter. It is safe for
// concurrent use.
type folderRateLimiter struct {
	mut  sync.Mutex
	send map[string]*rate.Limiter // folder ID -> outgoing limiter
	recv map[string]*rate.Limiter // folder ID -> incoming limiter
}

func newFolderRateLimiter() *folderRateLimiter {
	return &folderRateLimiter{
		send: make(map[string]*rate.Limiter),
		recv: make(map[string]*rate.Limiter),
	}
}

// setLimits updates the limiters to match the given folder configurations,
// creating limiters for new folders, updating changed ones and forgetting
// folders that are no longer present.
func (l *folderRateLimiter) setLimits(folders []config.FolderConfiguration) {
	l.mut.Lock()
	defer l.mut.Unlock()

	seen := make(map[string]struct{}, len(folders))
	for _, folder := range folders {
		seen[folder.ID] = struct{}{}
		setFolderLimit(l.send, folder.ID, folder.MaxSendKbps)
		setFolderLimit(l.recv, folder.ID, folder.MaxRecvKbps)
	}
	for id := range l.send {
		if _, ok := seen[id]; !ok {
			delete(l.send, id)
			delete(l.recv, id)
		}
	}
}

// waitSend blocks until n bytes may be sent for the given folder, or the
// context is cancelled.
func (l *folderRateLimiter) waitSend(ctx context.Context, folder string, n int) error {
	return waitFolder(ctx, l.send[folder], n)
}

// waitRecv blocks until n bytes may be received for the given folder, or the
// context is cancelled.
func (l *folderRateLimiter) waitRecv(ctx context.Context, folder string, n int) error {
	return waitFolder(ctx, l.recv[folder], n)
}

// setFolderLimit creates or updates the limiter for a folder. The rate is
// given in KiB/s; zero or less means unlimited.
func setFolderLimit(limiters map[string]*rate.Limiter, folder string, kbps int) {
	limit := rate.Inf
	if kbps > 0 {
		limit = rate.Limit(kbps) * 1024
	}
	if existing, ok := limiters[folder]; ok {
		existing.SetLimit(limit)
		return
	}
	limiters[folder] = rate.NewLimiter(limit, folderLimiterBurst)
}

func waitFolder(ctx context.Context, limiter *rate.Limiter, n int) error {
	if limiter == nil || limiter.Limit() == rate.Inf {
		return nil
	}
	// No single WaitN call may exceed the burst size or it never returns, so
	// we consume the tokens in burst sized chunks.
	for n > 0 {
		chunk := min(n, folderLimiterBurst)
		if err := limiter.WaitN(ctx, chunk); err != nil {
			return err
		}
		n -= chunk
	}
	return nil
}
