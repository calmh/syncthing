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

	"github.com/syncthing/syncthing/lib/structutil"

	intconfig "github.com/syncthing/syncthing/internal/config"
	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/protocol"
)

func TestFromLegacy(t *testing.T) {
	id := protocol.NewDeviceID([]byte("fromlegacy"))
	otherID := protocol.NewDeviceID([]byte("otherdevice"))
	ts := time.Date(2026, 9, 12, 12, 0, 0, 0, time.UTC)

	legacy := testLegacyConfiguration(id, otherID, ts)
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
		{"folder.filesystemType", got.GetFolders()[0].GetFilesystemType(), syncthingv2.FilesystemType_FILESYSTEM_TYPE_FAKE},
		{"folder.type", got.GetFolders()[0].GetType(), syncthingv2.FolderType_FOLDER_TYPE_RECEIVE_ONLY},
		{"folder.minDiskFree.gib", got.GetFolders()[0].GetMinDiskFree().GetGib(), 2.0},
		{"folder.order", got.GetFolders()[0].GetOrder(), syncthingv2.PullOrder_PULL_ORDER_OLDEST_FIRST},
		{"folder.blockPullOrder", got.GetFolders()[0].GetBlockPullOrder(), syncthingv2.BlockPullOrder_BLOCK_PULL_ORDER_IN_ORDER},
		{"folder.copyRangeMethod", got.GetFolders()[0].GetCopyRangeMethod(), syncthingv2.CopyRangeMethod_COPY_RANGE_METHOD_SEND_FILE},
		{"folder.versioning.type", got.GetFolders()[0].GetVersioning().GetType(), "simple"},
		{"folder.versioning.params", got.GetFolders()[0].GetVersioning().GetParams()["keep"], "3"},
		{"folder.versioning.fsType", got.GetFolders()[0].GetVersioning().GetFsType(), syncthingv2.FilesystemType_FILESYSTEM_TYPE_FAKE},
		{"folder.device.deviceID", got.GetFolders()[0].GetDevices()[0].GetDeviceId(), otherID},
		{"folder.device.introducedBy", got.GetFolders()[0].GetDevices()[0].GetIntroducedBy(), id},
		{"folder.device.encryptionPassword", got.GetFolders()[0].GetDevices()[0].GetEncryptionPassword(), "hunter2"},
		{"folder.xattrFilter.entry", got.GetFolders()[0].GetXattrFilter().GetEntries()[0].GetMatch(), "user.*"},
		{"folder.xattrFilter.maxSingleEntrySize", got.GetFolders()[0].GetXattrFilter().GetMaxSingleEntrySize(), int32(512)},
		{"folder.xattrFilter.maxTotalSize", got.GetFolders()[0].GetXattrFilter().GetMaxTotalSize(), int32(1024)},
		{"device.deviceID", got.GetDevices()[0].GetDeviceId(), otherID},
		{"device.compression", got.GetDevices()[0].GetCompression(), syncthingv2.Compression_COMPRESSION_ALWAYS},
		{"device.numConnections", got.GetDevices()[0].GetNumConnections(), int32(7)},
		{"device.remoteGUIPort", got.GetDevices()[0].GetRemoteGuiPort(), int32(1234)},
		{"device.ignoredFolders[0].id", got.GetDevices()[0].GetIgnoredFolders()[0].GetId(), "ig1"},
		{"gui.address", got.GetGui().GetAddress(), "127.0.0.1:9000"},
		{"gui.authMode", got.GetGui().GetAuthMode(), syncthingv2.AuthMode_AUTH_MODE_LDAP},
		{"gui.useTLS", got.GetGui().GetUseTls(), true},
		{"ldap.transport", got.GetLdap().GetTransport(), syncthingv2.LDAPTransport_LDAP_TRANSPORT_START_TLS},
		{"options.listenAddresses", got.GetOptions().GetListenAddresses(), []string{"default"}},
		{"options.maxSendKbps", got.GetOptions().GetMaxSendKbps(), int32(100)},
		{"options.urAccepted", got.GetOptions().GetUrAccepted(), int32(-1)},
		{"options.minHomeDiskFree.percent", got.GetOptions().GetMinHomeDiskFree().GetPercent(), 1.0},
		{"remoteIgnoredDevices[0].deviceID", got.GetRemoteIgnoredDevices()[0].GetDeviceId(), id},
		{"remoteIgnoredDevices[0].name", got.GetRemoteIgnoredDevices()[0].GetName(), "Friend"},
		{"defaults.folder.rescanIntervalS", got.GetDefaults().GetFolder().GetRescanIntervalS(), int32(10800)},
		{"defaults.device.compression", got.GetDefaults().GetDevice().GetCompression(), syncthingv2.Compression_COMPRESSION_NEVER},
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
	if err := intconfig.Validate(got); err != nil {
		t.Errorf("converted configuration failed validation: %v", err)
	}

	// The converted configuration round-trips through YAML.
	data, err := protoyaml.Marshal(got.Configuration)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var reparsed syncthingv2.Configuration
	if err := protoyaml.Unmarshal(data, &reparsed); err != nil {
		t.Fatalf("unmarshal: %v\n%s", err, data)
	}
	if !proto.Equal(got.Configuration, &reparsed) {
		t.Errorf("converted configuration did not survive YAML round trip:\n%s", data)
	}
}

