// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build integration

// Package integration contains end-to-end tests that drive real Syncthing
// processes via the REST API.
package integration

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/protocol"
	"github.com/syncthing/syncthing/lib/rc"
)

// syncthingBinary is the prebuilt Syncthing binary, located relative to the
// integration package directory (the working directory of `go test`).
const syncthingBinary = "../bin/syncthing"

// instance represents one Syncthing process under test, with its own home
// directory and a single shared folder backed by an in-memory fakefs.
type instance struct {
	index     int
	home      string
	fakefs    string // root URI for the fakefs-backed folder
	deviceID  protocol.DeviceID
	guiAddr   string // host:port
	listenURL string // tcp://host:port
	process   *rc.Process
}

// newInstance creates a fresh Syncthing home, runs `syncthing generate`,
// and rewrites the config to use the instance's assigned ports plus a
// single "default" folder backed by an in-memory fakefs at root URI
// fakefs (see lib/fs/fakefs.go for the supported parameters, e.g.
// "a?files=1000&sizeavg=4096&seed=1"). It does not start the process.
//
// Ports are derived from the index: GUI on 8080+index, BEP on 22000+index.
// Indexes must be unique within a test, and tests sharing indexes cannot
// run in parallel.
func newInstance(t *testing.T, index int, fakefs string) *instance {
	t.Helper()

	home := filepath.Join(t.TempDir(), "home")

	gen := exec.Command(syncthingBinary, "generate", "--home", home, "--no-port-probing")
	gen.Stdout = os.Stdout
	gen.Stderr = os.Stderr
	if err := gen.Run(); err != nil {
		t.Fatalf("syncthing generate (instance %d): %v", index, err)
	}

	cert, err := tls.LoadX509KeyPair(filepath.Join(home, "cert.pem"), filepath.Join(home, "key.pem"))
	if err != nil {
		t.Fatal(err)
	}
	deviceID := protocol.NewDeviceID(cert.Certificate[0])

	guiAddr := fmt.Sprintf("127.0.0.1:%d", 8080+index)
	listenURL := fmt.Sprintf("tcp://127.0.0.1:%d", 22000+index)

	cfgPath := filepath.Join(home, "config.xml")
	cfg := readConfigXML(t, cfgPath, deviceID)

	cfg.GUI.RawAddress = guiAddr
	cfg.GUI.APIKey = rc.APIKey
	cfg.GUI.User = ""
	cfg.GUI.Password = ""

	cfg.Options.RawListenAddresses = []string{listenURL}
	cfg.Options.GlobalAnnEnabled = false
	cfg.Options.LocalAnnEnabled = false
	cfg.Options.RelaysEnabled = false
	cfg.Options.NATEnabled = false
	cfg.Options.URAccepted = -1
	cfg.Options.AutoUpgradeIntervalH = 0
	cfg.Options.CREnabled = false
	cfg.Options.StartBrowser = false

	folder := cfg.Defaults.Folder.Copy()
	folder.ID = "default"
	folder.Label = "default"
	folder.FilesystemType = config.FilesystemTypeFake
	folder.Path = fakefs
	folder.Devices = []config.FolderDeviceConfiguration{{DeviceID: deviceID}}
	cfg.Folders = []config.FolderConfiguration{folder}

	writeConfigXML(t, cfgPath, &cfg)

	return &instance{
		index:     index,
		home:      home,
		fakefs:    fakefs,
		deviceID:  deviceID,
		guiAddr:   guiAddr,
		listenURL: listenURL,
	}
}

// start launches the Syncthing process and waits for it to finish initial
// scans. The process is automatically stopped at the end of the test.
func (i *instance) start(t *testing.T) {
	t.Helper()

	logFile := filepath.Join(i.home, "syncthing.log")
	p := rc.NewProcess(i.guiAddr)
	if err := p.LogTo(logFile); err != nil {
		t.Fatal(err)
	}
	if err := p.Start(syncthingBinary, "--home", i.home, "--no-browser"); err != nil {
		t.Fatalf("start instance %d: %v", i.index, err)
	}
	i.process = p

	awaitStartup(t, p)

	t.Cleanup(func() {
		if _, err := p.Stop(); err != nil {
			t.Logf("stop instance %d: %v (log: %s)", i.index, err, logFile)
		}
	})
}

// awaitStartup wraps rc.Process.AwaitStartup with a timeout so a stuck
// instance fails the test rather than hanging it.
func awaitStartup(t *testing.T, p *rc.Process) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		p.AwaitStartup()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("timed out waiting for Syncthing to start")
	}
}

// awaitSync waits until rc.InSync reports both processes are in sync on
// the named folder, or fails the test on timeout.
func awaitSync(t *testing.T, folder string, ps ...*rc.Process) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		rc.AwaitSync(folder, ps...)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(120 * time.Second):
		t.Fatalf("timed out waiting for folder %q to sync", folder)
	}
}

// connectTo pushes a configuration update via the REST API that adds peer
// as a known device and shares the "default" folder with it.
func (i *instance) connectTo(t *testing.T, peer *instance) {
	t.Helper()

	cfg, err := i.process.GetConfig()
	if err != nil {
		t.Fatalf("instance %d: get config: %v", i.index, err)
	}

	cfg.Devices = append(cfg.Devices, config.DeviceConfiguration{
		DeviceID:  peer.deviceID,
		Name:      fmt.Sprintf("instance-%d", peer.index),
		Addresses: []string{peer.listenURL},
	})

	for idx := range cfg.Folders {
		if cfg.Folders[idx].ID == "default" {
			cfg.Folders[idx].Devices = append(cfg.Folders[idx].Devices,
				config.FolderDeviceConfiguration{DeviceID: peer.deviceID})
		}
	}

	if err := i.process.PostConfig(cfg); err != nil {
		t.Fatalf("instance %d: post config: %v", i.index, err)
	}
}

// readConfigXML parses a Syncthing config file from disk.
func readConfigXML(t *testing.T, path string, myID protocol.DeviceID) config.Configuration {
	t.Helper()
	fd, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()
	cfg, _, err := config.ReadXML(fd, myID)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

// writeConfigXML writes a Syncthing config file to disk, replacing any
// existing file at the path.
func writeConfigXML(t *testing.T, path string, cfg *config.Configuration) {
	t.Helper()
	fd, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.WriteXML(fd); err != nil {
		fd.Close()
		t.Fatal(err)
	}
	if err := fd.Close(); err != nil {
		t.Fatal(err)
	}
}
