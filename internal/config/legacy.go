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

	"google.golang.org/protobuf/types/known/timestamppb"

	configpb "github.com/syncthing/syncthing/internal/gen/syncthing/v2/config"
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
func FromLegacy(cfg config.Configuration) Configuration {
	return Configuration{configpb.Configuration_builder{
		Version:              new(int32(cfg.Version)),
		Folders:              fromLegacyFolders(cfg.Folders),
		Devices:              fromLegacyDevices(cfg.Devices),
		Gui:                  fromLegacyGUI(cfg.GUI),
		Ldap:                 fromLegacyLDAP(cfg.LDAP),
		Options:              fromLegacyOptions(cfg.Options),
		RemoteIgnoredDevices: fromLegacyObservedDevices(cfg.IgnoredDevices),
		Defaults:             fromLegacyDefaults(cfg.Defaults),
	}.Build()}
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
		Id:                      new(f.ID),
		Label:                   new(f.Label),
		FilesystemType:          new(filesystemTypeFromLegacy(f.FilesystemType)),
		Path:                    new(f.Path),
		Type:                    new(configpb.FolderType(f.Type)),
		Devices:                 fromLegacyFolderDevices(f.Devices),
		Group:                   new(f.Group),
		RescanIntervalS:         new(int32(f.RescanIntervalS)),
		FsWatcherEnabled:        new(f.FSWatcherEnabled),
		FsWatcherDelayS:         new(f.FSWatcherDelayS),
		FsWatcherTimeoutS:       new(f.FSWatcherTimeoutS),
		IgnorePerms:             new(f.IgnorePerms),
		AutoNormalize:           new(f.AutoNormalize),
		MinDiskFree:             fromLegacySize(f.MinDiskFree),
		Versioning:              fromLegacyVersioning(f.Versioning),
		Copiers:                 new(int32(f.Copiers)),
		PullerMaxPendingKiB:     new(int32(f.PullerMaxPendingKiB)),
		Hashers:                 new(int32(f.Hashers)),
		Order:                   new(configpb.PullOrder(f.Order)),
		IgnoreDelete:            new(f.IgnoreDelete),
		ScanProgressIntervalS:   new(int32(f.ScanProgressIntervalS)),
		PullerPauseS:            new(int32(f.PullerPauseS)),
		PullerDelayS:            new(f.PullerDelayS),
		MaxConflicts:            new(int32(f.MaxConflicts)),
		DisableSparseFiles:      new(f.DisableSparseFiles),
		Paused:                  new(f.Paused),
		MarkerName:              new(f.MarkerName),
		CopyOwnershipFromParent: new(f.CopyOwnershipFromParent),
		ModTimeWindowS:          new(int32(f.RawModTimeWindowS)),
		MaxConcurrentWrites:     new(int32(f.MaxConcurrentWrites)),
		DisableFsync:            new(f.DisableFsync),
		BlockPullOrder:          new(configpb.BlockPullOrder(f.BlockPullOrder)),
		CopyRangeMethod:         new(configpb.CopyRangeMethod(f.CopyRangeMethod)),
		CaseSensitiveFs:         new(f.CaseSensitiveFS),
		JunctionsAsDirs:         new(f.JunctionsAsDirs),
		SyncOwnership:           new(f.SyncOwnership),
		SendOwnership:           new(f.SendOwnership),
		SyncXattrs:              new(f.SyncXattrs),
		SendXattrs:              new(f.SendXattrs),
		BlockIndexing:           new(f.BlockIndexing),
		XattrFilter:             fromLegacyXattrFilter(f.XattrFilter),
	}.Build()
}

func fromLegacyFolderDevices(devices []config.FolderDeviceConfiguration) []*configpb.FolderDeviceConfiguration {
	out := make([]*configpb.FolderDeviceConfiguration, len(devices))
	for i, device := range devices {
		out[i] = configpb.FolderDeviceConfiguration_builder{
			DeviceId:           new(device.DeviceID.String()),
			IntroducedBy:       new(device.IntroducedBy.String()),
			EncryptionPassword: new(device.EncryptionPassword),
		}.Build()
	}
	return out
}

func fromLegacyVersioning(v config.VersioningConfiguration) *configpb.VersioningConfiguration {
	return configpb.VersioningConfiguration_builder{
		Type:             new(v.Type),
		Params:           maps.Clone(v.Params),
		CleanupIntervalS: new(int32(v.CleanupIntervalS)),
		FsPath:           new(v.FSPath),
		FsType:           new(filesystemTypeFromLegacy(v.FSType)),
	}.Build()
}

