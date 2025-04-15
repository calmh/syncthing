// Copyright (C) 2016 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package model

import (
	"context"
	"os"
	"testing"
	"time"

	configv1 "github.com/syncthing/syncthing/internal/config/v1"

	"github.com/syncthing/syncthing/internal/db/sqlite"
	"github.com/syncthing/syncthing/lib/events"
	"github.com/syncthing/syncthing/lib/fs"
	"github.com/syncthing/syncthing/lib/ignore"
	"github.com/syncthing/syncthing/lib/protocol"
	"github.com/syncthing/syncthing/lib/protocol/mocks"
	"github.com/syncthing/syncthing/lib/rand"
)

var (
	myID, device1, device2  protocol.DeviceID
	defaultCfgWrapper       configv1.Wrapper
	defaultCfgWrapperCancel context.CancelFunc
	defaultFolderConfig     configv1.FolderConfiguration
	defaultCfg              configv1.Configuration
	defaultAutoAcceptCfg    configv1.Configuration
	device1Conn             = &mocks.Connection{}
	device2Conn             = &mocks.Connection{}
)

func init() {
	myID, _ = protocol.DeviceIDFromString("ZNWFSWE-RWRV2BD-45BLMCV-LTDE2UR-4LJDW6J-R5BPWEB-TXD27XJ-IZF5RA4")
	device1, _ = protocol.DeviceIDFromString("AIR6LPZ-7K4PTTV-UXQSMUU-CPQ5YWH-OEDFIIQ-JUG777G-2YQXXR5-YD6AWQR")
	device2, _ = protocol.DeviceIDFromString("GYRZZQB-IRNPV4Z-T7TC52W-EQYJ3TT-FDQW6MW-DFLMU42-SSSU6EM-FBK2VAY")
	device1Conn.DeviceIDReturns(device1)
	device1Conn.ConnectionIDReturns(rand.String(16))
	device2Conn.DeviceIDReturns(device2)
	device2Conn.ConnectionIDReturns(rand.String(16))

	cfg := configv1.New(myID)
	cfg.Options.MinHomeDiskFree.Value = 0 // avoids unnecessary free space checks
	defaultCfgWrapper, defaultCfgWrapperCancel = newConfigWrapper(cfg)

	defaultFolderConfig = newFolderConfig()

	waiter, _ := defaultCfgWrapper.Modify(func(cfg *configv1.Configuration) {
		cfg.SetDevice(newDeviceConfiguration(cfg.Defaults.Device, device1, "device1"))
		cfg.SetFolder(defaultFolderConfig)
		cfg.Options.KeepTemporariesH = 1
	})
	waiter.Wait()

	defaultCfg = defaultCfgWrapper.RawCopy()

	defaultAutoAcceptCfg = configv1.Configuration{
		Version: configv1.CurrentVersion,
		Devices: []configv1.DeviceConfiguration{
			{
				DeviceID: myID, // self
			},
			{
				DeviceID:          device1,
				AutoAcceptFolders: true,
			},
			{
				DeviceID:          device2,
				AutoAcceptFolders: true,
			},
		},
		Defaults: configv1.Defaults{
			Folder: configv1.FolderConfiguration{
				FilesystemType: configv1.FilesystemTypeFake,
				Path:           rand.String(32),
			},
		},
		Options: configv1.OptionsConfiguration{
			MinHomeDiskFree: configv1.Size{}, // avoids unnecessary free space checks
		},
	}
}

func newConfigWrapper(cfg configv1.Configuration) (configv1.Wrapper, context.CancelFunc) {
	wrapper := configv1.Wrap("", cfg, myID, events.NoopLogger)
	ctx, cancel := context.WithCancel(context.Background())
	go wrapper.Serve(ctx)
	return wrapper, cancel
}

func newDefaultCfgWrapper() (configv1.Wrapper, configv1.FolderConfiguration, context.CancelFunc) {
	w, cancel := newConfigWrapper(defaultCfgWrapper.RawCopy())
	fcfg := newFolderConfig()
	_, _ = w.Modify(func(cfg *configv1.Configuration) {
		cfg.SetFolder(fcfg)
	})
	return w, fcfg, cancel
}

