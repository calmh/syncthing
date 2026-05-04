// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package model

import (
	"slices"
	"testing"
	"time"

	"github.com/syncthing/syncthing/lib/protocol"
)

func TestDeviceActivity(t *testing.T) {
	n0 := Availability{protocol.DeviceID([32]byte{1, 2, 3, 4}), false}
	n1 := Availability{protocol.DeviceID([32]byte{5, 6, 7, 8}), true}
	n2 := Availability{protocol.DeviceID([32]byte{9, 10, 11, 12}), false}
	devices := []Availability{n0, n1, n2}

	t.Run("basic", func(t *testing.T) {
		na := newDeviceActivity()

		// making blocks take the assumed rate means we take device rate out
		// of the equation for this test, basing it only on the number of
		// outstanding blocks
		blockDur := time.Second * protocol.MinBlockSize / assumedRate

		if lb := na.leastBusy(devices); lb != 0 {
			t.Errorf("Least busy device should be n0 (%v) not %v", n0, lb)
		}
		if lb := na.leastBusy(devices); lb != 0 {
			t.Errorf("Least busy device should still be n0 (%v) not %v", n0, lb)
		}

		lb := na.leastBusy(devices)
		na.using(devices[lb].ID, protocol.MinBlockSize)
		if lb := na.leastBusy(devices); lb != 1 {
			t.Errorf("Least busy device should be n1 (%v) not %v", n1, lb)
		}
		lb = na.leastBusy(devices)
		na.using(devices[lb].ID, protocol.MinBlockSize)
		if lb := na.leastBusy(devices); lb != 2 {
			t.Errorf("Least busy device should be n2 (%v) not %v", n2, lb)
		}

		lb = na.leastBusy(devices)
		na.using(devices[lb].ID, protocol.MinBlockSize)
		if lb := na.leastBusy(devices); lb != 0 {
			t.Errorf("Least busy device should be n0 (%v) not %v", n0, lb)
		}

		na.done(n1.ID, protocol.MinBlockSize, blockDur)
		if lb := na.leastBusy(devices); lb != 1 {
			t.Errorf("Least busy device should be n1 (%v) not %v", n1, lb)
		}

		na.done(n2.ID, protocol.MinBlockSize, blockDur)
		if lb := na.leastBusy(devices); lb != 1 {
			t.Errorf("Least busy device should still be n1 (%v) not %v", n1, lb)
		}

		na.done(n0.ID, protocol.MinBlockSize, blockDur)
		if lb := na.leastBusy(devices); lb != 0 {
			t.Errorf("Least busy device should be n0 (%v) not %v", n0, lb)
		}
	})

	t.Run("rateBased", func(t *testing.T) {
		na := newDeviceActivity()

		// n0 has proven to be quick, averaging ten blocks per second
		na.using(n0.ID, protocol.MinBlockSize)
		na.done(n0.ID, protocol.MinBlockSize, time.Second/10)

		// n1 is a bit slower, averaging two blocks per second
		na.using(n1.ID, protocol.MinBlockSize)
		na.done(n1.ID, protocol.MinBlockSize, time.Second/2)

		// n2 is a yet slower, averaging one block per second
		na.using(n2.ID, protocol.MinBlockSize)
		na.done(n2.ID, protocol.MinBlockSize, time.Second/1)

		// Request one hundred blocks, and observe the distribution
		count := make([]int, 3)
		for range 100 {
			idx := na.leastBusy(devices)
			count[idx]++
			na.using(devices[idx].ID, protocol.MinBlockSize)
		}

		// n0 should have been assigned 10/13 ~= 78% of the blocks
		// n1 should have been assigned 2/13 ~= 15% of the blocks
		// n2 should have been assigned 1/13 ~= 7% of the blocks
		exp := []int{78, 15, 7}

		if !slices.Equal(count, exp) {
			t.Error("Unexpected results", count)
		}
	})
}