func fromLegacyXattrFilter(f config.XattrFilter) *configpb.XattrFilter {
	entries := make([]*configpb.XattrFilterEntry, len(f.Entries))
	for i, entry := range f.Entries {
		entries[i] = configpb.XattrFilterEntry_builder{
			Match:  new(entry.Match),
			Permit: new(entry.Permit),
		}.Build()
	}
	return configpb.XattrFilter_builder{
		Entries:            entries,
		MaxSingleEntrySize: new(int32(f.MaxSingleEntrySize)),
		MaxTotalSize:       new(int32(f.MaxTotalSize)),
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
		return configpb.Size_builder{Percent: new(s.Value)}.Build()
	}
	switch strings.ToLower(s.Unit) {
	case "mib":
		return configpb.Size_builder{Mib: new(s.Value)}.Build()
	case "gib":
		return configpb.Size_builder{Gib: new(s.Value)}.Build()
	case "kib":
		return configpb.Size_builder{Bytes: new(s.Value * 1024)}.Build()
	case "tib":
		return configpb.Size_builder{Bytes: new(s.Value * (1 << 40))}.Build()
	case "k", "kb":
		return configpb.Size_builder{Bytes: new(s.Value * 1000)}.Build()
	case "m", "mb":
		return configpb.Size_builder{Bytes: new(s.Value * 1000 * 1000)}.Build()
	case "g", "gb":
		return configpb.Size_builder{Bytes: new(s.Value * 1000 * 1000 * 1000)}.Build()
	case "t", "tb":
		return configpb.Size_builder{Bytes: new(s.Value * 1000 * 1000 * 1000 * 1000)}.Build()
	default:
		return configpb.Size_builder{Bytes: new(s.Value)}.Build()
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
		DeviceId:                 new(d.DeviceID.String()),
		Name:                     new(d.Name),
		Addresses:                slices.Clone(d.Addresses),
		Compression:              new(configpb.Compression(d.Compression)),
		CertName:                 new(d.CertName),
		Introducer:               new(d.Introducer),
		SkipIntroductionRemovals: new(d.SkipIntroductionRemovals),
		IntroducedBy:             new(d.IntroducedBy.String()),
		Paused:                   new(d.Paused),
		AllowedNetworks:          slices.Clone(d.AllowedNetworks),
		AutoAcceptFolders:        new(d.AutoAcceptFolders),
		MaxSendKbps:              new(int32(d.MaxSendKbps)),
		MaxRecvKbps:              new(int32(d.MaxRecvKbps)),
		IgnoredFolders:           fromLegacyObservedFolders(d.IgnoredFolders),
		MaxRequestKiB:            new(int32(d.MaxRequestKiB)),
		Untrusted:                new(d.Untrusted),
		RemoteGuiPort:            new(int32(d.RemoteGUIPort)),
		NumConnections:           new(int32(d.RawNumConnections)),
		Group:                    new(d.Group),
	}.Build()
}

func fromLegacyObservedFolders(folders []config.ObservedFolder) []*configpb.ObservedFolder {
	out := make([]*configpb.ObservedFolder, len(folders))
	for i, folder := range folders {
		out[i] = configpb.ObservedFolder_builder{
			Time:  timestamppb.New(folder.Time),
			Id:    new(folder.ID),
			Label: new(folder.Label),
		}.Build()
	}
	return out
}

func fromLegacyObservedDevices(devices []config.ObservedDevice) []*configpb.ObservedDevice {
	out := make([]*configpb.ObservedDevice, len(devices))
	for i, device := range devices {
		out[i] = configpb.ObservedDevice_builder{
			Time:     timestamppb.New(device.Time),
			DeviceId: new(device.ID.String()),
			Name:     new(device.Name),
			Address:  new(device.Address),
		}.Build()
	}
	return out
}

func fromLegacyGUI(g config.GUIConfiguration) *configpb.GUIConfiguration {
	return configpb.GUIConfiguration_builder{
		Enabled:                   new(g.Enabled),
		Address:                   new(g.RawAddress),
		UnixSocketPermissions:     new(g.RawUnixSocketPermissions),
		User:                      new(g.User),
		Password:                  new(g.Password),
		AuthMode:                  new(configpb.AuthMode(g.AuthMode)),
		MetricsWithoutAuth:        new(g.MetricsWithoutAuth),
		UseTls:                    new(g.RawUseTLS),
		ApiKey:                    new(g.APIKey),
		InsecureAdminAccess:       new(g.InsecureAdminAccess),
		Theme:                     new(g.Theme),
		InsecureSkipHostcheck:     new(g.InsecureSkipHostCheck),
		InsecureAllowFrameLoading: new(g.InsecureAllowFrameLoading),
		SendBasicAuthPrompt:       new(g.SendBasicAuthPrompt),
		SessionCookieDurationS:    new(int32(g.SessionCookieDurationS)),
		SessionCookiePath:         new(g.SessionCookiePath),
	}.Build()
}

