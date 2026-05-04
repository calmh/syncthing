// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package model

import (
	"sync"
	"time"

	"github.com/syncthing/syncthing/lib/protocol"
)

const (
	assumedRate = 1 << 20       // 1 MiB/s
	minimumRate = 128 << 10 / 8 // 128 KiB/s
)

// deviceActivity tracks the number of outstanding requests per device and can
// answer which device is least busy. It is safe for use from multiple
// goroutines.
type deviceActivity struct {
	act  map[protocol.DeviceID]int   // device ID -> outstanding bytes
	rate map[protocol.DeviceID][]int // device ID -> bytes/s
	mut  sync.Mutex
}

func newDeviceActivity() *deviceActivity {
	return &deviceActivity{
		act:  make(map[protocol.DeviceID]int),
		rate: make(map[protocol.DeviceID][]int),
	}
}

// Returns the index of the least busy device, or -1 if there are no
// available devices.
func (m *deviceActivity) leastBusy(availability []Availability) int {
	best := -1
	var shortestQueue float64

	m.mut.Lock()
	for i, a := range availability {
		if wt := m.waitTimeLocked(a.ID); shortestQueue == 0 || wt < shortestQueue {
			shortestQueue = wt
			best = i
		}
	}
	m.mut.Unlock()

	return best
}

func (m *deviceActivity) using(id protocol.DeviceID, bytes int) {
	m.mut.Lock()
	m.act[id] += bytes
	m.mut.Unlock()
}

func (m *deviceActivity) done(id protocol.DeviceID, bytes int, dur time.Duration) {
	m.mut.Lock()
	m.act[id] -= bytes
	if dur > 0 {
		rate := int(float64(bytes) / dur.Seconds())
		m.rate[id] = appendRate(m.rate[id], rate, 10)
	}
	m.mut.Unlock()
}

// waitTimeLocked returns how long we'll need to wait for the currently
// queued requests for the given device to resolve.
func (m *deviceActivity) waitTimeLocked(id protocol.DeviceID) float64 {
	rate := 0
	if len(m.rate[id]) > 0 {
		// rate is average observed device rate in bytes per second
		for _, r := range m.rate[id] {
			rate += r
		}
		rate /= len(m.rate[id])
		if rate < minimumRate {
			// Letting the rate get too low means we'll never schedule a
			// block on the device, which also means we'd never update the
			// rate. Set a lower boundary.
			rate = minimumRate
		}
	} else {
		// We have no rate stats yet; assume a rate so we have something to
		// calculate with, and so that we are likely to at least schedule a
		// block to this device to measure the rate.
		rate = assumedRate
	}

	// Calculate how long the current request queue is, in seconds. We add a
	// constant overhead factor so that the queue time remains dependent on
	// the observed rate even when there are no currently outstanding
	// requests, and so that we never return a zero wait time.
	return float64(m.act[id]+protocol.MinBlockSize) / float64(rate)
}

func appendRate(l []int, r int, maxL int) []int {
	if len(l) >= maxL {
		copy(l, l[1:])
		l[len(l)-1] = r
	} else {
		l = append(l, r)
	}
	return l
}
