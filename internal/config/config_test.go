// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	configpb "github.com/syncthing/syncthing/internal/gen/config"
)

// TestYAMLRoundTrip tests that a configuration survives being marshalled
// to YAML and unmarshalled back, including the distinction between unset
// fields and fields explicitly set to their zero value.
func TestYAMLRoundTrip(t *testing.T) {
	orig := testConfiguration()

	mo := protoyaml.MarshalOptions{
		Indent: 2,
	}
	data, err := mo.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	os.WriteFile("_testdata/roundtrip.yaml", data, 0o644)

	var parsed configpb.Configuration
	if err := protoyaml.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal: %v\n%s", err, data)
	}

	if !proto.Equal(orig, &parsed) {
		t.Errorf("configuration did not survive round trip:\ngot:  %v\nwant: %v", &parsed, orig)
	}

	// Fields explicitly set to their zero value must remain set.
	if folder := parsed.GetFolders()[0]; folder.Paused == nil {
		t.Errorf("explicitly set paused=false lost its presence in round trip")
	}

	// Enum values are marshalled as their value names, timestamps as
	// quoted RFC 3339 strings.
	for _, expected := range []string{
		"type: FOLDER_TYPE_SEND_ONLY",
		`time: "2026-08-01T12:00:00Z"`,
	} {
		if !strings.Contains(string(data), expected) {
			t.Errorf("expected %q in marshalled YAML:\n%s", expected, data)
		}
	}
}

// TestYAMLMarshalPresence tests that marshalling emits set fields and
// omits unset ones, so that defaults are not materialised into the file.
func TestYAMLMarshalPresence(t *testing.T) {
	cfg := &configpb.Configuration{Version: proto.Int32(52)}
	data, err := protoyaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(data), "version: 52") {
		t.Errorf("set field missing from YAML:\n%s", data)
	}
	if strings.Contains(string(data), "gui") {
		t.Errorf("unset field present in YAML:\n%s", data)
	}

	// An explicitly set zero value is emitted.
	cfg = &configpb.Configuration{Version: proto.Int32(0)}
	data, err = protoyaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(data), "version: 0") {
		t.Errorf("explicitly set zero value missing from YAML:\n%s", data)
	}

	// An unset field is not.
	cfg = &configpb.Configuration{}
	data, err = protoyaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(data), "version") {
		t.Errorf("unset field present in YAML:\n%s", data)
	}
}

// TestYAMLUnmarshalDefaults tests that fields absent from the YAML remain
// unset, and that getters return the declared default values.
func TestYAMLUnmarshalDefaults(t *testing.T) {
	const document = `
version: 52
folders:
  - id: minimal
`
	var cfg configpb.Configuration
	if err := protoyaml.Unmarshal([]byte(document), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if cfg.Gui != nil || cfg.Options != nil {
		t.Errorf("unset messages should remain unset")
	}

	// Getters on unset messages and fields return the declared defaults.
	// Check a representative spread of default types.
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"gui.enabled", cfg.GetGui().GetEnabled(), true},
		{"gui.address", cfg.GetGui().GetAddress(), "127.0.0.1:8384"},
		{"gui.theme", cfg.GetGui().GetTheme(), "default"},
		{"gui.sessionCookieDurationS", cfg.GetGui().GetSessionCookieDurationS(), int32(604800)},
		{"gui.sessionCookiePath", cfg.GetGui().GetSessionCookiePath(), "/"},
		{"options.localAnnouncePort", cfg.GetOptions().GetLocalAnnouncePort(), int32(21027)},
		{"options.localAnnounceMCAddr", cfg.GetOptions().GetLocalAnnounceMcAddr(), "[ff12::8384]:21027"},
		{"options.reconnectionIntervalS", cfg.GetOptions().GetReconnectionIntervalS(), int32(20)},
		{"options.urURL", cfg.GetOptions().GetUrUrl(), "https://data.syncthing.net/newdata"},
		{"options.releasesURL", cfg.GetOptions().GetReleasesUrl(), "https://upgrades.syncthing.net/meta.json"},
		{"options.crashReportingURL", cfg.GetOptions().GetCrUrl(), "https://crash.syncthing.net/newcrash"},
		{"options.connectionPriorityRelay", cfg.GetOptions().GetConnectionPriorityRelay(), int32(50)},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("%s: got %v, want default %v", test.name, test.got, test.want)
		}
	}

	folder := cfg.GetFolders()[0]
	if folder.RescanIntervalS != nil {
		t.Errorf("rescanIntervalS should remain unset")
	}
	if got := folder.GetRescanIntervalS(); got != 3600 {
		t.Errorf("folder.rescanIntervalS: got %d, want default 3600", got)
	}
	if got := folder.GetFilesystemType(); got != configpb.FilesystemType_FILESYSTEM_TYPE_BASIC {
		t.Errorf("folder.filesystemType: got %v, want default %v", got, configpb.FilesystemType_FILESYSTEM_TYPE_BASIC)
	}
	if got := folder.GetCopyRangeMethod(); got != configpb.CopyRangeMethod_COPY_RANGE_METHOD_STANDARD {
		t.Errorf("folder.copyRangeMethod: got %v, want default %v", got, configpb.CopyRangeMethod_COPY_RANGE_METHOD_STANDARD)
	}
	if !folder.GetAutoNormalize() || !folder.GetBlockIndexing() {
		t.Errorf("folder.autoNormalize and folder.blockIndexing should default to true")
	}
	if got := folder.GetFsWatcherDelayS(); got != 10 {
		t.Errorf("folder.fsWatcherDelayS: got %v, want default 10", got)
	}
	if got := folder.GetVersioning().GetCleanupIntervalS(); got != 3600 {
		t.Errorf("folder.versioning.cleanupIntervalS: got %d, want default 3600", got)
	}
	if got := folder.GetXattrFilter().GetMaxSingleEntrySize(); got != 1024 {
		t.Errorf("folder.xattrFilter.maxSingleEntrySize: got %d, want default 1024", got)
	}
	if got := folder.GetXattrFilter().GetMaxTotalSize(); got != 4096 {
		t.Errorf("folder.xattrFilter.maxTotalSize: got %d, want default 4096", got)
	}

	// Marshalling the result does not materialise the defaults into YAML.
	data, err := protoyaml.Marshal(&cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, unexpected := range []string{"gui", "options", "rescanIntervalS"} {
		if strings.Contains(string(data), unexpected) {
			t.Errorf("default value for unset field %q materialised into YAML:\n%s", unexpected, data)
		}
	}
}

