// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"reflect"
	"testing"
	"time"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/proto"

	configpb "github.com/syncthing/syncthing/internal/gen/config"
	config "github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/protocol"
)

func TestFromLegacy(t *testing.T) {
	id := protocol.NewDeviceID([]byte("fromlegacy"))
	otherID := protocol.NewDeviceID([]byte("otherdevice"))
	ts := time.Date(2026, 9, 12, 12, 0, 0, 0, time.UTC)

	legacy := config.Configuration{
		Version: 52,
		Folders: []config.FolderConfiguration{{
			ID:              "f1",
			Label:           "Folder One",
			Path:            "/srv/sync/f1",
			FilesystemType:  config.FilesystemTypeFake,
			Type:            config.FolderTypeReceiveOnly,
			RescanIntervalS: 0, // explicit zero: no periodic rescan
			Devices: []config.FolderDeviceConfiguration{{
				DeviceID:           otherID,
				IntroducedBy:       id,
				EncryptionPassword: "hunter2",
			}},
			MinDiskFree: config.Size{Value: 2, Unit: "GiB"},
			Versioning: config.VersioningConfiguration{
				Type:             "simple",
				Params:           map[string]string{"keep": "3"},
				CleanupIntervalS: 0, // explicit zero: no cleanup
				FSType:           config.FilesystemTypeFake,
			},
			MaxConflicts:    0, // explicit zero
			Copiers:         2,
			Hashers:         4,
			Order:           config.PullOrderOldestFirst,
			BlockPullOrder:  config.BlockPullOrderInOrder,
			CopyRangeMethod: config.CopyRangeMethodSendFile,
			XattrFilter: config.XattrFilter{
				Entries:            []config.XattrFilterEntry{{Match: "user.*", Permit: true}},
				MaxSingleEntrySize: 512,
				MaxTotalSize:       1024,
			},
		}},
		Devices: []config.DeviceConfiguration{{
			DeviceID:          otherID,
			Name:              "Other",
			Addresses:         []string{"dynamic"},
			Compression:       config.CompressionAlways,
			MaxSendKbps:       1000,
			IgnoredFolders:    []config.ObservedFolder{{Time: ts, ID: "ig1", Label: "Ignored"}},
			RemoteGUIPort:     1234,
			RawNumConnections: 7,
			Untrusted:         true,
		}},
		GUI: config.GUIConfiguration{
			RawAddress: "127.0.0.1:9000",
			AuthMode:   config.AuthModeLDAP,
			RawUseTLS:  true,
		},
		LDAP: config.LDAPConfiguration{
			Address:   "ldap.example.com:389",
			Transport: config.LDAPTransportStartTLS,
		},
		Options: config.OptionsConfiguration{
			RawListenAddresses: []string{"default"},
			MaxSendKbps:        100,
			URAccepted:         -1,
			MinHomeDiskFree:    config.Size{Value: 1, Unit: "%"},
		},
		IgnoredDevices: []config.ObservedDevice{{Time: ts, ID: id, Name: "Friend", Address: "192.0.2.1:22000"}},
		Defaults: config.Defaults{
			Folder:  config.FolderConfiguration{RescanIntervalS: 10800},
			Device:  config.DeviceConfiguration{Compression: config.CompressionNever},
			Ignores: config.Ignores{Lines: []string{"*.tmp"}},
		},
	}

	got := FromLegacy(legacy)

	// Explicit zeros must remain set, not fall back to defaults.
	folder := got.GetFolders()[0]
	if !folder.HasRescanIntervalS() || folder.GetRescanIntervalS() != 0 {
		t.Errorf("explicit rescanIntervalS 0 not preserved: %v", folder.GetRescanIntervalS())
	} else if !folder.GetVersioning().HasCleanupIntervalS() || folder.GetVersioning().GetCleanupIntervalS() != 0 {
		t.Errorf("explicit cleanupIntervalS 0 not preserved: %v", folder.GetVersioning().GetCleanupIntervalS())
	} else if !folder.HasMaxConflicts() || folder.GetMaxConflicts() != 0 {
		t.Errorf("explicit maxConflicts 0 not preserved: %v", folder.GetMaxConflicts())
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"version", got.GetVersion(), int32(52)},
		{"folder.id", got.GetFolders()[0].GetId(), "f1"},
		{"folder.filesystemType", got.GetFolders()[0].GetFilesystemType(), configpb.FilesystemType_FILESYSTEM_TYPE_FAKE},
		{"folder.type", got.GetFolders()[0].GetType(), configpb.FolderType_FOLDER_TYPE_RECEIVE_ONLY},
		{"folder.minDiskFree.gib", got.GetFolders()[0].GetMinDiskFree().GetGib(), 2.0},
		{"folder.order", got.GetFolders()[0].GetOrder(), configpb.PullOrder_PULL_ORDER_OLDEST_FIRST},
		{"folder.blockPullOrder", got.GetFolders()[0].GetBlockPullOrder(), configpb.BlockPullOrder_BLOCK_PULL_ORDER_IN_ORDER},
		{"folder.copyRangeMethod", got.GetFolders()[0].GetCopyRangeMethod(), configpb.CopyRangeMethod_COPY_RANGE_METHOD_SEND_FILE},
		{"folder.versioning.type", got.GetFolders()[0].GetVersioning().GetType(), "simple"},
		{"folder.versioning.params", got.GetFolders()[0].GetVersioning().GetParams()["keep"], "3"},
		{"folder.versioning.fsType", got.GetFolders()[0].GetVersioning().GetFsType(), configpb.FilesystemType_FILESYSTEM_TYPE_FAKE},
		{"folder.device.deviceID", got.GetFolders()[0].GetDevices()[0].GetDeviceId(), otherID},
		{"folder.device.introducedBy", got.GetFolders()[0].GetDevices()[0].GetIntroducedBy(), id},
		{"folder.device.encryptionPassword", got.GetFolders()[0].GetDevices()[0].GetEncryptionPassword(), "hunter2"},
		{"folder.xattrFilter.entry", got.GetFolders()[0].GetXattrFilter().GetEntries()[0].GetMatch(), "user.*"},
		{"folder.xattrFilter.maxSingleEntrySize", got.GetFolders()[0].GetXattrFilter().GetMaxSingleEntrySize(), int32(512)},
		{"folder.xattrFilter.maxTotalSize", got.GetFolders()[0].GetXattrFilter().GetMaxTotalSize(), int32(1024)},
		{"device.deviceID", got.GetDevices()[0].GetDeviceId(), otherID},
		{"device.compression", got.GetDevices()[0].GetCompression(), configpb.Compression_COMPRESSION_ALWAYS},
		{"device.numConnections", got.GetDevices()[0].GetNumConnections(), int32(7)},
		{"device.remoteGUIPort", got.GetDevices()[0].GetRemoteGuiPort(), int32(1234)},
		{"device.ignoredFolders[0].id", got.GetDevices()[0].GetIgnoredFolders()[0].GetId(), "ig1"},
		{"gui.address", got.GetGui().GetAddress(), "127.0.0.1:9000"},
		{"gui.authMode", got.GetGui().GetAuthMode(), configpb.AuthMode_AUTH_MODE_LDAP},
		{"gui.useTLS", got.GetGui().GetUseTls(), true},
		{"ldap.transport", got.GetLdap().GetTransport(), configpb.LDAPTransport_LDAP_TRANSPORT_START_TLS},
		{"options.listenAddresses", got.GetOptions().GetListenAddresses(), []string{"default"}},
		{"options.maxSendKbps", got.GetOptions().GetMaxSendKbps(), int32(100)},
		{"options.urAccepted", got.GetOptions().GetUrAccepted(), int32(-1)},
		{"options.minHomeDiskFree.percent", got.GetOptions().GetMinHomeDiskFree().GetPercent(), 1.0},
		{"remoteIgnoredDevices[0].deviceID", got.GetRemoteIgnoredDevices()[0].GetDeviceId(), id},
		{"remoteIgnoredDevices[0].name", got.GetRemoteIgnoredDevices()[0].GetName(), "Friend"},
		{"defaults.folder.rescanIntervalS", got.GetDefaults().GetFolder().GetRescanIntervalS(), int32(10800)},
		{"defaults.device.compression", got.GetDefaults().GetDevice().GetCompression(), configpb.Compression_COMPRESSION_NEVER},
		{"defaults.ignores.lines", got.GetDefaults().GetIgnores().GetLines(), []string{"*.tmp"}},
	}
	for _, test := range tests {
		if !reflect.DeepEqual(test.got, test.want) {
			t.Errorf("%s: got %v, want %v", test.name, test.got, test.want)
		}
	}

	if observed := got.GetDevices()[0].GetIgnoredFolders()[0]; !observed.GetTime().AsTime().Equal(ts) {
		t.Errorf("observed folder time: got %v, want %v", observed.GetTime().AsTime(), ts)
	}
	if observed := got.GetRemoteIgnoredDevices()[0]; !observed.GetTime().AsTime().Equal(ts) {
		t.Errorf("observed device time: got %v, want %v", observed.GetTime().AsTime(), ts)
	}

	// The converted configuration passes validation.
	if err := Validate(got); err != nil {
		t.Errorf("converted configuration failed validation: %v", err)
	}

	// The converted configuration round-trips through YAML.
	data, err := protoyaml.Marshal(got.Configuration)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var reparsed configpb.Configuration
	if err := protoyaml.Unmarshal(data, &reparsed); err != nil {
		t.Fatalf("unmarshal: %v\n%s", err, data)
	}
	if !proto.Equal(got.Configuration, &reparsed) {
		t.Errorf("converted configuration did not survive YAML round trip:\n%s", data)
	}
}

