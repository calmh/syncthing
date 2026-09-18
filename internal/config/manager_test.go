// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/proto"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/events"
	"github.com/syncthing/syncthing/lib/protocol"
)

// testManager runs a Manager until stopped.
type testManager struct {
	*Manager
	cancel context.CancelFunc
	done   chan struct{}
}

func startManager(m *Manager) *testManager {
	tm := &testManager{
		Manager: m,
		done:    make(chan struct{}),
	}
	ctx, cancel := context.WithCancel(context.Background())
	tm.cancel = cancel
	go func() {
		_ = m.Serve(ctx)
		close(tm.done)
	}()
	return tm
}

func (m *testManager) stop() {
	m.cancel()
	<-m.done
}

type requiresRestartCommitter struct {
	committed chan struct{}
}

func (requiresRestartCommitter) VerifyConfiguration(_, _ Configuration) error {
	return nil
}

func (c requiresRestartCommitter) CommitConfiguration(_, _ Configuration) bool {
	select {
	case c.committed <- struct{}{}:
	default:
	}
	return false
}

func (requiresRestartCommitter) String() string { return "requiresRestart" }

type validationErrorCommitter struct{}

func (validationErrorCommitter) VerifyConfiguration(_, _ Configuration) error {
	return errors.New("some error")
}

func (validationErrorCommitter) CommitConfiguration(_, _ Configuration) bool {
	return true
}

func (validationErrorCommitter) String() string { return "validationError" }

func replace(t *testing.T, w *Manager, to Configuration) {
	t.Helper()
	waiter, err := w.Modify(func(cfg *Configuration) {
		*cfg = to
	})
	if err != nil {
		t.Fatal(err)
	}
	waiter.Wait()
}

func TestManagerReplaceCommit(t *testing.T) {
	id := protocol.NewDeviceID([]byte("manager"))
	w := startManager(Manage("/dev/null", Configuration{syncthingv2.Configuration_builder{}.Build()}, id, events.NoopLogger))
	defer w.stop()

	if w.RawCopy().GetVersion() != 0 {
		t.Fatal("Config incorrect")
	}

	// Replace config. We should get back a clean response and the config
	// should change.

	replace(t, w.Manager, Configuration{syncthingv2.Configuration_builder{Version: new(int32(1))}.Build()})
	if w.RequiresRestart() {
		t.Fatal("Should not require restart")
	}
	if w.RawCopy().GetVersion() != 1 {
		t.Fatal("Config should have changed")
	}

	// Now with a subscriber requiring restart. We should get a clean
	// response but with the restart flag set, and the config should change.

	sub0 := requiresRestartCommitter{committed: make(chan struct{}, 1)}
	if cfg := w.Subscribe(sub0); cfg.GetVersion() != 1 {
		t.Fatal("Subscribe should return the current config")
	}

	replace(t, w.Manager, Configuration{syncthingv2.Configuration_builder{Version: new(int32(2))}.Build()})
	<-sub0.committed
	if !w.RequiresRestart() {
		t.Fatal("Should require restart")
	}
	if w.RawCopy().GetVersion() != 2 {
		t.Fatal("Config should have changed")
	}

	// Now with a subscriber that throws a validation error. The config
	// should not change.

	w.Subscribe(validationErrorCommitter{})

	_, err := w.Modify(func(cfg *Configuration) {
		cfg.SetVersion(3)
	})
	if err == nil {
		t.Fatal("Should have a validation error")
	}
	if !w.RequiresRestart() {
		t.Fatal("Should still require restart")
	}
	if w.RawCopy().GetVersion() != 2 {
		t.Fatal("Config should not have changed")
	}

	// A configuration that does not pass schema validation, here a device
	// ID with incorrect check digits, is likewise rejected.

	corrupt := []byte(protocol.NewDeviceID([]byte("managerbad")).String())
	if corrupt[0] == 'A' {
		corrupt[0] = 'B'
	} else {
		corrupt[0] = 'A'
	}
	if _, err := protocol.DeviceIDFromString(string(corrupt)); err == nil {
		t.Fatal("Test bug: corrupted device ID still parses")
	}
	_, err = w.Modify(func(cfg *Configuration) {
		cfg.SetDevices([]DeviceConfiguration{{
			DeviceConfiguration: syncthingv2.DeviceConfiguration_builder{
				DeviceId: new(string(corrupt)),
			}.Build(),
		}})
	})
	if err == nil {
		t.Fatal("Should have a validation error for bad device ID")
	}
	if w.RawCopy().GetVersion() != 2 {
		t.Fatal("Config should not have changed")
	}
}

