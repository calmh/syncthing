// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/internal/gen/syncthing/v2/syncthingv2connect"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/protocol"
)

// fakeConfigSource is a configSource returning a fixed configuration.
type fakeConfigSource struct {
	cfg config.Configuration
}

func (s fakeConfigSource) RawCopy() config.Configuration { return s.cfg }

func TestConfigServiceGetConfiguration(t *testing.T) {
	id := protocol.NewDeviceID([]byte("configservice"))
	svc := &configService{fakeConfigSource{cfg: config.New(id)}}

	// Mount the service at its default path on the root, the same way as
	// the main API mux.
	path, handler := syncthingv2connect.NewConfigurationServiceHandler(svc)
	mux := http.NewServeMux()
	mux.Handle(path, handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	client := syncthingv2connect.NewConfigurationServiceClient(server.Client(), server.URL)
	resp, err := client.GetConfiguration(context.Background(), connect.NewRequest(&syncthingv2.GetConfigRequest{}))
	if err != nil {
		t.Fatalf("GetConfiguration: %v", err)
	}

	got := resp.Msg.GetConfiguration()
	if v := got.GetVersion(); v != int32(config.CurrentVersion) {
		t.Errorf("version: got %d, want %d", v, config.CurrentVersion)
	}
	if devices := got.GetDevices(); len(devices) != 1 || devices[0].GetDeviceId() != id.String() {
		t.Errorf("devices: got %v, want single device %s", devices, id)
	}
	if !got.HasGui() || !got.HasOptions() {
		t.Errorf("gui and options should be set")
	}
	if v := got.GetOptions().GetReconnectionIntervalS(); v != 20 {
		t.Errorf("options.reconnectionIntervalS: got %d, want default 20", v)
	}
}
