// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"slices"
	"testing"

	"google.golang.org/protobuf/proto"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/protocol"
)

func TestPrepare(t *testing.T) {
	id := protocol.NewDeviceID([]byte("prepare"))
	other := protocol.NewDeviceID([]byte("prepareother"))
	untrusted := protocol.NewDeviceID([]byte("prepareuntrusted"))
	ghost := protocol.NewDeviceID([]byte("prepareghost"))

	cfg := Configuration{syncthingv2.Configuration_builder{
		Version: new(int32(52)),
		Folders: []*syncthingv2.FolderConfiguration{
			// Sorted second; shares with an unknown device and an
			// untrusted one without a password, has a duplicate device,
			// a negative rescan interval, and disables the watcher with
			// a zero delay.
			syncthingv2.FolderConfiguration_builder{
				Id:   new("zeta"),
				Path: new("/srv/sync/zeta"),
				Devices: []*syncthingv2.FolderDeviceConfiguration{
					syncthingv2.FolderDeviceConfiguration_builder{DeviceId: new(ghost.String())}.Build(),
					syncthingv2.FolderDeviceConfiguration_builder{DeviceId: new(untrusted.String())}.Build(),
					syncthingv2.FolderDeviceConfiguration_builder{DeviceId: new(other.String())}.Build(),
					syncthingv2.FolderDeviceConfiguration_builder{DeviceId: new(other.String())}.Build(),
				},
				RescanIntervalS: new(int32(-1)),
				FsWatcherDelayS: new(0.0),
				MaxConflicts:    new(int32(-1)),
			}.Build(),
			// Sorted first; receive-encrypted, so permissions are
			// ignored regardless of setting.
			syncthingv2.FolderConfiguration_builder{
				Id:           new("alpha"),
				Path:         new("/srv/sync/alpha"),
				Type:         new(syncthingv2.FolderType_FOLDER_TYPE_RECEIVE_ENCRYPTED),
				Devices:      []*syncthingv2.FolderDeviceConfiguration{syncthingv2.FolderDeviceConfiguration_builder{DeviceId: new(id.String())}.Build()},
				IgnorePerms:  new(false),
				MaxConflicts: new(int32(0)), // explicit zero stays
			}.Build(),
		},
		Devices: []*syncthingv2.DeviceConfiguration{
			syncthingv2.DeviceConfiguration_builder{DeviceId: new(other.String())}.Build(),
			// Duplicate of the above.
			syncthingv2.DeviceConfiguration_builder{DeviceId: new(other.String()), Name: new("dup")}.Build(),
			// Untrusted, yet an introducer and auto-accepting.
			syncthingv2.DeviceConfiguration_builder{
				DeviceId:          new(untrusted.String()),
				Untrusted:         new(true),
				Introducer:        new(true),
				AutoAcceptFolders: new(true),
				Addresses:         []string{"dynamic", "tcp://192.0.2.1:22000", "tcp://192.0.2.1:22000"},
				IgnoredFolders:    []*syncthingv2.ObservedFolder{syncthingv2.ObservedFolder_builder{Id: new("alpha")}.Build()},
			}.Build(),
			// Empty device ID, removed.
			syncthingv2.DeviceConfiguration_builder{}.Build(),
		},
		RemoteIgnoredDevices: []*syncthingv2.ObservedDevice{
			// Present in the device list, removed.
			syncthingv2.ObservedDevice_builder{DeviceId: new(other.String())}.Build(),
			// Not present, kept.
			syncthingv2.ObservedDevice_builder{DeviceId: new(ghost.String())}.Build(),
		},
		Options: syncthingv2.OptionsConfiguration_builder{
			ListenAddresses:       []string{"tcp://0.0.0.0:22000", "tcp://0.0.0.0:22000"},
			ReconnectionIntervalS: new(int32(1)),
			UrAccepted:            new(int32(2)),
		}.Build(),
	}.Build()}

	if err := Prepare(&cfg, id); err != nil {
		t.Fatalf("prepare: %v", err)
	}

	// The local device was added; devices are deduplicated and sorted.
	devices := cfg.GetDevices()
	if len(devices) != 3 {
		t.Fatalf("devices: got %d, want 3", len(devices))
	}
	if !slices.IsSortedFunc(devices, func(a, b DeviceConfiguration) int {
		return a.GetDeviceId().Compare(b.GetDeviceId())
	}) {
		t.Error("devices should be sorted by ID")
	}
	for _, want := range []protocol.DeviceID{untrusted, other, id} {
		if _, _, ok := cfg.Device(want); !ok {
			t.Errorf("device %s missing", want)
		}
	}

	// The untrusted device lost its introducer and auto-accept flags.
	// Device addresses are not deduplicated, and the ignored folder is
	// kept: the folder shared with the device in trusted mode was
	// unshared, so it is no longer shared.
	untrustedDev, _, ok := cfg.Device(untrusted)
	if !ok {
		t.Fatal("untrusted device missing")
	}
	if untrustedDev.GetIntroducer() || untrustedDev.GetAutoAcceptFolders() {
		t.Error("untrusted device should not be introducer or auto-accept")
	}
	if addrs := untrustedDev.GetAddresses(); len(addrs) != 3 {
		t.Errorf("addresses: got %v, want no entries removed", addrs)
	}
	if folders := untrustedDev.GetIgnoredFolders(); len(folders) != 1 || folders[0].GetId() != "alpha" {
		t.Errorf("ignored folders: got %v, want [alpha]", folders)
	}

	// The remote ignored devices list only has devices not in the config.
	if ignored := cfg.GetRemoteIgnoredDevices(); len(ignored) != 1 || !ignored[0].GetDeviceId().Equals(ghost) {
		t.Errorf("remote ignored devices: got %v, want [%s]", ignored, ghost)
	}

	// Folders are sorted by ID, the local device is sharing each folder,
	// and unknown or untrusted-without-password devices are unshared.
	// The zeta folder remains shared with the trusted "other" device.
	folders := cfg.GetFolders()
	if len(folders) != 2 || folders[0].GetId() != "alpha" || folders[1].GetId() != "zeta" {
		t.Fatalf("folders: got %v, want [alpha zeta]", folders)
	}
	if devs := folders[0].GetDevices(); len(devs) != 1 || !devs[0].GetDeviceId().Equals(id) {
		t.Errorf("folder alpha devices: got %v, want only local device", devs)
	}
	if devs := folders[1].GetDevices(); len(devs) != 2 || !devs[0].GetDeviceId().Equals(other) || !devs[1].GetDeviceId().Equals(id) {
		t.Errorf("folder zeta devices: got %v, want [other local]", devs)
	}

	// The receive-encrypted folder ignores permissions.
	if !folders[0].GetIgnorePerms() {
		t.Error("receive-encrypted folder should ignore permissions")
	}
	// Explicit zero values are preserved.
	if !folders[0].HasMaxConflicts() || folders[0].GetMaxConflicts() != 0 {
		t.Error("explicit maxConflicts 0 should be preserved")
	}
	// The negative rescan interval became zero (disabled).
	if !folders[1].HasRescanIntervalS() || folders[1].GetRescanIntervalS() != 0 {
		t.Error("negative rescan interval should become zero")
	}
	// The zero watcher delay disabled the watcher.
	if folders[1].GetFsWatcherEnabled() {
		t.Error("zero watcher delay should disable the watcher")
	}
	// Marker name and max concurrent writes are unset; the defaults come
	// from the schema via the getters.
	if folders[1].HasMarkerName() {
		t.Error("marker name should be unset")
	}
	if got := folders[1].GetMarkerName(); got != DefaultMarkerName {
		t.Errorf("marker name: got %q, want default %q", got, DefaultMarkerName)
	}
	if folders[1].HasMaxConcurrentWrites() {
		t.Error("max concurrent writes should be unset")
	}
	if got := folders[1].GetMaxConcurrentWrites(); got != 16 {
		t.Errorf("max concurrent writes: got %d, want default 16", got)
	}

	// The listen addresses are deduplicated, and the reconnection
	// interval and usage reporting ID are fixed up.
	opts := cfg.GetOptions()
	if addrs := opts.GetListenAddresses(); len(addrs) != 1 || addrs[0] != "tcp://0.0.0.0:22000" {
		t.Errorf("listen addresses: got %v", addrs)
	}
	if v := opts.GetReconnectionIntervalS(); v != 5 {
		t.Errorf("reconnection interval: got %d, want 5", v)
	}
	if v := opts.GetUrUniqueId(); v == "" {
		t.Error("usage reporting unique ID should be generated")
	}

	// The GUI was created and has an API key.
	if !cfg.HasGui() || cfg.GetGui().GetApiKey() == "" {
		t.Error("GUI should exist with an API key")
	}

	// Preparing again is a no-op.
	before := cfg.Copy()
	if err := Prepare(&cfg, id); err != nil {
		t.Fatalf("second prepare: %v", err)
	}
	if !proto.Equal(before.Configuration, cfg.Configuration) {
		t.Error("second prepare changed the configuration")
	}
}