func fromLegacyLDAP(l config.LDAPConfiguration) *configpb.LDAPConfiguration {
	return configpb.LDAPConfiguration_builder{
		Address:            new(l.Address),
		BindDn:             new(l.BindDN),
		Transport:          new(configpb.LDAPTransport(l.Transport)),
		InsecureSkipVerify: new(l.InsecureSkipVerify),
		SearchBaseDn:       new(l.SearchBaseDN),
		SearchFilter:       new(l.SearchFilter),
	}.Build()
}

func fromLegacyOptions(o config.OptionsConfiguration) *configpb.OptionsConfiguration {
	return configpb.OptionsConfiguration_builder{
		ListenAddresses:                     slices.Clone(o.RawListenAddresses),
		GlobalAnnounceServers:               slices.Clone(o.RawGlobalAnnServers),
		GlobalAnnounceEnabled:               new(o.GlobalAnnEnabled),
		LocalAnnounceEnabled:                new(o.LocalAnnEnabled),
		LocalAnnouncePort:                   new(int32(o.LocalAnnPort)),
		LocalAnnounceMcAddr:                 new(o.LocalAnnMCAddr),
		MaxSendKbps:                         new(int32(o.MaxSendKbps)),
		MaxRecvKbps:                         new(int32(o.MaxRecvKbps)),
		ReconnectionIntervalS:               new(int32(o.ReconnectIntervalS)),
		RelaysEnabled:                       new(o.RelaysEnabled),
		RelayReconnectIntervalM:             new(int32(o.RelayReconnectIntervalM)),
		StartBrowser:                        new(o.StartBrowser),
		NatEnabled:                          new(o.NATEnabled),
		NatLeaseMinutes:                     new(int32(o.NATLeaseM)),
		NatRenewalMinutes:                   new(int32(o.NATRenewalM)),
		NatTimeoutSeconds:                   new(int32(o.NATTimeoutS)),
		UrAccepted:                          new(int32(o.URAccepted)),
		UrSeen:                              new(int32(o.URSeen)),
		UrUniqueId:                          new(o.URUniqueID),
		UrUrl:                               new(o.URURL),
		UrPostInsecurely:                    new(o.URPostInsecurely),
		UrInitialDelayS:                     new(int32(o.URInitialDelayS)),
		AutoUpgradeIntervalH:                new(int32(o.AutoUpgradeIntervalH)),
		UpgradeToPreReleases:                new(o.UpgradeToPreReleases),
		KeepTemporariesH:                    new(int32(o.KeepTemporariesH)),
		ProgressUpdateIntervalS:             new(int32(o.ProgressUpdateIntervalS)),
		LimitBandwidthInLan:                 new(o.LimitBandwidthInLan),
		MinHomeDiskFree:                     fromLegacySize(o.MinHomeDiskFree),
		ReleasesUrl:                         new(o.ReleasesURL),
		AlwaysLocalNets:                     slices.Clone(o.AlwaysLocalNets),
		OverwriteRemoteDeviceNamesOnConnect: new(o.OverwriteRemoteDevNames),
		TempIndexMinBlocks:                  new(int32(o.TempIndexMinBlocks)),
		UnackedNotificationIds:              slices.Clone(o.UnackedNotificationIDs),
		TrafficClass:                        new(int32(o.TrafficClass)),
		SetLowPriority:                      new(o.SetLowPriority),
		MaxFolderConcurrency:                new(int32(o.RawMaxFolderConcurrency)),
		CrUrl:                               new(o.CRURL),
		CrashReportingEnabled:               new(o.CREnabled),
		StunKeepaliveStartS:                 new(int32(o.StunKeepaliveStartS)),
		StunKeepaliveMinS:                   new(int32(o.StunKeepaliveMinS)),
		StunServers:                         slices.Clone(o.RawStunServers),
		MaxConcurrentIncomingRequestKiB:     new(int32(o.RawMaxCIRequestKiB)),
		AnnounceLanAddresses:                new(o.AnnounceLANAddresses),
		SendFullIndexOnUpgrade:              new(o.SendFullIndexOnUpgrade),
		FeatureFlags:                        slices.Clone(o.FeatureFlags),
		AuditEnabled:                        new(o.AuditEnabled),
		AuditFile:                           new(o.AuditFile),
		ConnectionLimitEnough:               new(int32(o.ConnectionLimitEnough)),
		ConnectionLimitMax:                  new(int32(o.ConnectionLimitMax)),
		ConnectionPriorityTcpLan:            new(int32(o.ConnectionPriorityTCPLAN)),
		ConnectionPriorityQuicLan:           new(int32(o.ConnectionPriorityQUICLAN)),
		ConnectionPriorityTcpWan:            new(int32(o.ConnectionPriorityTCPWAN)),
		ConnectionPriorityQuicWan:           new(int32(o.ConnectionPriorityQUICWAN)),
		ConnectionPriorityRelay:             new(int32(o.ConnectionPriorityRelay)),
		ConnectionPriorityUpgradeThreshold:  new(int32(o.ConnectionPriorityUpgradeThreshold)),
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
