// Copyright (C) 2016 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package fs

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"

	"github.com/syncthing/syncthing/internal/slogutil"
)

type logFilesystem struct {
	Filesystem

	// Number of filesystem layers on top of logFilesystem to skip when looking
	// for the true caller of the filesystem
	layers int
}

func newLogFilesystem(fs Filesystem, layers int) *logFilesystem {
	return &logFilesystem{
		Filesystem: fs,
		layers:     layers,
	}
}

func (fs *logFilesystem) getCaller() string {
	_, file, line, ok := runtime.Caller(fs.layers + 1)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func (fs *logFilesystem) Chmod(name string, mode FileMode) error {
	err := fs.Filesystem.Chmod(name, mode)
	slog.Debug("Chmod", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Any("mode", mode), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	err := fs.Filesystem.Chtimes(name, atime, mtime)
	slog.Debug("Chtimes", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Time("atime", atime), slog.Time("mtime", mtime), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) Create(name string) (File, error) {
	file, err := fs.Filesystem.Create(name)
	slog.Debug("Create", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Any("file", file), slogutil.Error(err))
	return file, err
}

func (fs *logFilesystem) CreateSymlink(target, name string) error {
	err := fs.Filesystem.CreateSymlink(target, name)
	slog.Debug("CreateSymlink", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("target", target), slog.String("name", name), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) DirNames(name string) ([]string, error) {
	names, err := fs.Filesystem.DirNames(name)
	slog.Debug("DirNames", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Any("names", names), slogutil.Error(err))
	return names, err
}

func (fs *logFilesystem) Lstat(name string) (FileInfo, error) {
	info, err := fs.Filesystem.Lstat(name)
	slog.Debug("Lstat", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Any("info", info), slogutil.Error(err))
	return info, err
}

func (fs *logFilesystem) Mkdir(name string, perm FileMode) error {
	err := fs.Filesystem.Mkdir(name, perm)
	slog.Debug("Mkdir", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Any("perm", perm), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) MkdirAll(name string, perm FileMode) error {
	err := fs.Filesystem.MkdirAll(name, perm)
	slog.Debug("MkdirAll", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Any("perm", perm), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) Open(name string) (File, error) {
	file, err := fs.Filesystem.Open(name)
	slog.Debug("Open", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Any("file", file), slogutil.Error(err))
	return file, err
}

func (fs *logFilesystem) OpenFile(name string, flags int, mode FileMode) (File, error) {
	file, err := fs.Filesystem.OpenFile(name, flags, mode)
	slog.Debug("OpenFile", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Int("flags", flags), slog.Any("mode", mode), slog.Any("file", file), slogutil.Error(err))
	return file, err
}

func (fs *logFilesystem) ReadSymlink(name string) (string, error) {
	target, err := fs.Filesystem.ReadSymlink(name)
	slog.Debug("ReadSymlink", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.String("target", target), slogutil.Error(err))
	return target, err
}

func (fs *logFilesystem) Remove(name string) error {
	err := fs.Filesystem.Remove(name)
	slog.Debug("Remove", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) RemoveAll(name string) error {
	err := fs.Filesystem.RemoveAll(name)
	slog.Debug("RemoveAll", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) Rename(oldname, newname string) error {
	err := fs.Filesystem.Rename(oldname, newname)
	slog.Debug("Rename", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("oldName", oldname), slog.String("newName", newname), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) Stat(name string) (FileInfo, error) {
	info, err := fs.Filesystem.Stat(name)
	slog.Debug("Stat", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Any("info", info), slogutil.Error(err))
	return info, err
}

func (fs *logFilesystem) SymlinksSupported() bool {
	supported := fs.Filesystem.SymlinksSupported()
	slog.Debug("SymlinksSupported", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.Bool("supported", supported))
	return supported
}

func (fs *logFilesystem) Walk(root string, walkFn WalkFunc) error {
	err := fs.Filesystem.Walk(root, walkFn)
	slog.Debug("Walk", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("root", root), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) Watch(path string, ignore Matcher, ctx context.Context, ignorePerms bool) (<-chan Event, <-chan error, error) {
	evChan, errChan, err := fs.Filesystem.Watch(path, ignore, ctx, ignorePerms)
	slog.Debug("Watch", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("path", path), slog.Bool("ignorePerms", ignorePerms), slogutil.Error(err))
	return evChan, errChan, err
}

func (fs *logFilesystem) Unhide(name string) error {
	err := fs.Filesystem.Unhide(name)
	slog.Debug("Unhide", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) Hide(name string) error {
	err := fs.Filesystem.Hide(name)
	slog.Debug("Hide", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slogutil.Error(err))
	return err
}

func (fs *logFilesystem) Glob(name string) ([]string, error) {
	names, err := fs.Filesystem.Glob(name)
	slog.Debug("Glob", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("pattern", name), slog.Any("names", names), slogutil.Error(err))
	return names, err
}

func (fs *logFilesystem) Roots() ([]string, error) {
	roots, err := fs.Filesystem.Roots()
	slog.Debug("Roots", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.Any("roots", roots), slogutil.Error(err))
	return roots, err
}

func (fs *logFilesystem) Usage(name string) (Usage, error) {
	usage, err := fs.Filesystem.Usage(name)
	slog.Debug("Usage", slog.String("caller", fs.getCaller()), slog.Any("fsType", fs.Type()), slog.String("uri", fs.URI()), slog.String("name", name), slog.Any("usage", usage), slogutil.Error(err))
	return usage, err
}

func (fs *logFilesystem) underlying() (Filesystem, bool) {
	return fs.Filesystem, true
}
