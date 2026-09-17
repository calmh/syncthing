// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

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

// TestConfigServiceGetConfigurationMaterialized tests the
// materializeDefaults flag: fields that have their defaults declared by
// getters, and are unset in the configuration, are filled in.
func TestConfigServiceGetConfigurationMaterialized(t *testing.T) {
	id := protocol.NewDeviceID([]byte("configmaterialize"))
	// A sparse legacy configuration, where the options and device
	// addresses are unset rather than materialised.
	cfg := config.Configuration{
		Version: config.CurrentVersion,
		Devices: []config.DeviceConfiguration{{DeviceID: id}},
	}
	svc := &configService{fakeConfigSource{cfg: cfg}}

	path, handler := syncthingv2connect.NewConfigurationServiceHandler(svc)
	mux := http.NewServeMux()
	mux.Handle(path, handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	client := syncthingv2connect.NewConfigurationServiceClient(server.Client(), server.URL)

	// Without the flag, the getter-defaulted fields are unset.
	resp, err := client.GetConfiguration(context.Background(), connect.NewRequest(&syncthingv2.GetConfigRequest{}))
	if err != nil {
		t.Fatalf("GetConfiguration: %v", err)
	}
	if opts := resp.Msg.GetConfiguration().GetOptions(); len(opts.GetListenAddresses()) != 0 {
		t.Errorf("listenAddresses: got %v, want unset", opts.GetListenAddresses())
	}

	// With the flag, they are filled in.
	resp, err = client.GetConfiguration(context.Background(), connect.NewRequest(
		syncthingv2.GetConfigRequest_builder{MaterializeDefaults: new(true)}.Build(),
	))
	if err != nil {
		t.Fatalf("GetConfiguration: %v", err)
	}
	got := resp.Msg.GetConfiguration()
	if addrs := got.GetOptions().GetListenAddresses(); len(addrs) != 1 || addrs[0] != "default" {
		t.Errorf("listenAddresses: got %v, want [default]", addrs)
	}
	if addrs := got.GetDevices()[0].GetAddresses(); len(addrs) != 1 || addrs[0] != "dynamic" {
		t.Errorf("device addresses: got %v, want [dynamic]", addrs)
	}
	// The legacy zero size means the check is disabled, which converts
	// to an explicit zero rather than the default.
	if size := got.GetOptions().GetMinHomeDiskFree(); size == nil || !size.HasPercent() || size.GetPercent() != 0 {
		t.Errorf("minHomeDiskFree: got %v, want explicit zero", size)
	}
}

// TestConfigServiceGetConfigurationHTTPGet tests that the RPC is reachable
// with a plain HTTP GET request, per the NO_SIDE_EFFECTS idempotency
// annotation. The Connect protocol requires the response codec and the
// request message as query parameters on GET requests; the defaults
// wrapper fills them in so that a bare GET works and gets JSON.
func TestConfigServiceGetConfigurationHTTPGet(t *testing.T) {
	id := protocol.NewDeviceID([]byte("configservice-get"))
	svc := &configService{fakeConfigSource{cfg: config.New(id)}}

	// Mount as in the main API mux, including the GET defaults wrapper.
	path, handler := syncthingv2connect.NewConfigurationServiceHandler(svc)
	mux := http.NewServeMux()
	mux.Handle(path, withConnectGetDefaults(handler))
	server := httptest.NewServer(mux)
	defer server.Close()

	check := func(url, wantContentType string, unmarshal func([]byte) (*syncthingv2.GetConfigResponse, error)) {
		t.Helper()
		resp, err := server.Client().Get(url)
		if err != nil {
			t.Fatalf("GET: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GET: unexpected status %v", resp.Status)
		}
		if got := resp.Header.Get("Content-Type"); got != wantContentType {
			t.Errorf("content type: got %q, want %q", got, wantContentType)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		msg, err := unmarshal(body)
		if err != nil {
			t.Fatalf("unmarshal response: %v\n%s", err, body)
		}
		if v := msg.GetConfiguration().GetVersion(); v != int32(config.CurrentVersion) {
			t.Errorf("version: got %d, want %d", v, config.CurrentVersion)
		}
	}

	fromJSON := func(b []byte) (*syncthingv2.GetConfigResponse, error) {
		var msg syncthingv2.GetConfigResponse
		err := protojson.Unmarshal(b, &msg)
		return &msg, err
	}
	fromProto := func(b []byte) (*syncthingv2.GetConfigResponse, error) {
		var msg syncthingv2.GetConfigResponse
		err := proto.Unmarshal(b, &msg)
		return &msg, err
	}

	// A bare GET gets the defaults: JSON response of the empty request.
	check(server.URL+path+"GetConfiguration", "application/json", fromJSON)

	// Explicit parameters are respected, in any codec.
	check(server.URL+path+"GetConfiguration?encoding=json&message=%7B%7D", "application/json", fromJSON)
	check(server.URL+path+"GetConfiguration?encoding=proto&base64=1&message=", "application/proto", fromProto)
}
