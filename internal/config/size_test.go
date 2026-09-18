// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"testing"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/fs"
)

func TestSizeHelpers(t *testing.T) {
	tests := []struct {
		size       *syncthingv2.Size
		bytes      float64
		percentage bool
	}{
		{syncthingv2.Size_builder{Bytes: new(1024.0)}.Build(), 1024, false},
		{syncthingv2.Size_builder{Mib: new(2.0)}.Build(), 2 * (1 << 20), false},
		{syncthingv2.Size_builder{Gib: new(3.0)}.Build(), 3 * (1 << 30), false},
		{syncthingv2.Size_builder{Percent: new(10.0)}.Build(), 0, true},
		{syncthingv2.Size_builder{}.Build(), 0, false},
		{nil, 0, false},
	}
	for _, test := range tests {
		if got := SizeInBytes(test.size); got != test.bytes {
			t.Errorf("SizeInBytes(%v): got %v, want %v", test.size, got, test.bytes)
		}
		if got := SizeIsPercentage(test.size); got != test.percentage {
			t.Errorf("SizeIsPercentage(%v): got %v, want %v", test.size, got, test.percentage)
		}
	}
}

func TestCheckFreeSpace(t *testing.T) {
	usage := fs.Usage{Free: 500, Total: 1000}

	// A nil size or zero size means no requirement.
	if err := CheckFreeSpace(nil, usage); err != nil {
		t.Errorf("nil size: %v", err)
	}
	if err := CheckFreeSpace(syncthingv2.Size_builder{}.Build(), usage); err != nil {
		t.Errorf("empty size: %v", err)
	}
	if err := CheckFreeSpace(syncthingv2.Size_builder{Percent: new(0.0)}.Build(), usage); err != nil {
		t.Errorf("zero percent: %v", err)
	}
	if err := CheckFreeSpace(syncthingv2.Size_builder{Bytes: new(0.0)}.Build(), usage); err != nil {
		t.Errorf("zero bytes: %v", err)
	}

	// Fifty percent is free, so 10% is fine and 90% is not.
	if err := CheckFreeSpace(syncthingv2.Size_builder{Percent: new(10.0)}.Build(), usage); err != nil {
		t.Errorf("10%%: %v", err)
	}
	if err := CheckFreeSpace(syncthingv2.Size_builder{Percent: new(90.0)}.Build(), usage); err == nil {
		t.Error("90%: should fail")
	}

	// 500 bytes are free, so 100 is fine and 1000 is not.
	if err := CheckFreeSpace(syncthingv2.Size_builder{Bytes: new(100.0)}.Build(), usage); err != nil {
		t.Errorf("100B: %v", err)
	}
	if err := CheckFreeSpace(syncthingv2.Size_builder{Bytes: new(1000.0)}.Build(), usage); err == nil {
		t.Error("1000B: should fail")
	}
}
