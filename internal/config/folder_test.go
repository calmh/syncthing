// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"errors"
	"testing"
	"time"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/fs"
	"github.com/syncthing/syncthing/lib/protocol"
	"github.com/syncthing/syncthing/lib/rand"
)

func fakeFolder(id string) FolderConfiguration {
	// The fake filesystem is content addressed by the URI; use a unique
	// one per folder to get isolation between tests.
	return FolderConfiguration{syncthingv2.FolderConfiguration_builder{
		Id:             new(id),
		Path:           new(rand.String(16) + "?nostfolder=true"),
		FilesystemType: new(syncthingv2.FilesystemType_FILESYSTEM_TYPE_FAKE),
	}.Build()}
}

func TestFolderFilesystem(t *testing.T) {
	// A basic filesystem of the right path.
	folder := FolderConfiguration{syncthingv2.FolderConfiguration_builder{
		Path: new("/srv/sync/f1"),
	}.Build()}
	ffs := folder.Filesystem()
	if ffs.Type() != fs.FilesystemTypeBasic {
		t.Errorf("filesystem type: got %v, want basic", ffs.Type())
	}
	if ffs.URI() != "/srv/sync/f1" {
		t.Errorf("filesystem URI: got %q", ffs.URI())
	}

	// A fake filesystem.
	folder = FolderConfiguration{syncthingv2.FolderConfiguration_builder{
		Path:           new("fake-uri"),
		FilesystemType: new(syncthingv2.FilesystemType_FILESYSTEM_TYPE_FAKE),
	}.Build()}
	if ffs = folder.Filesystem(); ffs.Type() != fs.FilesystemTypeFake {
		t.Errorf("filesystem type: got %v, want fake", ffs.Type())
	}
}

func TestFolderMarker(t *testing.T) {
	folder := fakeFolder("marker")

	// The fake filesystem has a root directory, so the marker is what is
	// missing.
	if err := folder.CreateRoot(); err != nil {
		t.Fatalf("create root: %v", err)
	}
	if err := folder.CheckPath(); !errors.Is(err, ErrMarkerMissing) {
		t.Fatalf("check path: got %v, want ErrMarkerMissing", err)
	}
	if err := folder.CreateMarker(); err != nil {
		t.Fatalf("create marker: %v", err)
	}
	if err := folder.CheckPath(); err != nil {
		t.Fatalf("check path after marker: %v", err)
	}

	// Removing the marker brings back the error.
	if err := folder.RemoveMarker(); err != nil {
		t.Fatalf("remove marker: %v", err)
	}
	if err := folder.CheckPath(); !errors.Is(err, ErrMarkerMissing) {
		t.Fatalf("check path: got %v, want ErrMarkerMissing", err)
	}
}

func TestFolderComputedGetters(t *testing.T) {
	id := protocol.NewDeviceID([]byte("foldercomputed"))
	other := protocol.NewDeviceID([]byte("folderother"))
	folder := FolderConfiguration{syncthingv2.FolderConfiguration_builder{
		Id:   new("f1"),
		Path: new("/srv/sync/f1"),
		Devices: []*syncthingv2.FolderDeviceConfiguration{
			syncthingv2.FolderDeviceConfiguration_builder{DeviceId: new(id.String())}.Build(),
			syncthingv2.FolderDeviceConfiguration_builder{DeviceId: new(other.String())}.Build(),
		},
	}.Build()}

	if got := folder.DeviceIDs(); len(got) != 2 || !got[0].Equals(id) || !got[1].Equals(other) {
		t.Errorf("device IDs: got %v", got)
	}
	if !folder.SharedWith(other) {
		t.Error("folder should be shared with other")
	}
	if folder.SharedWith(protocol.NewDeviceID([]byte("stranger"))) {
		t.Error("folder should not be shared with stranger")
	}

	if got := folder.Description(); got != `"f1" (f1)` && got != "f1" {
		// Without a label the description is the ID.
		if got != "f1" {
			t.Errorf("description: got %q", got)
		}
	}
	folder.SetLabel("Folder One")
	if got := folder.Description(); got != `"Folder One" (f1)` {
		t.Errorf("description: got %q", got)
	}

	// The modification time window is the configured value in seconds.
	folder.SetModTimeWindowS(5)
	if got := folder.ModTimeWindow(); got != 5*time.Second {
		t.Errorf("mod time window: got %v, want 5s", got)
	}
	folder.ClearModTimeWindowS()
	if got := folder.ModTimeWindow(); got != 0 {
		t.Errorf("mod time window: got %v, want 0", got)
	}
}