func newFolderConfig() configv1.FolderConfiguration {
	cfg := newFolderConfiguration(defaultCfgWrapper, "default", "default", configv1.FilesystemTypeFake, rand.String(32)+"?content=true")
	cfg.FSWatcherEnabled = false
	cfg.PullerDelayS = 0
	cfg.Devices = append(cfg.Devices, configv1.FolderDeviceConfiguration{DeviceID: device1})
	return cfg
}

func setupModelWithConnection(t testing.TB) (*testModel, *fakeConnection, configv1.FolderConfiguration, context.CancelFunc) {
	t.Helper()
	w, fcfg, cancel := newDefaultCfgWrapper()
	m, fc := setupModelWithConnectionFromWrapper(t, w)
	return m, fc, fcfg, cancel
}

func setupModelWithConnectionFromWrapper(t testing.TB, w configv1.Wrapper) (*testModel, *fakeConnection) {
	t.Helper()
	m := setupModel(t, w)

	fc := addFakeConn(m, device1, "default")
	fc.folder = "default"

	_ = m.ScanFolder("default")

	return m, fc
}

func setupModel(t testing.TB, w configv1.Wrapper) *testModel {
	t.Helper()
	m := newModel(t, w, myID, nil)
	m.ServeBackground()
	<-m.started

	m.ScanFolders()

	return m
}

type testModel struct {
	*model
	t        testing.TB
	cancel   context.CancelFunc
	evCancel context.CancelFunc
	stopped  chan struct{}
}

func newModel(t testing.TB, cfg configv1.Wrapper, id protocol.DeviceID, protectedFiles []string) *testModel {
	t.Helper()
	evLogger := events.NewLogger()
	mdb, err := sqlite.OpenTemp()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		mdb.Close()
	})
	m := NewModel(cfg, id, mdb, protectedFiles, evLogger, protocol.NewKeyGenerator()).(*model)
	ctx, cancel := context.WithCancel(context.Background())
	go evLogger.Serve(ctx)
	return &testModel{
		model:    m,
		evCancel: cancel,
		stopped:  make(chan struct{}),
		t:        t,
	}
}

func (m *testModel) ServeBackground() {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	go func() {
		m.model.Serve(ctx)
		close(m.stopped)
	}()
	<-m.started
}

func (m *testModel) testCurrentFolderFile(folder string, file string) (protocol.FileInfo, bool) {
	f, ok, err := m.model.CurrentFolderFile(folder, file)
	must(m.t, err)
	return f, ok
}

func (m *testModel) testCompletion(device protocol.DeviceID, folder string) FolderCompletion {
	comp, err := m.Completion(device, folder)
	must(m.t, err)
	return comp
}

func cleanupModel(m *testModel) {
	if m.cancel != nil {
		m.cancel()
		<-m.stopped
	}
	m.evCancel()
	m.sdb.Close()
	os.Remove(m.cfg.ConfigPath())
}

func cleanupModelAndRemoveDir(m *testModel, dir string) {
	cleanupModel(m)
	os.RemoveAll(dir)
}

type alwaysChangedKey struct {
	fs   fs.Filesystem
	name string
}

// alwaysChanges is an ignore.ChangeDetector that always returns true on Changed()
type alwaysChanged struct {
	seen map[alwaysChangedKey]struct{}
}

func newAlwaysChanged() *alwaysChanged {
	return &alwaysChanged{
		seen: make(map[alwaysChangedKey]struct{}),
	}
}

func (c *alwaysChanged) Remember(fs fs.Filesystem, name string, _ time.Time) {
	c.seen[alwaysChangedKey{fs, name}] = struct{}{}
}

func (c *alwaysChanged) Reset() {
	c.seen = make(map[alwaysChangedKey]struct{})
}

func (c *alwaysChanged) Seen(fs fs.Filesystem, name string) bool {
	_, ok := c.seen[alwaysChangedKey{fs, name}]
	return ok
}

func (*alwaysChanged) Changed() bool {
	return true
}

// Reach in and update the ignore matcher to one that always does
// reloads when asked to, instead of checking file mtimes. This is
// because we will be changing the files on disk often enough that the
// mtimes will be unreliable to determine change status.
func folderIgnoresAlwaysReload(t testing.TB, m *testModel, fcfg configv1.FolderConfiguration) {
	t.Helper()
	m.removeFolder(fcfg)
	ignores := ignore.New(fcfg.Filesystem(), ignore.WithCache(true), ignore.WithChangeDetector(newAlwaysChanged()))
	m.mut.Lock()
	m.addAndStartFolderLockedWithIgnores(fcfg, ignores)
	m.mut.Unlock()
}