// TestYAMLUnmarshalUnknownField tests that unknown fields in the YAML are
// an error, so that typos don't get silently ignored.
func TestYAMLUnmarshalUnknownField(t *testing.T) {
	err := protoyaml.Unmarshal([]byte("version: 52\nnoSuchField: 1\n"), &configpb.Configuration{})
	if err == nil {
		t.Fatalf("expected error on unknown field")
	}
	if !strings.Contains(err.Error(), "noSuchField") {
		t.Errorf("error should mention the unknown field: %v", err)
	}
}

// testConfiguration returns a configuration populated with a
// representative spread of fields and values.
func testConfiguration() *configpb.Configuration {
	return &configpb.Configuration{
		Version: proto.Int32(52),
		Folders: []*configpb.FolderConfiguration{{
			Id:               proto.String("test1"),
			Label:            proto.String("Test One"),
			Path:             proto.String("/srv/sync/test1"),
			Type:             configpb.FolderType_FOLDER_TYPE_SEND_ONLY.Enum(),
			RescanIntervalS:  proto.Int32(7200),
			FsWatcherEnabled: proto.Bool(false), // explicit zero value
			FsWatcherDelayS:  proto.Float64(12.5),
			MinDiskFree:      &configpb.Size{Value: proto.Float64(5), Unit: proto.String("GiB")},
			MaxConflicts:     proto.Int32(0), // explicit zero value
			Paused:           proto.Bool(false),
			CopyRangeMethod:  configpb.CopyRangeMethod_COPY_RANGE_METHOD_IOCTL.Enum(),
			BlockPullOrder:   configpb.BlockPullOrder_BLOCK_PULL_ORDER_IN_ORDER.Enum(),
			Order:            configpb.PullOrder_PULL_ORDER_NEWEST_FIRST.Enum(),
			Versioning: &configpb.VersioningConfiguration{
				Type:             proto.String("simple"),
				Params:           map[string]string{"keep": "5", "cleanoutDays": "30"},
				CleanupIntervalS: proto.Int32(3600),
				FsPath:           proto.String("/srv/sync/test1/.stversions"),
				FsType:           configpb.FilesystemType_FILESYSTEM_TYPE_FAKE.Enum(),
			},
			XattrFilter: &configpb.XattrFilter{
				Entries: []*configpb.XattrFilterEntry{
					{Match: proto.String("user.sync.*"), Permit: proto.Bool(true)},
					{Match: proto.String("*"), Permit: proto.Bool(false)},
				},
				MaxSingleEntrySize: proto.Int32(2048),
				MaxTotalSize:       proto.Int32(8192),
			},
			Devices: []*configpb.FolderDeviceConfiguration{
				{
					DeviceId:           proto.String("AIR6LPZ-7ENK4MY-4VVODLE-VDGP6LY-Q3GSA2S-7SJMDTM-DLJ4BGV-K3ZP4QU"),
					IntroducedBy:       proto.String("BYR4TZD-LFDLPXV-KLS4HPL-XK25UQG-BUF6TCD-6FLFXCO-YV3FQXZ-DQ7SPQU"),
					EncryptionPassword: proto.String("hunter2"),
				},
				{
					DeviceId: proto.String("ZK6FOFT-TXHAKOT-DNJRW3B-7EDSH2F-C4QIYYI-ZAX2UXW-I3HZLGY-YGLAZQU"),
				},
			},
		}},
		Devices: []*configpb.DeviceConfiguration{
			{
				DeviceId:        proto.String("AIR6LPZ-7ENK4MY-4VVODLE-VDGP6LY-Q3GSA2S-7SJMDTM-DLJ4BGV-K3ZP4QU"),
				Name:            proto.String("Alpha"),
				Addresses:       []string{"dynamic", "tcp://192.0.2.1:22000"},
				Compression:     configpb.Compression_COMPRESSION_ALWAYS.Enum(),
				Introducer:      proto.Bool(true),
				Paused:          proto.Bool(false),
				MaxSendKbps:     proto.Int32(1000),
				MaxRecvKbps:     proto.Int32(2000),
				NumConnections:  proto.Int32(2),
				AllowedNetworks: []string{"192.168.0.0/16"},
				IgnoredFolders: []*configpb.ObservedFolder{{
					Time:  timestamppb.New(time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)),
					Id:    proto.String("ignored-folder"),
					Label: proto.String("Ignored Folder"),
				}},
			},
			{
				DeviceId: proto.String("ZK6FOFT-TXHAKOT-DNJRW3B-7EDSH2F-C4QIYYI-ZAX2UXW-I3HZLGY-YGLAZQU"),
				Name:     proto.String("Beta"),
			},
		},
		Gui: &configpb.GUIConfiguration{
			Enabled:  proto.Bool(true),
			Address:  proto.String("127.0.0.1:8384"),
			User:     proto.String("admin"),
			Password: proto.String("$2a$10$X5bFBmJHrHrrIYpuLvBvJ.uUGPmuffeYdWtMskYiCVnrUuLLzJBae"),
			AuthMode: configpb.AuthMode_AUTH_MODE_LDAP.Enum(),
			UseTls:   proto.Bool(true),
			ApiKey:   proto.String("kO8nJgP7tY2wZq4x"),
			Theme:    proto.String("dark"),
		},
		Ldap: &configpb.LDAPConfiguration{
			Address:            proto.String("ldap.example.com:389"),
			BindDn:             proto.String("cn=admin,dc=example,dc=com"),
			Transport:          configpb.LDAPTransport_LDAP_TRANSPORT_START_TLS.Enum(),
			InsecureSkipVerify: proto.Bool(true),
			SearchBaseDn:       proto.String("dc=example,dc=com"),
			SearchFilter:       proto.String("(uid=%s)"),
		},
		Options: &configpb.OptionsConfiguration{
			ListenAddresses:          []string{"default"},
			GlobalAnnounceServers:    []string{"default"},
			LocalAnnouncePort:        proto.Int32(21027),
			MaxSendKbps:              proto.Int32(5000),
			ReconnectionIntervalS:    proto.Int32(30),
			UrAccepted:               proto.Int32(-1), // negative value
			KeepTemporariesH:         proto.Int32(12),
			MinHomeDiskFree:          &configpb.Size{Value: proto.Float64(1), Unit: proto.String("%")},
			AlwaysLocalNets:          []string{"10.0.0.0/8"},
			UnackedNotificationIds:   []string{"authenticationUserAndPassword"},
			SetLowPriority:           proto.Bool(false), // explicit zero value
			FeatureFlags:             []string{"caves"},
			ConnectionPriorityTcpLan: proto.Int32(15),
			StunKeepaliveMinS:        proto.Int32(25),
		},
		RemoteIgnoredDevices: []*configpb.ObservedDevice{{
			Time:     timestamppb.New(time.Date(2026, 9, 1, 8, 30, 0, 0, time.UTC)),
			DeviceId: proto.String("E5OA2FV-LO2YUML-3SDRQ4C-JPJF3DQ-JSWF3LV-CKSPR4T-KKVMKWS-JVVBQAE"),
			Name:     proto.String("Gamma"),
			Address:  proto.String("192.0.2.7:22000"),
		}},
		Defaults: &configpb.Defaults{
			Folder: &configpb.FolderConfiguration{
				RescanIntervalS: proto.Int32(10800),
			},
			Device: &configpb.DeviceConfiguration{
				Compression: configpb.Compression_COMPRESSION_NEVER.Enum(),
			},
			Ignores: &configpb.Ignores{
				Lines: []string{"!qux", "baz/*"},
			},
		},
	}
}
