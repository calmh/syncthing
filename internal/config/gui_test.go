// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"testing"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
)

func guiCfg() GUIConfiguration {
	return GUIConfiguration{syncthingv2.GUIConfiguration_builder{}.Build()}
}

func TestGUIEnvOverrides(t *testing.T) {
	// Without override, the schema defaults apply.
	gui := guiCfg()
	if gui.IsOverridden() {
		t.Error("should not be overridden")
	}
	if gui.Address() != "127.0.0.1:8384" {
		t.Errorf("address: got %q, want default", gui.Address())
	}
	if gui.Network() != "tcp" {
		t.Errorf("network: got %q, want tcp", gui.Network())
	}
	if gui.UseTLS() {
		t.Error("TLS: got true, want default false")
	}
	if gui.URL() != "http://127.0.0.1:8384/" {
		t.Errorf("url: got %q", gui.URL())
	}

	// An address override replaces the address, network and scheme.
	t.Setenv("STGUIADDRESS", "0.0.0.0:9999")
	if !gui.IsOverridden() {
		t.Error("should be overridden")
	}
	if gui.Address() != "0.0.0.0:9999" {
		t.Errorf("address: got %q", gui.Address())
	}
	if gui.Network() != "tcp" {
		t.Errorf("network: got %q", gui.Network())
	}
	if gui.URL() != "http://127.0.0.1:9999/" {
		t.Errorf("url: got %q, want localhost substituted", gui.URL())
	}

	// An https override enables TLS.
	t.Setenv("STGUIADDRESS", "https://0.0.0.0:9999")
	if !gui.UseTLS() {
		t.Error("TLS should be enabled by https override")
	}
	if gui.Address() != "0.0.0.0:9999" {
		t.Errorf("address: got %q, want scheme stripped", gui.Address())
	}

	// A unix socket override selects the unix network.
	t.Setenv("STGUIADDRESS", "unix:///tmp/syncthing.sock")
	if gui.Network() != "unix" {
		t.Errorf("network: got %q, want unix", gui.Network())
	}
	if gui.Address() != "/tmp/syncthing.sock" {
		t.Errorf("address: got %q, want socket path", gui.Address())
	}
	if gui.URL() != "unix:///tmp/syncthing.sock" {
		t.Errorf("url: got %q", gui.URL())
	}

	// A stored UNIX socket path is recognised without override.
	gui = GUIConfiguration{syncthingv2.GUIConfiguration_builder{
		Address: new("/tmp/syncthing.sock"),
	}.Build()}
	if gui.Network() != "unix" {
		t.Errorf("network: got %q, want unix", gui.Network())
	}
}

func TestGUIAuth(t *testing.T) {
	// Without user and password, auth is disabled.
	gui := guiCfg()
	if gui.IsAuthEnabled() {
		t.Error("auth should be disabled")
	}

	// Setting a plaintext password stores the hash, and the hash
	// compares correctly.
	gui = GUIConfiguration{syncthingv2.GUIConfiguration_builder{
		User: new("admin"),
	}.Build()}
	if err := gui.SetPassword("hunter2"); err != nil {
		t.Fatalf("set password: %v", err)
	}
	if gui.GetPassword() == "hunter2" {
		t.Error("password should be hashed")
	}
	if !gui.IsAuthEnabled() {
		t.Error("auth should be enabled")
	}
	if err := gui.CompareHashedPassword("hunter2"); err != nil {
		t.Errorf("correct password: %v", err)
	}
	if err := gui.CompareHashedPassword("wrong"); err == nil {
		t.Error("wrong password should not compare")
	}

	// A pre-hashed password is stored as-is.
	hashed := gui.GetPassword()
	if err := gui.SetPassword(hashed); err != nil {
		t.Fatalf("set hashed password: %v", err)
	}
	if gui.GetPassword() != hashed {
		t.Error("hash should be kept")
	}

	// LDAP mode enables auth without user and password.
	gui = GUIConfiguration{syncthingv2.GUIConfiguration_builder{
		AuthMode: new(syncthingv2.AuthMode_AUTH_MODE_LDAP),
	}.Build()}
	if !gui.IsAuthEnabled() {
		t.Error("auth should be enabled in LDAP mode")
	}

	// The API key is valid from the config or the environment, but never
	// empty.
	gui = GUIConfiguration{syncthingv2.GUIConfiguration_builder{
		ApiKey: new("kO8nJgP7tY2wZq4x"),
	}.Build()}
	if !gui.IsValidAPIKey("kO8nJgP7tY2wZq4x") {
		t.Error("stored API key should be valid")
	}
	if gui.IsValidAPIKey("wrong") {
		t.Error("wrong API key should not be valid")
	}
	if gui.IsValidAPIKey("") {
		t.Error("empty API key should not be valid")
	}
	t.Setenv("STGUIAPIKEY", "env-key")
	if !gui.IsValidAPIKey("env-key") {
		t.Error("environment API key should be valid")
	}
	if !gui.IsValidAPIKey("kO8nJgP7tY2wZq4x") {
		t.Error("stored API key should still be valid")
	}
}