func TestFromLegacyNewConfig(t *testing.T) {
	id := protocol.NewDeviceID([]byte("newconfig"))

	got := FromLegacy(New(id))

	if v := got.GetVersion(); v != int32(CurrentVersion) {
		t.Errorf("version: got %d, want %d", v, CurrentVersion)
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
	if err := intconfig.Validate(got); err != nil {
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
		{"FolderType.SendReceive", int64(FolderTypeSendReceive), int64(syncthingv2.FolderType_FOLDER_TYPE_SEND_RECEIVE)},
		{"FolderType.SendOnly", int64(FolderTypeSendOnly), int64(syncthingv2.FolderType_FOLDER_TYPE_SEND_ONLY)},
		{"FolderType.ReceiveOnly", int64(FolderTypeReceiveOnly), int64(syncthingv2.FolderType_FOLDER_TYPE_RECEIVE_ONLY)},
		{"FolderType.ReceiveEncrypted", int64(FolderTypeReceiveEncrypted), int64(syncthingv2.FolderType_FOLDER_TYPE_RECEIVE_ENCRYPTED)},
		{"Compression.Metadata", int64(CompressionMetadata), int64(syncthingv2.Compression_COMPRESSION_METADATA)},
		{"Compression.Never", int64(CompressionNever), int64(syncthingv2.Compression_COMPRESSION_NEVER)},
		{"Compression.Always", int64(CompressionAlways), int64(syncthingv2.Compression_COMPRESSION_ALWAYS)},
		{"AuthMode.Static", int64(AuthModeStatic), int64(syncthingv2.AuthMode_AUTH_MODE_STATIC)},
		{"AuthMode.LDAP", int64(AuthModeLDAP), int64(syncthingv2.AuthMode_AUTH_MODE_LDAP)},
		{"BlockPullOrder.Standard", int64(BlockPullOrderStandard), int64(syncthingv2.BlockPullOrder_BLOCK_PULL_ORDER_STANDARD)},
		{"BlockPullOrder.Random", int64(BlockPullOrderRandom), int64(syncthingv2.BlockPullOrder_BLOCK_PULL_ORDER_RANDOM)},
		{"BlockPullOrder.InOrder", int64(BlockPullOrderInOrder), int64(syncthingv2.BlockPullOrder_BLOCK_PULL_ORDER_IN_ORDER)},
		{"CopyRangeMethod.Standard", int64(CopyRangeMethodStandard), int64(syncthingv2.CopyRangeMethod_COPY_RANGE_METHOD_STANDARD)},
		{"CopyRangeMethod.Ioctl", int64(CopyRangeMethodIoctl), int64(syncthingv2.CopyRangeMethod_COPY_RANGE_METHOD_IOCTL)},
		{"CopyRangeMethod.CopyFileRange", int64(CopyRangeMethodCopyFileRange), int64(syncthingv2.CopyRangeMethod_COPY_RANGE_METHOD_COPY_FILE_RANGE)},
		{"CopyRangeMethod.SendFile", int64(CopyRangeMethodSendFile), int64(syncthingv2.CopyRangeMethod_COPY_RANGE_METHOD_SEND_FILE)},
		{"CopyRangeMethod.DuplicateExtents", int64(CopyRangeMethodDuplicateExtents), int64(syncthingv2.CopyRangeMethod_COPY_RANGE_METHOD_DUPLICATE_EXTENTS)},
		{"CopyRangeMethod.AllWithFallback", int64(CopyRangeMethodAllWithFallback), int64(syncthingv2.CopyRangeMethod_COPY_RANGE_METHOD_ALL_WITH_FALLBACK)},
		{"LDAPTransport.Plain", int64(LDAPTransportPlain), int64(syncthingv2.LDAPTransport_LDAP_TRANSPORT_PLAIN)},
		{"LDAPTransport.TLS", int64(LDAPTransportTLS), int64(syncthingv2.LDAPTransport_LDAP_TRANSPORT_TLS)},
		{"LDAPTransport.StartTLS", int64(LDAPTransportStartTLS), int64(syncthingv2.LDAPTransport_LDAP_TRANSPORT_START_TLS)},
		{"PullOrder.Random", int64(PullOrderRandom), int64(syncthingv2.PullOrder_PULL_ORDER_RANDOM)},
		{"PullOrder.Alphabetic", int64(PullOrderAlphabetic), int64(syncthingv2.PullOrder_PULL_ORDER_ALPHABETIC)},
		{"PullOrder.SmallestFirst", int64(PullOrderSmallestFirst), int64(syncthingv2.PullOrder_PULL_ORDER_SMALLEST_FIRST)},
		{"PullOrder.LargestFirst", int64(PullOrderLargestFirst), int64(syncthingv2.PullOrder_PULL_ORDER_LARGEST_FIRST)},
		{"PullOrder.OldestFirst", int64(PullOrderOldestFirst), int64(syncthingv2.PullOrder_PULL_ORDER_OLDEST_FIRST)},
		{"PullOrder.NewestFirst", int64(PullOrderNewestFirst), int64(syncthingv2.PullOrder_PULL_ORDER_NEWEST_FIRST)},
	}
	for _, test := range tests {
		if test.legacy != test.want {
			t.Errorf("%s: legacy value %d does not match new value %d", test.desc, test.legacy, test.want)
		}
	}

	// The legacy string based filesystem type maps to the new enum, with
	// the empty (zero) value meaning "basic".
	for _, test := range []struct {
		legacy FilesystemType
		want   syncthingv2.FilesystemType
	}{
		{FilesystemTypeBasic, syncthingv2.FilesystemType_FILESYSTEM_TYPE_BASIC},
		{FilesystemTypeFake, syncthingv2.FilesystemType_FILESYSTEM_TYPE_FAKE},
		{"", syncthingv2.FilesystemType_FILESYSTEM_TYPE_BASIC},
	} {
		if got := filesystemTypeFromLegacy(test.legacy); got != test.want {
			t.Errorf("filesystemTypeFromLegacy(%q): got %v, want %v", test.legacy, got, test.want)
		}
	}
}

// testLegacyConfiguration returns a legacy configuration populated
// with a representative spread of fields and values.
func testLegacyConfiguration(id, otherID protocol.DeviceID, ts time.Time) Configuration {
	return Configuration{
		Version: 52,
		Folders: []FolderConfiguration{{
			ID:              "f1",
			Label:           "Folder One",
			Path:            "/srv/sync/f1",
			FilesystemType:  FilesystemTypeFake,
			Type:            FolderTypeReceiveOnly,
			RescanIntervalS: 0, // explicit zero: no periodic rescan
			Devices: []FolderDeviceConfiguration{{
				DeviceID:           otherID,
				IntroducedBy:       id,
				EncryptionPassword: "hunter2",
			}},
			MinDiskFree: Size{Value: 2, Unit: "GiB"},
			Versioning: VersioningConfiguration{
				Type:             "simple",
				Params:           map[string]string{"keep": "3"},
				CleanupIntervalS: 0, // explicit zero: no cleanup
				FSType:           FilesystemTypeFake,
			},
			MaxConflicts:    0, // explicit zero
			Copiers:         2,
			Hashers:         4,
			Order:           PullOrderOldestFirst,
			BlockPullOrder:  BlockPullOrderInOrder,
			CopyRangeMethod: CopyRangeMethodSendFile,
			XattrFilter: XattrFilter{
				Entries:            []XattrFilterEntry{{Match: "user.*", Permit: true}},
				MaxSingleEntrySize: 512,
				MaxTotalSize:       1024,
			},
		}},
		Devices: []DeviceConfiguration{{
			DeviceID:          otherID,
			Name:              "Other",
			Addresses:         []string{"dynamic"},
			Compression:       CompressionAlways,
			MaxSendKbps:       1000,
			IgnoredFolders:    []ObservedFolder{{Time: ts, ID: "ig1", Label: "Ignored"}},
			RemoteGUIPort:     1234,
			RawNumConnections: 7,
			Untrusted:         true,
		}},
		GUI: GUIConfiguration{
			RawAddress: "127.0.0.1:9000",
			AuthMode:   AuthModeLDAP,
			RawUseTLS:  true,
		},
		LDAP: LDAPConfiguration{
			Address:   "ldap.example.com:389",
			Transport: LDAPTransportStartTLS,
		},
		Options: OptionsConfiguration{
			RawListenAddresses: []string{"default"},
			MaxSendKbps:        100,
			URAccepted:         -1,
			MinHomeDiskFree:    Size{Value: 1, Unit: "%"},
		},
		IgnoredDevices: []ObservedDevice{{Time: ts, ID: id, Name: "Friend", Address: "192.0.2.1:22000"}},
		Defaults: Defaults{
			Folder:  FolderConfiguration{RescanIntervalS: 10800},
			Device:  DeviceConfiguration{Compression: CompressionNever},
			Ignores: Ignores{Lines: []string{"*.tmp"}},
		},
	}
}

func TestToLegacy(t *testing.T) {
	id := protocol.NewDeviceID([]byte("tolegacy"))

	cfg := intconfig.Configuration{Configuration: syncthingv2.Configuration_builder{
		Version: new(int32(52)),
		Folders: []*syncthingv2.FolderConfiguration{
			syncthingv2.FolderConfiguration_builder{
				Id:   new("f1"),
				Path: new("/srv/sync/f1"),
				Type: new(syncthingv2.FolderType_FOLDER_TYPE_SEND_ONLY),
				// Most fields unset; they should get their defaults.
			}.Build(),
			syncthingv2.FolderConfiguration_builder{
				Id:          new("f2"),
				Path:        new("/srv/sync/f2"),
				MinDiskFree: syncthingv2.Size_builder{Percent: new(0.0)}.Build(), // explicit zero: disabled
			}.Build(),
		},
		Devices: []*syncthingv2.DeviceConfiguration{syncthingv2.DeviceConfiguration_builder{
			DeviceId: new(id.String()),
		}.Build()},
	}.Build()}

	got := ToLegacy(cfg)

	// Unset fields get their defaults, like a loaded legacy
	// configuration.
	folder := got.Folders[0]
	if folder.RescanIntervalS != 3600 {
		t.Errorf("rescanIntervalS: got %d, want default 3600", folder.RescanIntervalS)
	}
	if folder.MarkerName != DefaultMarkerName {
		t.Errorf("markerName: got %q, want default %q", folder.MarkerName, DefaultMarkerName)
	}
	if folder.MaxConcurrentWrites != 16 {
		t.Errorf("maxConcurrentWrites: got %d, want default 16", folder.MaxConcurrentWrites)
	}
	if folder.MinDiskFree != (Size{Value: 1, Unit: "%"}) {
		t.Errorf("minDiskFree: got %v, want default one percent", folder.MinDiskFree)
	}
	if folder.FilesystemType != FilesystemTypeBasic {
		t.Errorf("filesystemType: got %q, want basic", folder.FilesystemType)
	}
	// Set values are kept.
	if folder.ID != "f1" || folder.Path != "/srv/sync/f1" || folder.Type != FolderTypeSendOnly {
		t.Errorf("folder values not kept: %+v", folder)
	}

	// An explicit zero size disables the check, as the legacy zero size.
	if got.Folders[1].MinDiskFree != (Size{}) {
		t.Errorf("minDiskFree: got %v, want zero size", got.Folders[1].MinDiskFree)
	}

	// Device IDs are native and addresses default to dynamic discovery.
	device := got.Devices[0]
	if device.DeviceID != id {
		t.Errorf("deviceID: got %v, want %v", device.DeviceID, id)
	}
	if !reflect.DeepEqual(device.Addresses, []string{"dynamic"}) {
		t.Errorf("addresses: got %v, want [dynamic]", device.Addresses)
	}

	// The GUI and options sections get their defaults.
	if !got.GUI.Enabled || got.GUI.RawAddress != "127.0.0.1:8384" || got.GUI.Theme != DefaultTheme {
		t.Errorf("GUI not defaulted: %+v", got.GUI)
	}
	if !reflect.DeepEqual(got.Options.RawListenAddresses, []string{"default"}) {
		t.Errorf("listenAddresses: got %v, want [default]", got.Options.RawListenAddresses)
	}
	if got.Options.MinHomeDiskFree != (Size{Value: 1, Unit: "%"}) {
		t.Errorf("minHomeDiskFree: got %v, want default one percent", got.Options.MinHomeDiskFree)
	}
	if got.Options.ReconnectIntervalS != 20 {
		t.Errorf("reconnectionIntervalS: got %d, want default 20", got.Options.ReconnectIntervalS)
	}
}

func TestLegacyRoundTrip(t *testing.T) {
	id := protocol.NewDeviceID([]byte("roundtrip"))
	otherID := protocol.NewDeviceID([]byte("roundtripother"))
	ts := time.Date(2026, 9, 12, 12, 0, 0, 0, time.UTC)

	// A loaded legacy configuration has its defaults materialised, and
	// ToLegacy produces effective values; materialise the fixture the
	// same way before comparing. The legacy default mechanism covers
	// neither top level folders nor slice fields, so those are set
	// explicitly.
	orig := testLegacyConfiguration(id, otherID, ts)
	structutil.SetDefaults(&orig)
	for i := range orig.Folders {
		structutil.SetDefaults(&orig.Folders[i])
	}
	orig.Options.RawGlobalAnnServers = []string{"default"}
	orig.Options.RawStunServers = []string{"default"}
	orig.Defaults.Device.Addresses = []string{"dynamic"}

	got := ToLegacy(FromLegacy(orig))

	if !reflect.DeepEqual(orig, got) {
		t.Errorf("legacy configuration did not survive round trip:\ngot:  %+v\nwant: %+v", got, orig)
	}
}