func TestPrepareErrors(t *testing.T) {
	id := protocol.NewDeviceID([]byte("prepareerr"))

	for _, tc := range []struct {
		name string
		cfg  *syncthingv2.Configuration
	}{
		{
			"empty folder ID",
			syncthingv2.Configuration_builder{
				Folders: []*syncthingv2.FolderConfiguration{syncthingv2.FolderConfiguration_builder{
					Path: new("/srv/sync/f1"),
				}.Build()},
			}.Build(),
		},
		{
			"empty folder path",
			syncthingv2.Configuration_builder{
				Folders: []*syncthingv2.FolderConfiguration{syncthingv2.FolderConfiguration_builder{
					Id: new("f1"),
				}.Build()},
			}.Build(),
		},
		{
			"duplicate folder ID",
			syncthingv2.Configuration_builder{
				Folders: []*syncthingv2.FolderConfiguration{
					syncthingv2.FolderConfiguration_builder{Id: new("f1"), Path: new("/1")}.Build(),
					syncthingv2.FolderConfiguration_builder{Id: new("f1"), Path: new("/2")}.Build(),
				},
			}.Build(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Configuration{tc.cfg}
			if err := Prepare(&cfg, id); err == nil {
				t.Errorf("prepare should have failed")
			}
		})
	}
}
