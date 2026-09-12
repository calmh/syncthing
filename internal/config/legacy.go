// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"maps"
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	configpb "github.com/syncthing/syncthing/internal/gen/config"
	config "github.com/syncthing/syncthing/lib/config"
)

// FromLegacy returns a new format configuration with the same contents as
// the given legacy configuration.
//
// All fields are copied verbatim. Since the legacy types cannot distinguish
// between unset fields and fields set to their zero value, every field is
// set in the result; zero values that are meaningful in the legacy format
// (e.g. rescanIntervalS 0 for no periodic rescan, or a disabled GUI) are
// preserved as such. Legacy fields that were deprecated have no counterpart
// in the new format and are dropped.
func FromLegacy(cfg config.Configuration) *configpb.Configuration {
	return configpb.Configuration_builder{
		Version:              proto.Int32(int32(cfg.Version)),
		Folders:              fromLegacyFolders(cfg.Folders),
		Devices:              fromLegacyDevices(cfg.Devices),
		Gui:                  fromLegacyGUI(cfg.GUI),
		Ldap:                 fromLegacyLDAP(cfg.LDAP),
		Options:              fromLegacyOptions(cfg.Options),
		RemoteIgnoredDevices: fromLegacyObservedDevices(cfg.IgnoredDevices),
		Defaults:             fromLegacyDefaults(cfg.Defaults),
	}.Build()
}

func fromLegacyFolders(folders []config.FolderConfiguration) []*configpb.FolderConfiguration {
	out := make([]*configpb.FolderConfiguration, len(folders))
	for i, folder := range folders {
		out[i] = fromLegacyFolder(folder)
	}
	return out
}

func fromLegacyFolder(f config.FolderConfiguration) *configpb.FolderConfiguration {
	return configpb.FolderConfiguration_builder{
		Id:                      proto.String(f.ID),
		Label:                   proto.String(f.Label),
		FilesystemType:          filesystemTypeFromLegacy(f.FilesystemType).Enum(),
		Path:                    proto.String(f.Path),
		Type:                    configpb.FolderType(f.Type).Enum(),
		Devices:                 fromLegacyFolderDevices(f.Devices),
		Group:                   proto.String(f.Group),
		RescanIntervalS:         proto.Int32(int32(f.RescanIntervalS)),
		FsWatcherEnabled:        proto.Bool(f.FSWatcherEnabled),
		FsWatcherDelayS:         proto.Float64(f.FSWatcherDelayS),
		FsWatcherTimeoutS:       proto.Float64(f.FSWatcherTimeoutS),
		IgnorePerms:             proto.Bool(f.IgnorePerms),
		AutoNormalize:           proto.Bool(f.AutoNormalize),
		MinDiskFree:             fromLegacySize(f.MinDiskFree),
		Versioning:              fromLegacyVersioning(f.Versioning),
		Copiers:                 proto.Int32(int32(f.Copiers)),
		PullerMaxPendingKiB:     proto.Int32(int32(f.PullerMaxPendingKiB)),
		Hashers:                 proto.Int32(int32(f.Hashers)),
		Order:                   configpb.PullOrder(f.Order).Enum(),
		IgnoreDelete:            proto.Bool(f.IgnoreDelete),
		ScanProgressIntervalS:   proto.Int32(int32(f.ScanProgressIntervalS)),
		PullerPauseS:            proto.Int32(int32(f.PullerPauseS)),
		PullerDelayS:            proto.Float64(f.PullerDelayS),
		MaxConflicts:            proto.Int32(int32(f.MaxConflicts)),
		DisableSparseFiles:      proto.Bool(f.DisableSparseFiles),
		Paused:                  proto.Bool(f.Paused),
		MarkerName:              proto.String(f.MarkerName),
		CopyOwnershipFromParent: proto.Bool(f.CopyOwnershipFromParent),
		ModTimeWindowS:          proto.Int32(int32(f.RawModTimeWindowS)),
		MaxConcurrentWrites:     proto.Int32(int32(f.MaxConcurrentWrites)),
		DisableFsync:            proto.Bool(f.DisableFsync),
		BlockPullOrder:          configpb.BlockPullOrder(f.BlockPullOrder).Enum(),
		CopyRangeMethod:         configpb.CopyRangeMethod(f.CopyRangeMethod).Enum(),
		CaseSensitiveFs:         proto.Bool(f.CaseSensitiveFS),
		JunctionsAsDirs:         proto.Bool(f.JunctionsAsDirs),
		SyncOwnership:           proto.Bool(f.SyncOwnership),
		SendOwnership:           proto.Bool(f.SendOwnership),
		SyncXattrs:              proto.Bool(f.SyncXattrs),
		SendXattrs:              proto.Bool(f.SendXattrs),
		BlockIndexing:           proto.Bool(f.BlockIndexing),
		XattrFilter:             fromLegacyXattrFilter(f.XattrFilter),
	}.Build()
}

