// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"slices"
	"testing"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/protocol"
)

func optsCfg(opts *syncthingv2.OptionsConfiguration) OptionsConfiguration {
	if opts == nil {
		opts = syncthingv2.OptionsConfiguration_builder{}.Build()
	}
	return OptionsConfiguration{opts}
}

func TestEffectiveAddresses(t *testing.T) {
	// Unset address lists expand to the defaults.
	opts := optsCfg(nil)
	if got := opts.EffectiveListenAddresses(); !slices.Equal(got, DefaultListenAddresses) {
		t.Errorf("listen addresses: got %v, want %v", got, DefaultListenAddresses)
	}
	wantDiscovery := []string{
		"https://discovery-lookup.syncthing.net/v2/?noannounce",
		"https://discovery-announce-v4.syncthing.net/v2/?nolookup",
		"https://discovery-announce-v6.syncthing.net/v2/?nolookup",
	}
	if got := opts.EffectiveGlobalAnnounceServers(); !slices.Equal(got, wantDiscovery) {
		t.Errorf("global announce servers: got %v, want %v", got, wantDiscovery)
	}

	// The "default" keyword expands, explicit entries are kept, and the
	// result is deduplicated: the tcp address is already in the default
	// set, and the lookup server is in both the v4 and v6 sets.
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{
		ListenAddresses:       []string{"default", "tcp://0.0.0.0:22000"},
		GlobalAnnounceServers: []string{"default-v4", "default-v6", "https://discovery.example.com/"},
	}.Build())
	got := opts.EffectiveListenAddresses()
	if !slices.Equal(got, DefaultListenAddresses) {
		t.Errorf("listen addresses: got %v, want %v", got, DefaultListenAddresses)
	}
	got = opts.EffectiveGlobalAnnounceServers()
	if !slices.Equal(got, append(wantDiscovery, "https://discovery.example.com/")) {
		t.Errorf("global announce servers: got %v", got)
	}

	// Explicit STUN servers are used as-is.
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{
		StunServers: []string{"stun.example.com:3478"},
	}.Build())
	if got := opts.EffectiveStunServers(); !slices.Equal(got, []string{"stun.example.com:3478"}) {
		t.Errorf("stun servers: got %v", got)
	}
}

func TestOptionsComputedGetters(t *testing.T) {
	// MaxFolderConcurrency: unset means the number of CPUs, negative
	// means unlimited (zero), positive is used as-is.
	opts := optsCfg(nil)
	if got := opts.MaxFolderConcurrency(); got <= 0 {
		t.Errorf("max folder concurrency: got %d, want CPU count", got)
	}
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{MaxFolderConcurrency: new(int32(-1))}.Build())
	if got := opts.MaxFolderConcurrency(); got != 0 {
		t.Errorf("max folder concurrency: got %d, want 0", got)
	}
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{MaxFolderConcurrency: new(int32(7))}.Build())
	if got := opts.MaxFolderConcurrency(); got != 7 {
		t.Errorf("max folder concurrency: got %d, want 7", got)
	}

	// MaxConcurrentIncomingRequestKiB: unset means the default, negative
	// means disabled, too-small values are raised to the minimum.
	if got := optsCfg(nil).MaxConcurrentIncomingRequestKiB(); got != 256*1024 {
		t.Errorf("max concurrent incoming request: got %d, want 256 MiB", got)
	}
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{MaxConcurrentIncomingRequestKiB: new(int32(-1))}.Build())
	if got := opts.MaxConcurrentIncomingRequestKiB(); got != 0 {
		t.Errorf("max concurrent incoming request: got %d, want 0", got)
	}
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{MaxConcurrentIncomingRequestKiB: new(int32(1))}.Build())
	if got := opts.MaxConcurrentIncomingRequestKiB(); got <= 1 {
		t.Errorf("max concurrent incoming request: got %d, want raised to minimum", got)
	}

	// AutoUpgradeEnabled is driven by the interval.
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{AutoUpgradeIntervalH: new(int32(12))}.Build())
	if !opts.AutoUpgradeEnabled() {
		t.Error("auto upgrade should be enabled")
	}
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{AutoUpgradeIntervalH: new(int32(0))}.Build())
	if opts.AutoUpgradeEnabled() {
		t.Error("auto upgrade should be disabled")
	}

	// Feature flags.
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{FeatureFlags: []string{"caves"}}.Build())
	if !opts.FeatureFlag("caves") {
		t.Error("feature flag should be set")
	}
	if opts.FeatureFlag("other") {
		t.Error("feature flag should not be set")
	}

	// IsStunDisabled.
	opts = optsCfg(nil)
	if opts.IsStunDisabled() {
		t.Error("stun should not be disabled by default")
	}
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{StunKeepaliveMinS: new(int32(0))}.Build())
	if !opts.IsStunDisabled() {
		t.Error("stun should be disabled with zero keepalive")
	}

	// LowestConnectionLimit picks the lower of the two set limits.
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{
		ConnectionLimitEnough: new(int32(10)),
		ConnectionLimitMax:    new(int32(5)),
	}.Build())
	if got := opts.LowestConnectionLimit(); got != 5 {
		t.Errorf("lowest connection limit: got %d, want 5", got)
	}
	opts = optsCfg(syncthingv2.OptionsConfiguration_builder{
		ConnectionLimitEnough: new(int32(3)),
	}.Build())
	if got := opts.LowestConnectionLimit(); got != 3 {
		t.Errorf("lowest connection limit: got %d, want 3", got)
	}
}

func TestDeviceComputedGetters(t *testing.T) {
	id := protocol.NewDeviceID([]byte("computeddevice"))
	other := protocol.NewDeviceID([]byte("computedother"))
	device := DeviceConfiguration{syncthingv2.DeviceConfiguration_builder{
		DeviceId: new(id.String()),
		Name:     new("Alpha"),
		IgnoredFolders: []*syncthingv2.ObservedFolder{syncthingv2.ObservedFolder_builder{
			Id: new("ignored"),
		}.Build()},
	}.Build()}

	// NumConnections: zero means the default of three, negative means
	// one, positive is used as-is.
	if got := device.NumConnections(); got != 3 {
		t.Errorf("num connections: got %d, want 3", got)
	}
	device.SetNumConnections(-1)
	if got := device.NumConnections(); got != 1 {
		t.Errorf("num connections: got %d, want 1", got)
	}
	device.SetNumConnections(7)
	if got := device.NumConnections(); got != 7 {
		t.Errorf("num connections: got %d, want 7", got)
	}

	if !device.IgnoredFolder("ignored") {
		t.Error("folder should be ignored")
	}
	if device.IgnoredFolder("other") {
		t.Error("folder should not be ignored")
	}

	if got := device.Description(); got != "Alpha ("+id.Short().String()+")" {
		t.Errorf("description: got %q", got)
	}
	unnamed := DeviceConfiguration{syncthingv2.DeviceConfiguration_builder{
		DeviceId: new(other.String()),
	}.Build()}
	if got := unnamed.Description(); got != other.Short().String() {
		t.Errorf("description: got %q, want short ID", got)
	}
}
