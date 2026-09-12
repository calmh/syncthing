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

	configpb "github.com/syncthing/syncthing/internal/gen/syncthing/v2/config"
	"github.com/syncthing/syncthing/lib/protocol"
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
	if folder := parsed.GetFolders()[0]; !folder.HasPaused() {
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
	cfg := configpb.Configuration_builder{Version: new(int32(52))}.Build()
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
	cfg = configpb.Configuration_builder{Version: new(int32(0))}.Build()
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

	if cfg.HasGui() || cfg.HasOptions() {
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
	if folder.HasRescanIntervalS() {
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

// TestYAMLUnmarshalSize tests unmarshalling of sizes in each unit, and
// that an unset size means zero.
func TestYAMLUnmarshalSize(t *testing.T) {
	tests := []struct {
		unit      string
		yamlValue string
		want      float64
	}{
		{"percent", "10", 10},
		{"bytes", "1024", 1024},
		{"mib", "2", 2},
		{"gib", "5", 5},
	}
	for _, test := range tests {
		document := "folders:\n  - id: f1\n    minDiskFree:\n      " + test.unit + ": " + test.yamlValue + "\n"
		cfg, err := Unmarshal([]byte(document))
		if err != nil {
			t.Errorf("unmarshal %s: %v", test.unit, err)
			continue
		}
		size := cfg.GetFolders()[0].GetMinDiskFree()
		var got float64
		switch {
		case size.HasPercent():
			got = size.GetPercent()
		case size.HasBytes():
			got = size.GetBytes()
		case size.HasMib():
			got = size.GetMib()
		case size.HasGib():
			got = size.GetGib()
		default:
			t.Errorf("size for %s not set", test.unit)
			continue
		}
		if got != test.want {
			t.Errorf("size %s: got %v, want %v", test.unit, got, test.want)
		}
	}

	// An absent size is unset, meaning zero.
	cfg, err := Unmarshal([]byte("folders:\n  - id: f1\n"))
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if size := cfg.GetFolders()[0].GetMinDiskFree(); size != nil {
		t.Errorf("absent size should be unset, got %v", size)
	}
}

// TestYAMLUnmarshalValidation tests that validation rules are enforced at
// load time.
func TestYAMLUnmarshalValidation(t *testing.T) {
	valid := []string{
		"folders:\n  - id: f1\n    minDiskFree:\n      gib: 5\n",
		"folders:\n  - id: f1\n    minDiskFree:\n      percent: 10\n",
		// Empty device IDs are allowed, e.g. for device defaults.
		"devices:\n  - deviceId: \"\"\n",
	}
	for _, document := range valid {
		if _, err := Unmarshal([]byte(document)); err != nil {
			t.Errorf("valid document rejected: %v\n%s", err, document)
		}
	}

	invalid := []string{
		// Percent must be within 0-100.
		"folders:\n  - id: f1\n    minDiskFree:\n      percent: 150\n",
		// Sizes must not be negative.
		"folders:\n  - id: f1\n    minDiskFree:\n      bytes: -1\n",
		"folders:\n  - id: f1\n    minDiskFree:\n      mib: -1\n",
		"folders:\n  - id: f1\n    minDiskFree:\n      gib: -1\n",
		"options:\n  minHomeDiskFree:\n    percent: -1\n",
		// Device IDs must parse as device IDs (or be empty).
		"devices:\n  - deviceId: not-a-device-id\n",
		"folders:\n  - id: f1\n    devices:\n      - deviceId: abc\n",
	}
	for _, document := range invalid {
		if _, err := Unmarshal([]byte(document)); err == nil {
			t.Errorf("invalid document accepted:\n%s", document)
		}
	}
}

// TestYAMLUnmarshalDeviceIDValidation tests that device IDs are verified
// to parse with matching check digits, beyond the shape checked by the
// pattern rules.
func TestYAMLUnmarshalDeviceIDValidation(t *testing.T) {
	good := protocol.NewDeviceID([]byte("gooddevice")).String()

	// Corrupt one character, keeping the shape of the ID intact: the
	// check digits no longer match.
	corrupt := []byte(good)
	if corrupt[0] == 'A' {
		corrupt[0] = 'B'
	} else {
		corrupt[0] = 'A'
	}
	bad := string(corrupt)
	if _, err := protocol.DeviceIDFromString(bad); err == nil {
		t.Fatalf("test bug: corrupted device ID %q still parses", bad)
	}

	if _, err := Unmarshal([]byte("devices:\n  - deviceId: " + good + "\n")); err != nil {
		t.Errorf("valid device ID rejected: %v", err)
	}
	if _, err := Unmarshal([]byte("folders:\n  - id: f1\n    devices:\n      - deviceId: " + good + "\n")); err != nil {
		t.Errorf("valid device ID rejected: %v", err)
	}
	// Lenient forms that still parse: lowercase, and without dashes.
	if _, err := Unmarshal([]byte("devices:\n  - deviceId: " + strings.ToLower(good) + "\n")); err != nil {
		t.Errorf("lowercase device ID rejected: %v", err)
	}
	if _, err := Unmarshal([]byte("devices:\n  - deviceId: " + strings.ReplaceAll(good, "-", "") + "\n")); err != nil {
		t.Errorf("device ID without dashes rejected: %v", err)
	}

	for _, document := range []string{
		"devices:\n  - deviceId: " + bad + "\n",
		"folders:\n  - id: f1\n    devices:\n      - deviceId: " + bad + "\n",
	} {
		_, err := Unmarshal([]byte(document))
		if err == nil {
			t.Errorf("device ID with bad check digits accepted:\n%s", document)
			continue
		}
		if !strings.Contains(err.Error(), "invalid device ID") {
			t.Errorf("error should mention the device ID: %v", err)
		}
	}
}

// testConfiguration returns a configuration populated with a
// representative spread of fields and values.
func testConfiguration() *configpb.Configuration {
	return configpb.Configuration_builder{
		Version: new(int32(52)),
		Folders: []*configpb.FolderConfiguration{configpb.FolderConfiguration_builder{
			Id:               new("test1"),
			Label:            new("Test One"),
			Path:             new("/srv/sync/test1"),
			Type:             new(configpb.FolderType_FOLDER_TYPE_SEND_ONLY),
			RescanIntervalS:  new(int32(7200)),
			FsWatcherEnabled: new(false), // explicit zero value
			FsWatcherDelayS:  new(12.5),
			MinDiskFree:      configpb.Size_builder{Gib: new(5.0)}.Build(),
			MaxConflicts:     new(int32(0)), // explicit zero value
			Paused:           new(false),
			CopyRangeMethod:  new(configpb.CopyRangeMethod_COPY_RANGE_METHOD_IOCTL),
			BlockPullOrder:   new(configpb.BlockPullOrder_BLOCK_PULL_ORDER_IN_ORDER),
			Order:            new(configpb.PullOrder_PULL_ORDER_NEWEST_FIRST),
			Versioning: configpb.VersioningConfiguration_builder{
				Type:             new("simple"),
				Params:           map[string]string{"keep": "5", "cleanoutDays": "30"},
				CleanupIntervalS: new(int32(3600)),
				FsPath:           new("/srv/sync/test1/.stversions"),
				FsType:           new(configpb.FilesystemType_FILESYSTEM_TYPE_FAKE),
			}.Build(),
			XattrFilter: configpb.XattrFilter_builder{
				Entries: []*configpb.XattrFilterEntry{
					configpb.XattrFilterEntry_builder{Match: new("user.sync.*"), Permit: new(true)}.Build(),
					configpb.XattrFilterEntry_builder{Match: new("*"), Permit: new(false)}.Build(),
				},
				MaxSingleEntrySize: new(int32(2048)),
				MaxTotalSize:       new(int32(8192)),
			}.Build(),
			Devices: []*configpb.FolderDeviceConfiguration{
				configpb.FolderDeviceConfiguration_builder{
					DeviceId:           new("AIR6LPZ-7ENK4MY-4VVODLE-VDGP6LY-Q3GSA2S-7SJMDTM-DLJ4BGV-K3ZP4QU"),
					IntroducedBy:       new("BYR4TZD-LFDLPXV-KLS4HPL-XK25UQG-BUF6TCD-6FLFXCO-YV3FQXZ-DQ7SPQU"),
					EncryptionPassword: new("hunter2"),
				}.Build(),
				configpb.FolderDeviceConfiguration_builder{
					DeviceId: new("ZK6FOFT-TXHAKOT-DNJRW3B-7EDSH2F-C4QIYYI-ZAX2UXW-I3HZLGY-YGLAZQU"),
				}.Build(),
			},
		}.Build()},
		Devices: []*configpb.DeviceConfiguration{
			configpb.DeviceConfiguration_builder{
				DeviceId:        new("AIR6LPZ-7ENK4MY-4VVODLE-VDGP6LY-Q3GSA2S-7SJMDTM-DLJ4BGV-K3ZP4QU"),
				Name:            new("Alpha"),
				Addresses:       []string{"dynamic", "tcp://192.0.2.1:22000"},
				Compression:     new(configpb.Compression_COMPRESSION_ALWAYS),
				Introducer:      new(true),
				Paused:          new(false),
				MaxSendKbps:     new(int32(1000)),
				MaxRecvKbps:     new(int32(2000)),
				NumConnections:  new(int32(2)),
				AllowedNetworks: []string{"192.168.0.0/16"},
				IgnoredFolders: []*configpb.ObservedFolder{configpb.ObservedFolder_builder{
					Time:  timestamppb.New(time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)),
					Id:    new("ignored-folder"),
					Label: new("Ignored Folder"),
				}.Build()},
			}.Build(),
			configpb.DeviceConfiguration_builder{
				DeviceId: new("ZK6FOFT-TXHAKOT-DNJRW3B-7EDSH2F-C4QIYYI-ZAX2UXW-I3HZLGY-YGLAZQU"),
				Name:     new("Beta"),
			}.Build(),
		},
		Gui: configpb.GUIConfiguration_builder{
			Enabled:  new(true),
			Address:  new("127.0.0.1:8384"),
			User:     new("admin"),
			Password: new("$2a$10$X5bFBmJHrHrrIYpuLvBvJ.uUGPmuffeYdWtMskYiCVnrUuLLzJBae"),
			AuthMode: new(configpb.AuthMode_AUTH_MODE_LDAP),
			UseTls:   new(true),
			ApiKey:   new("kO8nJgP7tY2wZq4x"),
			Theme:    new("dark"),
		}.Build(),
		Ldap: configpb.LDAPConfiguration_builder{
			Address:            new("ldap.example.com:389"),
			BindDn:             new("cn=admin,dc=example,dc=com"),
			Transport:          new(configpb.LDAPTransport_LDAP_TRANSPORT_START_TLS),
			InsecureSkipVerify: new(true),
			SearchBaseDn:       new("dc=example,dc=com"),
			SearchFilter:       new("(uid=%s)"),
		}.Build(),
		Options: configpb.OptionsConfiguration_builder{
			ListenAddresses:          []string{"default"},
			GlobalAnnounceServers:    []string{"default"},
			LocalAnnouncePort:        new(int32(21027)),
			MaxSendKbps:              new(int32(5000)),
			ReconnectionIntervalS:    new(int32(30)),
			UrAccepted:               new(int32(-1)), // negative value
			KeepTemporariesH:         new(int32(12)),
			MinHomeDiskFree:          configpb.Size_builder{Percent: new(1.0)}.Build(),
			AlwaysLocalNets:          []string{"10.0.0.0/8"},
			UnackedNotificationIds:   []string{"authenticationUserAndPassword"},
			SetLowPriority:           new(false), // explicit zero value
			FeatureFlags:             []string{"caves"},
			ConnectionPriorityTcpLan: new(int32(15)),
			StunKeepaliveMinS:        new(int32(25)),
		}.Build(),
		RemoteIgnoredDevices: []*configpb.ObservedDevice{configpb.ObservedDevice_builder{
			Time:     timestamppb.New(time.Date(2026, 9, 1, 8, 30, 0, 0, time.UTC)),
			DeviceId: new("E5OA2FV-LO2YUML-3SDRQ4C-JPJF3DQ-JSWF3LV-CKSPR4T-KKVMKWS-JVVBQAE"),
			Name:     new("Gamma"),
			Address:  new("192.0.2.7:22000"),
		}.Build()},
		Defaults: configpb.Defaults_builder{
			Folder: configpb.FolderConfiguration_builder{
				RescanIntervalS: new(int32(10800)),
			}.Build(),
			Device: configpb.DeviceConfiguration_builder{
				Compression: new(configpb.Compression_COMPRESSION_NEVER),
			}.Build(),
			Ignores: configpb.Ignores_builder{
				Lines: []string{"!qux", "baz/*"},
			}.Build(),
		}.Build(),
	}.Build()
}