func fromLegacyFolderDevices(devices []config.FolderDeviceConfiguration) []*configpb.FolderDeviceConfiguration {
	out := make([]*configpb.FolderDeviceConfiguration, len(devices))
	for i, device := range devices {
		out[i] = configpb.FolderDeviceConfiguration_builder{
			DeviceId:           proto.String(device.DeviceID.String()),
			IntroducedBy:       proto.String(device.IntroducedBy.String()),
			EncryptionPassword: proto.String(device.EncryptionPassword),
		}.Build()
	}
	return out
}

func fromLegacyVersioning(v config.VersioningConfiguration) *configpb.VersioningConfiguration {
	return configpb.VersioningConfiguration_builder{
		Type:             proto.String(v.Type),
		Params:           maps.Clone(v.Params),
		CleanupIntervalS: proto.Int32(int32(v.CleanupIntervalS)),
		FsPath:           proto.String(v.FSPath),
		FsType:           filesystemTypeFromLegacy(v.FSType).Enum(),
	}.Build()
}

func fromLegacyXattrFilter(f config.XattrFilter) *configpb.XattrFilter {
	entries := make([]*configpb.XattrFilterEntry, len(f.Entries))
	for i, entry := range f.Entries {
		entries[i] = configpb.XattrFilterEntry_builder{
			Match:  proto.String(entry.Match),
			Permit: proto.Bool(entry.Permit),
		}.Build()
	}
	return configpb.XattrFilter_builder{
		Entries:            entries,
		MaxSingleEntrySize: proto.Int32(int32(f.MaxSingleEntrySize)),
		MaxTotalSize:       proto.Int32(int32(f.MaxTotalSize)),
	}.Build()
}

// fromLegacySize converts a legacy size to the new representation. A
// legacy size of zero or less meant "no minimum" and converts to an unset
// size. Units with an "iB"-suffix are binary, other unit prefixes decimal,
// matching the intent of the legacy string format.
func fromLegacySize(s config.Size) *configpb.Size {
	if s.Value <= 0 {
		return nil
	}
	if strings.Contains(s.Unit, "%") {
		return configpb.Size_builder{Percent: proto.Float64(s.Value)}.Build()
	}
	switch strings.ToLower(s.Unit) {
	case "mib":
		return configpb.Size_builder{Mib: proto.Float64(s.Value)}.Build()
	case "gib":
		return configpb.Size_builder{Gib: proto.Float64(s.Value)}.Build()
	case "kib":
		return configpb.Size_builder{Bytes: proto.Float64(s.Value * 1024)}.Build()
	case "tib":
		return configpb.Size_builder{Bytes: proto.Float64(s.Value * (1 << 40))}.Build()
	case "k", "kb":
		return configpb.Size_builder{Bytes: proto.Float64(s.Value * 1000)}.Build()
	case "m", "mb":
		return configpb.Size_builder{Bytes: proto.Float64(s.Value * 1000 * 1000)}.Build()
	case "g", "gb":
		return configpb.Size_builder{Bytes: proto.Float64(s.Value * 1000 * 1000 * 1000)}.Build()
	case "t", "tb":
		return configpb.Size_builder{Bytes: proto.Float64(s.Value * 1000 * 1000 * 1000 * 1000)}.Build()
	default:
		return configpb.Size_builder{Bytes: proto.Float64(s.Value)}.Build()
	}
}

func fromLegacyDevices(devices []config.DeviceConfiguration) []*configpb.DeviceConfiguration {
	out := make([]*configpb.DeviceConfiguration, len(devices))
	for i, device := range devices {
		out[i] = fromLegacyDevice(device)
	}
	return out
}

func fromLegacyDevice(d config.DeviceConfiguration) *configpb.DeviceConfiguration {
	return configpb.DeviceConfiguration_builder{
		DeviceId:                 proto.String(d.DeviceID.String()),
		Name:                     proto.String(d.Name),
		Addresses:                slices.Clone(d.Addresses),
		Compression:              configpb.Compression(d.Compression).Enum(),
		CertName:                 proto.String(d.CertName),
		Introducer:               proto.Bool(d.Introducer),
		SkipIntroductionRemovals: proto.Bool(d.SkipIntroductionRemovals),
		IntroducedBy:             proto.String(d.IntroducedBy.String()),
		Paused:                   proto.Bool(d.Paused),
		AllowedNetworks:          slices.Clone(d.AllowedNetworks),
		AutoAcceptFolders:        proto.Bool(d.AutoAcceptFolders),
		MaxSendKbps:              proto.Int32(int32(d.MaxSendKbps)),
		MaxRecvKbps:              proto.Int32(int32(d.MaxRecvKbps)),
		IgnoredFolders:           fromLegacyObservedFolders(d.IgnoredFolders),
		MaxRequestKiB:            proto.Int32(int32(d.MaxRequestKiB)),
		Untrusted:                proto.Bool(d.Untrusted),
		RemoteGuiPort:            proto.Int32(int32(d.RemoteGUIPort)),
		NumConnections:           proto.Int32(int32(d.RawNumConnections)),
		Group:                    proto.String(d.Group),
	}.Build()
}