func TestManagerSaveLoad(t *testing.T) {
	id := protocol.NewDeviceID([]byte("saveload"))
	path := filepath.Join(t.TempDir(), "config.yaml")

	// The API key is set explicitly, as preparation generates a random
	// one when it is unset.
	cfg := Configuration{syncthingv2.Configuration_builder{
		Version: new(int32(52)),
		Devices: []*syncthingv2.DeviceConfiguration{syncthingv2.DeviceConfiguration_builder{
			DeviceId: new(id.String()),
			Name:     new("myself"),
		}.Build()},
		Gui: syncthingv2.GUIConfiguration_builder{
			ApiKey: new("kO8nJgP7tY2wZq4x"),
		}.Build(),
	}.Build()}
	w := startManager(Manage(path, cfg, id, events.NoopLogger))

	waiter, err := w.Modify(func(cfg *Configuration) {
		cfg.SetFolders([]FolderConfiguration{{
			FolderConfiguration: syncthingv2.FolderConfiguration_builder{
				Id:   new("f1"),
				Path: new("/srv/sync/f1"),
			}.Build(),
		}})
	})
	if err != nil {
		t.Fatal(err)
	}
	waiter.Wait()

	// Stopping the manager saves the configuration.
	w.stop()

	loaded, err := Load(path, id, events.NoopLogger)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	// Loading prepares the configuration; do the same to the expected
	// value before comparing.
	want := w.RawCopy()
	if err := Prepare(&want, id); err != nil {
		t.Fatal(err)
	}
	got := loaded.RawCopy()
	if !proto.Equal(want.Configuration, got.Configuration) {
		t.Errorf("configuration did not survive save/load:\ngot:  %v\nwant: %v", got, want)
	}
	if _, ok := loaded.Folder("f1"); !ok {
		t.Error("folder f1 missing after load")
	}
	// Unset fields stay unset; the getter still reports the default.
	if loaded.RawCopy().GetFolders()[0].HasRescanIntervalS() {
		t.Error("rescanIntervalS should be unset after load")
	}
	if v := loaded.RawCopy().GetFolders()[0].GetRescanIntervalS(); v != 3600 {
		t.Errorf("rescanIntervalS: got %d, want default 3600", v)
	}
}

func TestManagerAccessorsAndRemove(t *testing.T) {
	id := protocol.NewDeviceID([]byte("remove"))
	other := protocol.NewDeviceID([]byte("removeother"))
	cfg := Configuration{syncthingv2.Configuration_builder{
		Version: new(int32(52)),
		Folders: []*syncthingv2.FolderConfiguration{syncthingv2.FolderConfiguration_builder{
			Id:   new("f1"),
			Path: new("/srv/sync/f1"),
			Devices: []*syncthingv2.FolderDeviceConfiguration{syncthingv2.FolderDeviceConfiguration_builder{
				DeviceId:           new(id.String()),
				EncryptionPassword: new("secret"),
			}.Build()},
		}.Build()},
		Devices: []*syncthingv2.DeviceConfiguration{
			syncthingv2.DeviceConfiguration_builder{DeviceId: new(id.String())}.Build(),
			syncthingv2.DeviceConfiguration_builder{
				DeviceId: new(other.String()),
				IgnoredFolders: []*syncthingv2.ObservedFolder{syncthingv2.ObservedFolder_builder{
					Id: new("ignored"),
				}.Build()},
			}.Build(),
		},
		RemoteIgnoredDevices: []*syncthingv2.ObservedDevice{syncthingv2.ObservedDevice_builder{
			DeviceId: new(other.String()),
		}.Build()},
	}.Build()}
	w := startManager(Manage("", cfg, id, events.NoopLogger))
	defer w.stop()

	// Accessors.
	if _, ok := w.Device(other); !ok {
		t.Error("device should be present")
	}
	devices := w.Devices()
	if len(devices) != 2 {
		t.Error("devices map incorrect")
	} else if _, ok := devices[other]; !ok {
		t.Error("devices map incorrect")
	}
	if len(w.DeviceList()) != 2 {
		t.Error("device list incorrect")
	}
	if _, ok := w.Folder("f1"); !ok {
		t.Error("folder should be present")
	}
	if len(w.Folders()) != 1 {
		t.Error("folders map incorrect")
	}
	if pw := w.FolderPasswords(id)["f1"]; pw != "secret" {
		t.Errorf("folder password: got %q, want %q", pw, "secret")
	}
	if !w.IgnoredDevice(other) {
		t.Error("device should be ignored")
	}
	if !w.IgnoredFolder(other, "ignored") {
		t.Error("folder should be ignored")
	}
	if w.DefaultIgnores() != nil {
		t.Error("default ignores should be unset")
	}
	if w.RawCopy().HasGui() || w.RawCopy().HasOptions() || w.LDAP() != nil {
		t.Error("unset sections should be nil")
	}

	// Removing a device and a folder takes effect and is validated.
	waiter, err := w.RemoveDevice(other)
	if err != nil {
		t.Fatal(err)
	}
	waiter.Wait()
	if _, ok := w.Device(other); ok {
		t.Error("device should be gone")
	}
	if len(w.DeviceList()) != 1 {
		t.Error("device list should have one entry")
	}

	waiter, err = w.RemoveFolder("f1")
	if err != nil {
		t.Fatal(err)
	}
	waiter.Wait()
	if _, ok := w.Folder("f1"); ok {
		t.Error("folder should be gone")
	}
	if len(w.FolderList()) != 0 {
		t.Error("folder list should be empty")
	}
}
