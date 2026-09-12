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

// TestWrapperDeviceIDs tests that the wrapper types expose device IDs as
// the native type throughout the configuration tree.
func TestWrapperDeviceIDs(t *testing.T) {
	id := protocol.NewDeviceID([]byte("wrapper"))
	other := protocol.NewDeviceID([]byte("otherdevice"))

	cfg := Configuration{syncthingv2.Configuration_builder{
		Folders: []*syncthingv2.FolderConfiguration{syncthingv2.FolderConfiguration_builder{
			Id: new("f1"),
			Devices: []*syncthingv2.FolderDeviceConfiguration{syncthingv2.FolderDeviceConfiguration_builder{
				DeviceId:     new(id.String()),
				IntroducedBy: new(other.String()),
			}.Build()},
		}.Build()},
		Devices: []*syncthingv2.DeviceConfiguration{syncthingv2.DeviceConfiguration_builder{
			DeviceId:     new(id.String()),
			IntroducedBy: new(other.String()),
		}.Build()},
		RemoteIgnoredDevices: []*syncthingv2.ObservedDevice{syncthingv2.ObservedDevice_builder{
			DeviceId: new(id.String()),
		}.Build()},
		Defaults: syncthingv2.Defaults_builder{}.Build(),
	}.Build()}

	tests := []struct {
		name string
		got  protocol.DeviceID
		want protocol.DeviceID
	}{
		{"folder device", cfg.GetFolders()[0].GetDevices()[0].GetDeviceId(), id},
		{"folder device introducer", cfg.GetFolders()[0].GetDevices()[0].GetIntroducedBy(), other},
		{"device", cfg.GetDevices()[0].GetDeviceId(), id},
		{"device introducer", cfg.GetDevices()[0].GetIntroducedBy(), other},
		{"ignored device", cfg.GetRemoteIgnoredDevices()[0].GetDeviceId(), id},
		// The device defaults have no device ID set.
		{"defaults device", cfg.GetDefaults().GetDevice().GetDeviceId(), protocol.EmptyDeviceID},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("%s: got %s, want %s", test.name, test.got, test.want)
		}
	}

	// Setters store the canonical string form in the underlying message.
	device := cfg.GetFolders()[0].GetDevices()[0]
	device.SetDeviceId(other)
	device.SetIntroducedBy(protocol.EmptyDeviceID)
	pb := cfg.Configuration.GetFolders()[0].GetDevices()[0]
	if got := pb.GetDeviceId(); got != other.String() {
		t.Errorf("set device ID: got %q, want %q", got, other.String())
	}
	if got := pb.GetIntroducedBy(); got != "" {
		t.Errorf("set introduced by: got %q, want empty", got)
	}

	// Parent setters accept the wrapper types.
	cfg.SetDevices([]DeviceConfiguration{{
		DeviceConfiguration: syncthingv2.DeviceConfiguration_builder{
			DeviceId: new(other.String()),
		}.Build(),
	}})
	if got := cfg.Configuration.GetDevices()[0].GetDeviceId(); got != other.String() {
		t.Errorf("set devices: got %q, want %q", got, other.String())
	}

	// The wrappers are views over the same underlying messages; the
	// configuration still marshals through the wrapper API.
	if _, err := Marshal(cfg); err != nil {
		t.Errorf("marshal through wrapper: %v", err)
	}
}