func fromLegacyObservedFolders(folders []config.ObservedFolder) []*configpb.ObservedFolder {
	out := make([]*configpb.ObservedFolder, len(folders))
	for i, folder := range folders {
		out[i] = configpb.ObservedFolder_builder{
			Time:  timestamppb.New(folder.Time),
			Id:    proto.String(folder.ID),
			Label: proto.String(folder.Label),
		}.Build()
	}
	return out
}

func fromLegacyObservedDevices(devices []config.ObservedDevice) []*configpb.ObservedDevice {
	out := make([]*configpb.ObservedDevice, len(devices))
	for i, device := range devices {
		out[i] = configpb.ObservedDevice_builder{
			Time:     timestamppb.New(device.Time),
			DeviceId: proto.String(device.ID.String()),
			Name:     proto.String(device.Name),
			Address:  proto.String(device.Address),
		}.Build()
	}
	return out
}

func fromLegacyGUI(g config.GUIConfiguration) *configpb.GUIConfiguration {
	return configpb.GUIConfiguration_builder{
		Enabled:                   proto.Bool(g.Enabled),
		Address:                   proto.String(g.RawAddress),
		UnixSocketPermissions:     proto.String(g.RawUnixSocketPermissions),
		User:                      proto.String(g.User),
		Password:                  proto.String(g.Password),
		AuthMode:                  configpb.AuthMode(g.AuthMode).Enum(),
		MetricsWithoutAuth:        proto.Bool(g.MetricsWithoutAuth),
		UseTls:                    proto.Bool(g.RawUseTLS),
		ApiKey:                    proto.String(g.APIKey),
		InsecureAdminAccess:       proto.Bool(g.InsecureAdminAccess),
		Theme:                     proto.String(g.Theme),
		InsecureSkipHostcheck:     proto.Bool(g.InsecureSkipHostCheck),
		InsecureAllowFrameLoading: proto.Bool(g.InsecureAllowFrameLoading),
		SendBasicAuthPrompt:       proto.Bool(g.SendBasicAuthPrompt),
		SessionCookieDurationS:    proto.Int32(int32(g.SessionCookieDurationS)),
		SessionCookiePath:         proto.String(g.SessionCookiePath),
	}.Build()
}

func fromLegacyLDAP(l config.LDAPConfiguration) *configpb.LDAPConfiguration {
	return configpb.LDAPConfiguration_builder{
		Address:            proto.String(l.Address),
		BindDn:             proto.String(l.BindDN),
		Transport:          configpb.LDAPTransport(l.Transport).Enum(),
		InsecureSkipVerify: proto.Bool(l.InsecureSkipVerify),
		SearchBaseDn:       proto.String(l.SearchBaseDN),
		SearchFilter:       proto.String(l.SearchFilter),
	}.Build()
}

