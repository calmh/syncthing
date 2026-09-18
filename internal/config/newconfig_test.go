// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"testing"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/protocol"
)

func TestNewConfiguration(t *testing.T) {
	id := protocol.NewDeviceID([]byte("newcfg"))
	cfg := New(id)

	if v := cfg.GetVersion(); v != int32(CurrentVersion) {
		t.Errorf("version: got %d, want %d", v, CurrentVersion)
	}
	// The local device is present, and the GUI section has been created
	// with an API key.
	if _, _, ok := cfg.Device(id); !ok {
		t.Error("local device missing")
	}
	if cfg.GetGui().GetApiKey() == "" {
		t.Error("API key missing")
	}
	if ids := cfg.GetOptions().GetUnackedNotificationIds(); len(ids) != 1 || ids[0] != "authenticationUserAndPassword" {
		t.Errorf("unacked notification IDs: got %v", ids)
	}
	// The configuration is valid.
	if err := Validate(cfg); err != nil {
		t.Errorf("validate: %v", err)
	}
}

func TestConfigurationMapsAndSetters(t *testing.T) {
	id := protocol.NewDeviceID([]byte("maps"))
	other := protocol.NewDeviceID([]byte("mapsother"))
	f1 := FolderConfiguration{syncthingv2.FolderConfiguration_builder{
		Id:   new("f1"),
		Path: new("/srv/sync/f1"),
	}.Build()}
	f2 := FolderConfiguration{syncthingv2.FolderConfiguration_builder{
		Id:   new("f2"),
		Path: new("/srv/sync/f2"),
	}.Build()}
	cfg := Configuration{syncthingv2.Configuration_builder{
		Folders: []*syncthingv2.FolderConfiguration{f1.FolderConfiguration},
		Devices: []*syncthingv2.DeviceConfiguration{syncthingv2.DeviceConfiguration_builder{
			DeviceId: new(id.String()),
		}.Build()},
	}.Build()}

	// Maps are keyed by ID and device ID.
	if m := cfg.FolderMap(); len(m) != 1 || m["f1"].GetId() != "f1" {
		t.Errorf("folder map: got %v", m)
	}
	if m := cfg.DeviceMap(); len(m) != 1 || !m[id].GetDeviceId().Equals(id) {
		t.Errorf("device map: got %v", m)
	}

	// SetFolder appends new folders and replaces existing ones.
	cfg.SetFolder(f2)
	if len(cfg.GetFolders()) != 2 {
		t.Fatalf("folders: got %d, want 2", len(cfg.GetFolders()))
	}
	f2.SetLabel("Replaced")
	cfg.SetFolder(f2)
	if len(cfg.GetFolders()) != 2 {
		t.Fatalf("folders: got %d, want 2 after replace", len(cfg.GetFolders()))
	}
	if cfg.FolderMap()["f2"].GetLabel() != "Replaced" {
		t.Error("folder not replaced")
	}

	// SetDevice appends new devices and replaces existing ones.
	cfg.SetDevice(DeviceConfiguration{syncthingv2.DeviceConfiguration_builder{
		DeviceId: new(other.String()),
	}.Build()})
	if len(cfg.GetDevices()) != 2 {
		t.Fatalf("devices: got %d, want 2", len(cfg.GetDevices()))
	}
	cfg.SetDevice(DeviceConfiguration{syncthingv2.DeviceConfiguration_builder{
		DeviceId: new(other.String()),
		Name:     new("Replaced"),
	}.Build()})
	if len(cfg.GetDevices()) != 2 {
		t.Fatalf("devices: got %d, want 2 after replace", len(cfg.GetDevices()))
	}
	if cfg.DeviceMap()[other].GetName() != "Replaced" {
		t.Error("device not replaced")
	}
}

func TestProbeFreePorts(t *testing.T) {
	id := protocol.NewDeviceID([]byte("probeports"))
	cfg := New(id)

	if err := cfg.ProbeFreePorts(); err != nil {
		t.Fatalf("probe free ports: %v", err)
	}

	// The GUI address and listen addresses are set to concrete values
	// (the defaults, if available).
	if !cfg.GetGui().HasAddress() {
		t.Error("GUI address not set")
	}
	if addrs := cfg.GetOptions().GetListenAddresses(); len(addrs) == 0 {
		t.Error("listen addresses not set")
	}
}
