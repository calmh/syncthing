// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"testing"

	"google.golang.org/protobuf/proto"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
)

func TestMaterializeDefaults(t *testing.T) {
	cfg := Configuration{syncthingv2.Configuration_builder{
		Version: new(int32(52)),
		Folders: []*syncthingv2.FolderConfiguration{syncthingv2.FolderConfiguration_builder{
			Id: new("f1"),
			// rescan_interval_s unset (schema default 3600),
			// min_disk_free unset (getter default one percent),
			// max_conflicts explicitly zero.
			Path:         new("/srv/sync/f1"),
			MaxConflicts: new(int32(0)),
		}.Build()},
		Devices: []*syncthingv2.DeviceConfiguration{
			// Addresses unset (getter default "dynamic").
			syncthingv2.DeviceConfiguration_builder{DeviceId: new("AIR6LPZ-7ENK4MY-4VVODLE-VDGP6LY-Q3GSA2S-7SJMDTM-DLJ4BGV-K3ZP4QU")}.Build(),
			// Addresses set, kept as they are.
			syncthingv2.DeviceConfiguration_builder{
				DeviceId:  new("ZK6FOFT-TXHAKOT-DNJRW3B-7EDSH2F-C4QIYYI-ZAX2UXW-I3HZLGY-YGLAZQU"),
				Addresses: []string{"tcp://192.0.2.1:22000"},
			}.Build(),
		},
		Options: syncthingv2.OptionsConfiguration_builder{
			// Listen addresses unset (getter default "default"), a
			// global announce server set and kept.
			GlobalAnnounceServers: []string{"https://discovery.example.com/"},
		}.Build(),
		Defaults: syncthingv2.Defaults_builder{
			Folder: syncthingv2.FolderConfiguration_builder{
				Path: new("/srv/sync/default"),
			}.Build(),
		}.Build(),
	}.Build()}

	got := MaterializeDefaults(cfg)

	folder := got.GetFolders()[0]
	// Schema defaults are filled in.
	if !folder.HasRescanIntervalS() || folder.GetRescanIntervalS() != 3600 {
		t.Errorf("rescanIntervalS: got %d (set %v), want default 3600", folder.GetRescanIntervalS(), folder.HasRescanIntervalS())
	}
	// Explicitly set zero values are kept.
	if !folder.HasMaxConflicts() || folder.GetMaxConflicts() != 0 {
		t.Errorf("maxConflicts: got %d (set %v), want explicit 0", folder.GetMaxConflicts(), folder.HasMaxConflicts())
	}
	// Getter defaults are filled in.
	if size := folder.GetMinDiskFree(); size == nil || !size.HasPercent() || size.GetPercent() != 1 {
		t.Errorf("minDiskFree: got %v, want one percent", folder.GetMinDiskFree())
	}

	// The device without addresses gets the default; the one with
	// addresses keeps them.
	if addrs := got.GetDevices()[0].GetAddresses(); len(addrs) != 1 || addrs[0] != "dynamic" {
		t.Errorf("device 0 addresses: got %v, want [dynamic]", addrs)
	}
	if addrs := got.GetDevices()[1].GetAddresses(); len(addrs) != 1 || addrs[0] != "tcp://192.0.2.1:22000" {
		t.Errorf("device 1 addresses: got %v, want [tcp://192.0.2.1:22000]", addrs)
	}

	opts := got.GetOptions()
	// Getter defaults for the repeated address fields.
	if addrs := opts.GetListenAddresses(); len(addrs) != 1 || addrs[0] != "default" {
		t.Errorf("listenAddresses: got %v, want [default]", addrs)
	}
	if addrs := opts.GetStunServers(); len(addrs) != 1 || addrs[0] != "default" {
		t.Errorf("stunServers: got %v, want [default]", addrs)
	}
	if size := opts.GetMinHomeDiskFree(); size == nil || !size.HasPercent() || size.GetPercent() != 1 {
		t.Errorf("minHomeDiskFree: got %v, want one percent", size)
	}
	// Set values are kept.
	if servers := opts.GetGlobalAnnounceServers(); len(servers) != 1 || servers[0] != "https://discovery.example.com/" {
		t.Errorf("globalAnnounceServers: got %v", servers)
	}

	// Getter defaults apply inside the defaults template as well.
	if size := got.GetDefaults().GetFolder().GetMinDiskFree(); size == nil || !size.HasPercent() || size.GetPercent() != 1 {
		t.Errorf("defaults.folder.minDiskFree: got %v, want one percent", size)
	}

	// The original configuration is untouched.
	if orig := cfg.GetFolders()[0]; orig.HasRescanIntervalS() || orig.GetMinDiskFree() != nil {
		t.Error("materialization modified the original configuration")
	}
	if orig := cfg.GetOptions(); len(orig.GetListenAddresses()) != 0 {
		t.Error("materialization modified the original configuration")
	}

	// Materializing again is a no-op.
	again := MaterializeDefaults(got)
	if !proto.Equal(got.Configuration, again.Configuration) {
		t.Error("second materialization changed the configuration")
	}
}