func fromLegacyOptions(o config.OptionsConfiguration) *configpb.OptionsConfiguration {
	return configpb.OptionsConfiguration_builder{
		ListenAddresses:                     slices.Clone(o.RawListenAddresses),
		GlobalAnnounceServers:               slices.Clone(o.RawGlobalAnnServers),
		GlobalAnnounceEnabled:               proto.Bool(o.GlobalAnnEnabled),
		LocalAnnounceEnabled:                proto.Bool(o.LocalAnnEnabled),
		LocalAnnouncePort:                   proto.Int32(int32(o.LocalAnnPort)),
		LocalAnnounceMcAddr:                 proto.String(o.LocalAnnMCAddr),
		MaxSendKbps:                         proto.Int32(int32(o.MaxSendKbps)),
		MaxRecvKbps:                         proto.Int32(int32(o.MaxRecvKbps)),
		ReconnectionIntervalS:               proto.Int32(int32(o.ReconnectIntervalS)),
		RelaysEnabled:                       proto.Bool(o.RelaysEnabled),
		RelayReconnectIntervalM:             proto.Int32(int32(o.RelayReconnectIntervalM)),
		StartBrowser:                        proto.Bool(o.StartBrowser),
		NatEnabled:                          proto.Bool(o.NATEnabled),
		NatLeaseMinutes:                     proto.Int32(int32(o.NATLeaseM)),
		NatRenewalMinutes:                   proto.Int32(int32(o.NATRenewalM)),
		NatTimeoutSeconds:                   proto.Int32(int32(o.NATTimeoutS)),
		UrAccepted:                          proto.Int32(int32(o.URAccepted)),
		UrSeen:                              proto.Int32(int32(o.URSeen)),
		UrUniqueId:                          proto.String(o.URUniqueID),
		UrUrl:                               proto.String(o.URURL),
		UrPostInsecurely:                    proto.Bool(o.URPostInsecurely),
		UrInitialDelayS:                     proto.Int32(int32(o.URInitialDelayS)),
		AutoUpgradeIntervalH:                proto.Int32(int32(o.AutoUpgradeIntervalH)),
		UpgradeToPreReleases:                proto.Bool(o.UpgradeToPreReleases),
		KeepTemporariesH:                    proto.Int32(int32(o.KeepTemporariesH)),
		ProgressUpdateIntervalS:             proto.Int32(int32(o.ProgressUpdateIntervalS)),
		LimitBandwidthInLan:                 proto.Bool(o.LimitBandwidthInLan),
		MinHomeDiskFree:                     fromLegacySize(o.MinHomeDiskFree),
		ReleasesUrl:                         proto.String(o.ReleasesURL),
		AlwaysLocalNets:                     slices.Clone(o.AlwaysLocalNets),
		OverwriteRemoteDeviceNamesOnConnect: proto.Bool(o.OverwriteRemoteDevNames),
		TempIndexMinBlocks:                  proto.Int32(int32(o.TempIndexMinBlocks)),
		UnackedNotificationIds:              slices.Clone(o.UnackedNotificationIDs),
		TrafficClass:                        proto.Int32(int32(o.TrafficClass)),
		SetLowPriority:                      proto.Bool(o.SetLowPriority),
		MaxFolderConcurrency:                proto.Int32(int32(o.RawMaxFolderConcurrency)),
		CrUrl:                               proto.String(o.CRURL),
		CrashReportingEnabled:               proto.Bool(o.CREnabled),
		StunKeepaliveStartS:                 proto.Int32(int32(o.StunKeepaliveStartS)),
		StunKeepaliveMinS:                   proto.Int32(int32(o.StunKeepaliveMinS)),
		StunServers:                         slices.Clone(o.RawStunServers),
		MaxConcurrentIncomingRequestKiB:     proto.Int32(int32(o.RawMaxCIRequestKiB)),
		AnnounceLanAddresses:                proto.Bool(o.AnnounceLANAddresses),
		SendFullIndexOnUpgrade:              proto.Bool(o.SendFullIndexOnUpgrade),
		FeatureFlags:                        slices.Clone(o.FeatureFlags),
		AuditEnabled:                        proto.Bool(o.AuditEnabled),
		AuditFile:                           proto.String(o.AuditFile),
		ConnectionLimitEnough:               proto.Int32(int32(o.ConnectionLimitEnough)),
		ConnectionLimitMax:                  proto.Int32(int32(o.ConnectionLimitMax)),
		ConnectionPriorityTcpLan:            proto.Int32(int32(o.ConnectionPriorityTCPLAN)),
		ConnectionPriorityQuicLan:           proto.Int32(int32(o.ConnectionPriorityQUICLAN)),
		ConnectionPriorityTcpWan:            proto.Int32(int32(o.ConnectionPriorityTCPWAN)),
		ConnectionPriorityQuicWan:           proto.Int32(int32(o.ConnectionPriorityQUICWAN)),
		ConnectionPriorityRelay:             proto.Int32(int32(o.ConnectionPriorityRelay)),
		ConnectionPriorityUpgradeThreshold:  proto.Int32(int32(o.ConnectionPriorityUpgradeThreshold)),
	}.Build()
}

func fromLegacyDefaults(d config.Defaults) *configpb.Defaults {
	return configpb.Defaults_builder{
		Folder:  fromLegacyFolder(d.Folder),
		Device:  fromLegacyDevice(d.Device),
		Ignores: configpb.Ignores_builder{Lines: slices.Clone(d.Ignores.Lines)}.Build(),
	}.Build()
}

// filesystemTypeFromLegacy maps the legacy string based filesystem type to
// the new enum. The legacy zero value ("") means "basic".
func filesystemTypeFromLegacy(t config.FilesystemType) configpb.FilesystemType {
	if t == config.FilesystemTypeFake {
		return configpb.FilesystemType_FILESYSTEM_TYPE_FAKE
	}
	return configpb.FilesystemType_FILESYSTEM_TYPE_BASIC
}