func basicClusterConfig(local, remote protocol.DeviceID, folders ...string) *protocol.ClusterConfig {
	var cc protocol.ClusterConfig
	for _, folder := range folders {
		cc.Folders = append(cc.Folders, protocol.Folder{
			ID: folder,
			Devices: []protocol.Device{
				{
					ID: local,
				},
				{
					ID: remote,
				},
			},
		})
	}
	return &cc
}

func localIndexUpdate(m *testModel, folder string, fs []protocol.FileInfo) {
	m.sdb.Update(folder, protocol.LocalDeviceID, fs)
	seq, err := m.sdb.GetDeviceSequence(folder, protocol.LocalDeviceID)
	if err != nil {
		panic(err)
	}
	filenames := make([]string, len(fs))
	for i, file := range fs {
		filenames[i] = file.Name
	}
	m.evLogger.Log(events.LocalIndexUpdated, map[string]interface{}{
		"folder":    folder,
		"items":     len(fs),
		"filenames": filenames,
		"sequence":  seq,
		"version":   seq, // legacy for sequence
	})
}

func newDeviceConfiguration(defaultCfg configv1.DeviceConfiguration, id protocol.DeviceID, name string) configv1.DeviceConfiguration {
	cfg := defaultCfg.Copy()
	cfg.DeviceID = id
	cfg.Name = name
	return cfg
}

func replace(t testing.TB, w configv1.Wrapper, to configv1.Configuration) {
	t.Helper()
	waiter, err := w.Modify(func(cfg *configv1.Configuration) {
		*cfg = to
	})
	if err != nil {
		t.Fatal(err)
	}
	waiter.Wait()
}

func pauseFolder(t testing.TB, w configv1.Wrapper, id string, paused bool) {
	t.Helper()
	waiter, err := w.Modify(func(cfg *configv1.Configuration) {
		_, i, _ := cfg.Folder(id)
		cfg.Folders[i].Paused = paused
	})
	if err != nil {
		t.Fatal(err)
	}
	waiter.Wait()
}

func setFolder(t testing.TB, w configv1.Wrapper, fcfg configv1.FolderConfiguration) {
	t.Helper()
	waiter, err := w.Modify(func(cfg *configv1.Configuration) {
		cfg.SetFolder(fcfg)
	})
	if err != nil {
		t.Fatal(err)
	}
	waiter.Wait()
}

func pauseDevice(t testing.TB, w configv1.Wrapper, id protocol.DeviceID, paused bool) {
	t.Helper()
	waiter, err := w.Modify(func(cfg *configv1.Configuration) {
		_, i, _ := cfg.Device(id)
		cfg.Devices[i].Paused = paused
	})
	if err != nil {
		t.Fatal(err)
	}
	waiter.Wait()
}

func setDevice(t testing.TB, w configv1.Wrapper, device configv1.DeviceConfiguration) {
	t.Helper()
	waiter, err := w.Modify(func(cfg *configv1.Configuration) {
		cfg.SetDevice(device)
	})
	if err != nil {
		t.Fatal(err)
	}
	waiter.Wait()
}

func addDevice2(t testing.TB, w configv1.Wrapper, fcfg configv1.FolderConfiguration) {
	waiter, err := w.Modify(func(cfg *configv1.Configuration) {
		cfg.SetDevice(newDeviceConfiguration(cfg.Defaults.Device, device2, "device2"))
		fcfg.Devices = append(fcfg.Devices, configv1.FolderDeviceConfiguration{DeviceID: device2})
		cfg.SetFolder(fcfg)
	})
	must(t, err)
	waiter.Wait()
}

func writeFile(t testing.TB, filesystem fs.Filesystem, name string, data []byte) {
	t.Helper()
	fd, err := filesystem.Create(name)
	must(t, err)
	defer fd.Close()
	_, err = fd.Write(data)
	must(t, err)
}

func writeFilePerm(t testing.TB, filesystem fs.Filesystem, name string, data []byte, perm fs.FileMode) {
	t.Helper()
	writeFile(t, filesystem, name, data)
	must(t, filesystem.Chmod(name, perm))
}
