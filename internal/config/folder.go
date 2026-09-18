// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/disk"

	"github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/internal/slogutil"
	"github.com/syncthing/syncthing/lib/build"
	"github.com/syncthing/syncthing/lib/fs"
	"github.com/syncthing/syncthing/lib/protocol"
)

var (
	ErrPathNotDirectory = errors.New("folder path not a directory")
	ErrPathMissing      = errors.New("folder path missing")
	ErrMarkerMissing    = errors.New("folder marker missing (this indicates potential data loss, search docs/forum to get information about how to proceed)")
)

// Filesystem creates a filesystem for the path and options of this
// folder.
//
// This is intentionally not a pointer method, because things like
// cfg.GetFolders()[0].Filesystem() should be valid.
func (f FolderConfiguration) Filesystem(extraOpts ...fs.Option) fs.Filesystem {
	var opts []fs.Option
	if f.GetFilesystemType() == syncthingv2.FilesystemType_FILESYSTEM_TYPE_BASIC && f.GetJunctionsAsDirs() {
		opts = append(opts, new(fs.OptionJunctionsAsDirs))
	}
	if !f.GetCaseSensitiveFs() {
		opts = append(opts, new(fs.OptionDetectCaseConflicts))
	}
	opts = append(opts, extraOpts...)
	return fs.NewFilesystem(filesystemTypeToFS(f.GetFilesystemType()), f.GetPath(), opts...)
}

// ModTimeWindow returns the allowed modification timestamp difference
// when comparing files for equivalence. On Android, a default of two
// seconds is used for filesystems with unstable timestamps.
func (f FolderConfiguration) ModTimeWindow() time.Duration {
	dur := time.Duration(f.GetModTimeWindowS()) * time.Second
	if f.GetModTimeWindowS() < 1 && build.IsAndroid {
		if usage, err := disk.Usage(f.Filesystem().URI()); err != nil {
			dur = 2 * time.Second
			slog.Debug("Detecting FS on android: setting mtime window to 2s", slogutil.FilePath(f.GetPath()), slogutil.Error(err))
		} else if strings.HasPrefix(strings.ToLower(usage.Fstype), "ext2") || strings.HasPrefix(strings.ToLower(usage.Fstype), "ext3") || strings.HasPrefix(strings.ToLower(usage.Fstype), "ext4") {
			slog.Debug("Detecting FS on android: leaving mtime window at 0", slogutil.FilePath(f.GetPath()), "fstype", usage.Fstype)
		} else {
			dur = 2 * time.Second
			slog.Debug("Detecting FS on android: setting mtime window to 2s", slogutil.FilePath(f.GetPath()), "fstype", usage.Fstype)
		}
	}
	return dur
}

// CreateMarker creates the folder marker, if it is missing. This should
// only be called when the folder uses the default marker name.
func (f *FolderConfiguration) CreateMarker() error {
	if err := f.CheckPath(); !errors.Is(err, ErrMarkerMissing) {
		return err
	}
	if f.GetMarkerName() != DefaultMarkerName {
		// Folder uses a non-default marker so we shouldn't mess with it.
		// Pretend we created it and let the subsequent health checks sort
		// out the actual situation.
		return nil
	}

	ffs := f.Filesystem()

	// Create the marker as a directory
	err := ffs.Mkdir(DefaultMarkerName, fs.ModePerm)
	if err != nil {
		return err
	}

	// Create a file inside it, reducing the risk of the marker directory
	// being removed by automated cleanup tools.
	markerFile := filepath.Join(DefaultMarkerName, f.markerFilename())
	if err := fs.WriteFile(ffs, markerFile, f.markerContents(), 0o666); err != nil {
		return err
	}

	// Sync & hide the containing directory
	if dir, err := ffs.Open("."); err != nil {
		slog.Debug("Folder marker: open . failed", slogutil.Error(err))
	} else if err := dir.Sync(); err != nil {
		slog.Debug("Folder marker: fsync . failed", slogutil.Error(err))
	}
	ffs.Hide(DefaultMarkerName)

	return nil
}

// RemoveMarker removes the folder marker, if it exists.
func (f *FolderConfiguration) RemoveMarker() error {
	ffs := f.Filesystem()
	_ = ffs.Remove(filepath.Join(DefaultMarkerName, f.markerFilename()))
	return ffs.Remove(DefaultMarkerName)
}