func TestFromLegacyNewConfig(t *testing.T) {
	id := protocol.NewDeviceID([]byte("newconfig"))

	got := FromLegacy(config.New(id))

	if v := got.GetVersion(); v != int32(config.CurrentVersion) {
		t.Errorf("version: got %d, want %d", v, config.CurrentVersion)
	}
	// The device list contains ourselves; a fresh legacy config has
	// materialised defaults which must carry over unchanged.
	if devices := got.GetDevices(); len(devices) != 1 || devices[0].GetDeviceId() != id {
		t.Errorf("devices: got %v, want single device %s", devices, id)
	}
	if v := got.GetGui().GetAddress(); v != "127.0.0.1:8384" {
		t.Errorf("gui.address: got %q, want default %q", v, "127.0.0.1:8384")
	}
	if v := got.GetOptions().GetListenAddresses(); len(v) != 1 || v[0] != "default" {
		t.Errorf("options.listenAddresses: got %v, want [default]", v)
	}
	if v := got.GetOptions().GetUnackedNotificationIds(); len(v) != 1 || v[0] != "authenticationUserAndPassword" {
		t.Errorf("options.unackedNotificationIDs: got %v", v)
	}
	if !got.HasDefaults() {
		t.Errorf("defaults not converted")
	}
	// Materialised legacy sizes ("1 %") convert to percentages.
	if got := got.GetOptions().GetMinHomeDiskFree().GetPercent(); got != 1 {
		t.Errorf("options.minHomeDiskFree.percent: got %v, want 1", got)
	}
	if got := got.GetDefaults().GetFolder().GetMinDiskFree().GetPercent(); got != 1 {
		t.Errorf("defaults.folder.minDiskFree.percent: got %v, want 1", got)
	}
	// The converted configuration passes validation, including the device
	// ID patterns (the device defaults have an empty device ID).
	if err := Validate(got); err != nil {
		t.Errorf("converted configuration failed validation: %v", err)
	}

	// Marshals without error.
	if _, err := protoyaml.Marshal(got.Configuration); err != nil {
		t.Fatalf("marshal: %v", err)
	}
}