func (f *FolderConfiguration) markerFilename() string {
	h := sha256.Sum256([]byte(f.GetId()))
	return fmt.Sprintf("syncthing-folder-%x.txt", h[:3])
}

func (f *FolderConfiguration) markerContents() []byte {
	var buf bytes.Buffer
	buf.WriteString("# This directory is a Syncthing folder marker.\n# Do not delete.\n\n")
	fmt.Fprintf(&buf, "folderID: %s\n", f.GetId())
	fmt.Fprintf(&buf, "created: %s\n", time.Now().Format(time.RFC3339))
	return buf.Bytes()
}

// CheckPath returns nil if the folder root exists and contains the
// marker file.
func (f *FolderConfiguration) CheckPath() error {
	return f.checkFilesystemPath(f.Filesystem(), ".")
}

func (f *FolderConfiguration) checkFilesystemPath(ffs fs.Filesystem, path string) error {
	fi, err := ffs.Stat(path)
	if err != nil {
		if !fs.IsNotExist(err) {
			return err
		}
		return ErrPathMissing
	}

	// Users might have the root directory as a symlink or reparse point.
	// Furthermore, OneDrive bullcrap uses a magic reparse point to the
	// cloudz... Yet it's impossible for this to happen, as filesystem
	// adds a trailing path separator to the root, so even if you point
	// the filesystem at a file Stat ends up calling stat on C:\dir\file\
	// which, fails with "is not a directory" in the error check above,
	// and we don't even get to here.
	if !fi.IsDir() && !fi.IsSymlink() {
		return ErrPathNotDirectory
	}

	_, err = ffs.Stat(filepath.Join(path, f.GetMarkerName()))
	if err != nil {
		if !fs.IsNotExist(err) {
			return err
		}
		return ErrMarkerMissing
	}

	return nil
}

// CreateRoot creates the folder root directory, if it is missing.
func (f *FolderConfiguration) CreateRoot() error {
	filesystem := f.Filesystem()

	if _, err := filesystem.Stat("."); fs.IsNotExist(err) {
		return filesystem.MkdirAll(".", fs.ModePerm)
	} else if err != nil {
		return err
	}
	return nil
}

// Description returns a human readable representation of the folder.
func (f FolderConfiguration) Description() string {
	if f.GetLabel() == "" {
		return f.GetId()
	}
	return fmt.Sprintf("%q (%s)", f.GetLabel(), f.GetId())
}

// LogAttr returns a log attribute for the folder.
func (f FolderConfiguration) LogAttr() slog.Attr {
	if f.GetLabel() == "" || f.GetLabel() == f.GetId() {
		return slog.Group("folder", slog.String("id", f.GetId()), slog.String("type", f.GetType().String()))
	}
	return slog.Group("folder", slog.String("label", f.GetLabel()), slog.String("id", f.GetId()), slog.String("type", f.GetType().String()))
}

// DeviceIDs returns the IDs of the devices the folder is shared with.
func (f *FolderConfiguration) DeviceIDs() []protocol.DeviceID {
	deviceIDs := make([]protocol.DeviceID, len(f.GetDevices()))
	for i, n := range f.GetDevices() {
		deviceIDs[i] = n.GetDeviceId()
	}
	return deviceIDs
}

// SharedWith returns whether the folder is shared with the given device.
func (f FolderConfiguration) SharedWith(device protocol.DeviceID) bool {
	_, ok := f.Device(device)
	return ok
}

// CheckAvailableSpace checks that there is enough free space on the
// folder filesystem to receive req more bytes, and that the free space
// does not fall below the configured minimum.
func (f FolderConfiguration) CheckAvailableSpace(req uint64) error {
	val := SizeInBytes(f.MinDiskFree())
	if val <= 0 && !SizeIsPercentage(f.MinDiskFree()) {
		return nil
	}
	ffs := f.Filesystem()
	usage, err := ffs.Usage(".")
	if err != nil {
		return nil //nolint:nilerr
	}
	if err := checkAvailableSpace(req, f.MinDiskFree(), usage); err != nil {
		return fmt.Errorf("insufficient space in folder %v (%v): %w", f.Description(), ffs.URI(), err)
	}
	return nil
}

// filesystemTypeToFS maps the filesystem type to the filesystem
// implementation.
func filesystemTypeToFS(t syncthingv2.FilesystemType) fs.FilesystemType {
	if t == syncthingv2.FilesystemType_FILESYSTEM_TYPE_FAKE {
		return fs.FilesystemTypeFake
	}
	return fs.FilesystemTypeBasic
}