// TestLegacyEnumValues pins the invariant that the conversion relies on:
// corresponding constants in the legacy and new enums must have the same
// numeric value.
func TestLegacyEnumValues(t *testing.T) {
	tests := []struct {
		desc   string
		legacy int64
		want   int64
	}{
		{"FolderType.SendReceive", int64(config.FolderTypeSendReceive), int64(configpb.FolderType_FOLDER_TYPE_SEND_RECEIVE)},
		{"FolderType.SendOnly", int64(config.FolderTypeSendOnly), int64(configpb.FolderType_FOLDER_TYPE_SEND_ONLY)},
		{"FolderType.ReceiveOnly", int64(config.FolderTypeReceiveOnly), int64(configpb.FolderType_FOLDER_TYPE_RECEIVE_ONLY)},
		{"FolderType.ReceiveEncrypted", int64(config.FolderTypeReceiveEncrypted), int64(configpb.FolderType_FOLDER_TYPE_RECEIVE_ENCRYPTED)},
		{"Compression.Metadata", int64(config.CompressionMetadata), int64(configpb.Compression_COMPRESSION_METADATA)},
		{"Compression.Never", int64(config.CompressionNever), int64(configpb.Compression_COMPRESSION_NEVER)},
		{"Compression.Always", int64(config.CompressionAlways), int64(configpb.Compression_COMPRESSION_ALWAYS)},
		{"AuthMode.Static", int64(config.AuthModeStatic), int64(configpb.AuthMode_AUTH_MODE_STATIC)},
		{"AuthMode.LDAP", int64(config.AuthModeLDAP), int64(configpb.AuthMode_AUTH_MODE_LDAP)},
		{"BlockPullOrder.Standard", int64(config.BlockPullOrderStandard), int64(configpb.BlockPullOrder_BLOCK_PULL_ORDER_STANDARD)},
		{"BlockPullOrder.Random", int64(config.BlockPullOrderRandom), int64(configpb.BlockPullOrder_BLOCK_PULL_ORDER_RANDOM)},
		{"BlockPullOrder.InOrder", int64(config.BlockPullOrderInOrder), int64(configpb.BlockPullOrder_BLOCK_PULL_ORDER_IN_ORDER)},
		{"CopyRangeMethod.Standard", int64(config.CopyRangeMethodStandard), int64(configpb.CopyRangeMethod_COPY_RANGE_METHOD_STANDARD)},
		{"CopyRangeMethod.Ioctl", int64(config.CopyRangeMethodIoctl), int64(configpb.CopyRangeMethod_COPY_RANGE_METHOD_IOCTL)},
		{"CopyRangeMethod.CopyFileRange", int64(config.CopyRangeMethodCopyFileRange), int64(configpb.CopyRangeMethod_COPY_RANGE_METHOD_COPY_FILE_RANGE)},
		{"CopyRangeMethod.SendFile", int64(config.CopyRangeMethodSendFile), int64(configpb.CopyRangeMethod_COPY_RANGE_METHOD_SEND_FILE)},
		{"CopyRangeMethod.DuplicateExtents", int64(config.CopyRangeMethodDuplicateExtents), int64(configpb.CopyRangeMethod_COPY_RANGE_METHOD_DUPLICATE_EXTENTS)},
		{"CopyRangeMethod.AllWithFallback", int64(config.CopyRangeMethodAllWithFallback), int64(configpb.CopyRangeMethod_COPY_RANGE_METHOD_ALL_WITH_FALLBACK)},
		{"LDAPTransport.Plain", int64(config.LDAPTransportPlain), int64(configpb.LDAPTransport_LDAP_TRANSPORT_PLAIN)},
		{"LDAPTransport.TLS", int64(config.LDAPTransportTLS), int64(configpb.LDAPTransport_LDAP_TRANSPORT_TLS)},
		{"LDAPTransport.StartTLS", int64(config.LDAPTransportStartTLS), int64(configpb.LDAPTransport_LDAP_TRANSPORT_START_TLS)},
		{"PullOrder.Random", int64(config.PullOrderRandom), int64(configpb.PullOrder_PULL_ORDER_RANDOM)},
		{"PullOrder.Alphabetic", int64(config.PullOrderAlphabetic), int64(configpb.PullOrder_PULL_ORDER_ALPHABETIC)},
		{"PullOrder.SmallestFirst", int64(config.PullOrderSmallestFirst), int64(configpb.PullOrder_PULL_ORDER_SMALLEST_FIRST)},
		{"PullOrder.LargestFirst", int64(config.PullOrderLargestFirst), int64(configpb.PullOrder_PULL_ORDER_LARGEST_FIRST)},
		{"PullOrder.OldestFirst", int64(config.PullOrderOldestFirst), int64(configpb.PullOrder_PULL_ORDER_OLDEST_FIRST)},
		{"PullOrder.NewestFirst", int64(config.PullOrderNewestFirst), int64(configpb.PullOrder_PULL_ORDER_NEWEST_FIRST)},
	}
	for _, test := range tests {
		if test.legacy != test.want {
			t.Errorf("%s: legacy value %d does not match new value %d", test.desc, test.legacy, test.want)
		}
	}

	// The legacy string based filesystem type maps to the new enum, with
	// the empty (zero) value meaning "basic".
	for _, test := range []struct {
		legacy config.FilesystemType
		want   configpb.FilesystemType
	}{
		{config.FilesystemTypeBasic, configpb.FilesystemType_FILESYSTEM_TYPE_BASIC},
		{config.FilesystemTypeFake, configpb.FilesystemType_FILESYSTEM_TYPE_FAKE},
		{"", configpb.FilesystemType_FILESYSTEM_TYPE_BASIC},
	} {
		if got := filesystemTypeFromLegacy(test.legacy); got != test.want {
			t.Errorf("filesystemTypeFromLegacy(%q): got %v, want %v", test.legacy, got, test.want)
		}
	}
}
