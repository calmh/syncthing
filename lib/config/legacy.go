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

	intconfig "github.com/syncthing/syncthing/internal/config"
	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
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
func FromLegacy(cfg Configuration) intconfig.Configuration {
	return intconfig.Configuration{Configuration: syncthingv2.Configuration_builder{
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

func fromLegacyFolders(folders []FolderConfiguration) []*syncthingv2.FolderConfiguration {
	out := make([]*syncthingv2.FolderConfiguration, len(folders))
	for i, folder := range folders {
		out[i] = fromLegacyFolder(folder)
	}
	return out
}

func fromLegacyFolder(f FolderConfiguration) *syncthingv2.FolderConfiguration {
	return syncthingv2.FolderConfiguration_builder{
		Id:                      new(f.ID),
		Label:                   new(f.Label),
		FilesystemType:          new(filesystemTypeFromLegacy(f.FilesystemType)),
		Path:                    new(f.Path),
		Type:                    new(syncthingv2.FolderType(f.Type)),
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
		Order:                   new(syncthingv2.PullOrder(f.Order)),
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
		BlockPullOrder:          new(syncthingv2.BlockPullOrder(f.BlockPullOrder)),
		CopyRangeMethod:         new(syncthingv2.CopyRangeMethod(f.CopyRangeMethod)),
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

func fromLegacyFolderDevices(devices []FolderDeviceConfiguration) []*syncthingv2.FolderDeviceConfiguration {
	out := make([]*syncthingv2.FolderDeviceConfiguration, len(devices))
	for i, device := range devices {
		out[i] = syncthingv2.FolderDeviceConfiguration_builder{
			DeviceId:           new(device.DeviceID.String()),
			IntroducedBy:       new(device.IntroducedBy.String()),
			EncryptionPassword: new(device.EncryptionPassword),
		}.Build()
	}
	return out
}

func fromLegacyVersioning(v VersioningConfiguration) *syncthingv2.VersioningConfiguration {
	return syncthingv2.VersioningConfiguration_builder{
		Type:             new(v.Type),
		Params:           maps.Clone(v.Params),
		CleanupIntervalS: new(int32(v.CleanupIntervalS)),
		FsPath:           new(v.FSPath),
		FsType:           new(filesystemTypeFromLegacy(v.FSType)),
	}.Build()
}

func fromLegacyXattrFilter(f XattrFilter) *syncthingv2.XattrFilter {
	entries := make([]*syncthingv2.XattrFilterEntry, len(f.Entries))
	for i, entry := range f.Entries {
		entries[i] = syncthingv2.XattrFilterEntry_builder{
			Match:  new(entry.Match),
			Permit: new(entry.Permit),
		}.Build()
	}
	return syncthingv2.XattrFilter_builder{
		Entries:            entries,
		MaxSingleEntrySize: new(int32(f.MaxSingleEntrySize)),
		MaxTotalSize:       new(int32(f.MaxTotalSize)),
	}.Build()
}

// fromLegacySize converts a legacy size to the new representation. A
// legacy size of zero or less meant "no minimum" and converts to an unset
// size. Units with an "iB"-suffix are binary, other unit prefixes decimal,
// matching the intent of the legacy string format.
func fromLegacySize(s Size) *syncthingv2.Size {
	if s.Value <= 0 {
		// A legacy size of zero or less disabled the check. Express it
		// as an explicit zero, as an unset size means the default.
		return syncthingv2.Size_builder{Percent: new(0.0)}.Build()
	}
	if strings.Contains(s.Unit, "%") {
		return syncthingv2.Size_builder{Percent: new(s.Value)}.Build()
	}
	switch strings.ToLower(s.Unit) {
	case "mib":
		return syncthingv2.Size_builder{Mib: new(s.Value)}.Build()
	case "gib":
		return syncthingv2.Size_builder{Gib: new(s.Value)}.Build()
	case "kib":
		return syncthingv2.Size_builder{Bytes: new(s.Value * 1024)}.Build()
	case "tib":
		return syncthingv2.Size_builder{Bytes: new(s.Value * (1 << 40))}.Build()
	case "k", "kb":
		return syncthingv2.Size_builder{Bytes: new(s.Value * 1000)}.Build()
	case "m", "mb":
		return syncthingv2.Size_builder{Bytes: new(s.Value * 1000 * 1000)}.Build()
	case "g", "gb":
		return syncthingv2.Size_builder{Bytes: new(s.Value * 1000 * 1000 * 1000)}.Build()
	case "t", "tb":
		return syncthingv2.Size_builder{Bytes: new(s.Value * 1000 * 1000 * 1000 * 1000)}.Build()
	default:
		return syncthingv2.Size_builder{Bytes: new(s.Value)}.Build()
	}
}

func fromLegacyDevices(devices []DeviceConfiguration) []*syncthingv2.DeviceConfiguration {
	out := make([]*syncthingv2.DeviceConfiguration, len(devices))
	for i, device := range devices {
		out[i] = fromLegacyDevice(device)
	}
	return out
}

func fromLegacyDevice(d DeviceConfiguration) *syncthingv2.DeviceConfiguration {
	return syncthingv2.DeviceConfiguration_builder{
		DeviceId:                 new(d.DeviceID.String()),
		Name:                     new(d.Name),
		Addresses:                slices.Clone(d.Addresses),
		Compression:              new(syncthingv2.Compression(d.Compression)),
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

func fromLegacyObservedFolders(folders []ObservedFolder) []*syncthingv2.ObservedFolder {
	out := make([]*syncthingv2.ObservedFolder, len(folders))
	for i, folder := range folders {
		out[i] = syncthingv2.ObservedFolder_builder{
			Time:  timestamppb.New(folder.Time),
			Id:    new(folder.ID),
			Label: new(folder.Label),
		}.Build()
	}
	return out
}

func fromLegacyObservedDevices(devices []ObservedDevice) []*syncthingv2.ObservedDevice {
	out := make([]*syncthingv2.ObservedDevice, len(devices))
	for i, device := range devices {
		out[i] = syncthingv2.ObservedDevice_builder{
			Time:     timestamppb.New(device.Time),
			DeviceId: new(device.ID.String()),
			Name:     new(device.Name),
			Address:  new(device.Address),
		}.Build()
	}
	return out
}

func fromLegacyGUI(g GUIConfiguration) *syncthingv2.GUIConfiguration {
	return syncthingv2.GUIConfiguration_builder{
		Enabled:                   new(g.Enabled),
		Address:                   new(g.RawAddress),
		UnixSocketPermissions:     new(g.RawUnixSocketPermissions),
		User:                      new(g.User),
		Password:                  new(g.Password),
		AuthMode:                  new(syncthingv2.AuthMode(g.AuthMode)),
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

func fromLegacyLDAP(l LDAPConfiguration) *syncthingv2.LDAPConfiguration {
	return syncthingv2.LDAPConfiguration_builder{
		Address:            new(l.Address),
		BindDn:             new(l.BindDN),
		Transport:          new(syncthingv2.LDAPTransport(l.Transport)),
		InsecureSkipVerify: new(l.InsecureSkipVerify),
		SearchBaseDn:       new(l.SearchBaseDN),
		SearchFilter:       new(l.SearchFilter),
	}.Build()
}

func fromLegacyOptions(o OptionsConfiguration) *syncthingv2.OptionsConfiguration {
	return syncthingv2.OptionsConfiguration_builder{
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

func fromLegacyDefaults(d Defaults) *syncthingv2.Defaults {
	return syncthingv2.Defaults_builder{
		Folder:  fromLegacyFolder(d.Folder),
		Device:  fromLegacyDevice(d.Device),
		Ignores: syncthingv2.Ignores_builder{Lines: slices.Clone(d.Ignores.Lines)}.Build(),
	}.Build()
}

// filesystemTypeFromLegacy maps the legacy string based filesystem type to
// the new enum. The legacy zero value ("") means "basic".
func filesystemTypeFromLegacy(t FilesystemType) syncthingv2.FilesystemType {
	if t == FilesystemTypeFake {
		return syncthingv2.FilesystemType_FILESYSTEM_TYPE_FAKE
	}
	return syncthingv2.FilesystemType_FILESYSTEM_TYPE_BASIC
}

// ToLegacy returns a legacy format configuration with the same effective
// contents as the given configuration. Values are read through the
// getters, so fields that are unset in the new format get their default
// values, like a loaded legacy configuration would have.
func ToLegacy(cfg intconfig.Configuration) Configuration {
	return Configuration{
		Version:        int(cfg.GetVersion()),
		Folders:        toLegacyFolders(cfg.GetFolders()),
		Devices:        toLegacyDevices(cfg.GetDevices()),
		GUI:            toLegacyGUI(cfg.GetGui()),
		LDAP:           toLegacyLDAP(cfg.GetLdap()),
		Options:        toLegacyOptions(cfg.GetOptions()),
		IgnoredDevices: toLegacyObservedDevices(cfg.GetRemoteIgnoredDevices()),
		Defaults:       toLegacyDefaults(cfg.GetDefaults()),
	}
}

func toLegacyFolders(folders []intconfig.FolderConfiguration) []FolderConfiguration {
	if len(folders) == 0 {
		return nil
	}
	out := make([]FolderConfiguration, len(folders))
	for i, folder := range folders {
		out[i] = toLegacyFolder(folder)
	}
	return out
}

func toLegacyFolder(f intconfig.FolderConfiguration) FolderConfiguration {
	return FolderConfiguration{
		ID:                      f.GetId(),
		Label:                   f.GetLabel(),
		FilesystemType:          filesystemTypeToLegacy(f.GetFilesystemType()),
		Path:                    f.GetPath(),
		Type:                    FolderType(f.GetType()),
		Devices:                 toLegacyFolderDevices(f.GetDevices()),
		Group:                   f.GetGroup(),
		RescanIntervalS:         int(f.GetRescanIntervalS()),
		FSWatcherEnabled:        f.GetFsWatcherEnabled(),
		FSWatcherDelayS:         f.GetFsWatcherDelayS(),
		FSWatcherTimeoutS:       f.GetFsWatcherTimeoutS(),
		IgnorePerms:             f.GetIgnorePerms(),
		AutoNormalize:           f.GetAutoNormalize(),
		MinDiskFree:             toLegacySize(f.MinDiskFree()),
		Versioning:              toLegacyVersioning(f.GetVersioning()),
		Copiers:                 int(f.GetCopiers()),
		PullerMaxPendingKiB:     int(f.GetPullerMaxPendingKiB()),
		Hashers:                 int(f.GetHashers()),
		Order:                   PullOrder(f.GetOrder()),
		IgnoreDelete:            f.GetIgnoreDelete(),
		ScanProgressIntervalS:   int(f.GetScanProgressIntervalS()),
		PullerPauseS:            int(f.GetPullerPauseS()),
		PullerDelayS:            f.GetPullerDelayS(),
		MaxConflicts:            int(f.GetMaxConflicts()),
		DisableSparseFiles:      f.GetDisableSparseFiles(),
		Paused:                  f.GetPaused(),
		MarkerName:              f.GetMarkerName(),
		CopyOwnershipFromParent: f.GetCopyOwnershipFromParent(),
		RawModTimeWindowS:       int(f.GetModTimeWindowS()),
		MaxConcurrentWrites:     int(f.GetMaxConcurrentWrites()),
		DisableFsync:            f.GetDisableFsync(),
		BlockPullOrder:          BlockPullOrder(f.GetBlockPullOrder()),
		CopyRangeMethod:         CopyRangeMethod(f.GetCopyRangeMethod()),
		CaseSensitiveFS:         f.GetCaseSensitiveFs(),
		JunctionsAsDirs:         f.GetJunctionsAsDirs(),
		SyncOwnership:           f.GetSyncOwnership(),
		SendOwnership:           f.GetSendOwnership(),
		SyncXattrs:              f.GetSyncXattrs(),
		SendXattrs:              f.GetSendXattrs(),
		BlockIndexing:           f.GetBlockIndexing(),
		XattrFilter:             toLegacyXattrFilter(f.GetXattrFilter()),
	}
}

func toLegacyFolderDevices(devices []intconfig.FolderDeviceConfiguration) []FolderDeviceConfiguration {
	if len(devices) == 0 {
		return nil
	}
	out := make([]FolderDeviceConfiguration, len(devices))
	for i, device := range devices {
		out[i] = FolderDeviceConfiguration{
			DeviceID:           device.GetDeviceId(),
			IntroducedBy:       device.GetIntroducedBy(),
			EncryptionPassword: device.GetEncryptionPassword(),
		}
	}
	return out
}

func toLegacyVersioning(v *syncthingv2.VersioningConfiguration) VersioningConfiguration {
	return VersioningConfiguration{
		Type:             v.GetType(),
		Params:           maps.Clone(v.GetParams()),
		CleanupIntervalS: int(v.GetCleanupIntervalS()),
		FSPath:           v.GetFsPath(),
		FSType:           filesystemTypeToLegacy(v.GetFsType()),
	}
}

func toLegacyXattrFilter(f *syncthingv2.XattrFilter) XattrFilter {
	var entries []XattrFilterEntry
	if len(f.GetEntries()) > 0 {
		entries = make([]XattrFilterEntry, len(f.GetEntries()))
		for i, entry := range f.GetEntries() {
			entries[i] = XattrFilterEntry{
				Match:  entry.GetMatch(),
				Permit: entry.GetPermit(),
			}
		}
	}
	return XattrFilter{
		Entries:            entries,
		MaxSingleEntrySize: int(f.GetMaxSingleEntrySize()),
		MaxTotalSize:       int(f.GetMaxTotalSize()),
	}
}

// toLegacySize converts a size to the legacy value and unit form. A size
// with a zero value is the legacy zero size, disabling the check.
func toLegacySize(size *intconfig.Size) Size {
	var s Size
	switch {
	case size.HasPercent():
		s = Size{Value: size.GetPercent(), Unit: "%"}
	case size.HasMib():
		s = Size{Value: size.GetMib(), Unit: "MiB"}
	case size.HasGib():
		s = Size{Value: size.GetGib(), Unit: "GiB"}
	case size.HasBytes():
		s = Size{Value: size.GetBytes(), Unit: "B"}
	}
	if s.Value == 0 {
		// Zero disables the check; the legacy zero size has no unit.
		return Size{}
	}
	return s
}

func toLegacyDevices(devices []intconfig.DeviceConfiguration) []DeviceConfiguration {
	if len(devices) == 0 {
		return nil
	}
	out := make([]DeviceConfiguration, len(devices))
	for i, device := range devices {
		out[i] = toLegacyDevice(device)
	}
	return out
}

func toLegacyDevice(d intconfig.DeviceConfiguration) DeviceConfiguration {
	return DeviceConfiguration{
		DeviceID:                 d.GetDeviceId(),
		Name:                     d.GetName(),
		Addresses:                slices.Clone(d.Addresses()),
		Compression:              Compression(d.GetCompression()),
		CertName:                 d.GetCertName(),
		Introducer:               d.GetIntroducer(),
		SkipIntroductionRemovals: d.GetSkipIntroductionRemovals(),
		IntroducedBy:             d.GetIntroducedBy(),
		Paused:                   d.GetPaused(),
		AllowedNetworks:          slices.Clone(d.GetAllowedNetworks()),
		AutoAcceptFolders:        d.GetAutoAcceptFolders(),
		MaxSendKbps:              int(d.GetMaxSendKbps()),
		MaxRecvKbps:              int(d.GetMaxRecvKbps()),
		IgnoredFolders:           toLegacyObservedFolders(d.GetIgnoredFolders()),
		MaxRequestKiB:            int(d.GetMaxRequestKiB()),
		Untrusted:                d.GetUntrusted(),
		RemoteGUIPort:            int(d.GetRemoteGuiPort()),
		RawNumConnections:        int(d.GetNumConnections()),
		Group:                    d.GetGroup(),
	}
}

func toLegacyObservedFolders(folders []*syncthingv2.ObservedFolder) []ObservedFolder {
	if len(folders) == 0 {
		return nil
	}
	out := make([]ObservedFolder, len(folders))
	for i, folder := range folders {
		out[i] = ObservedFolder{
			Time:  folder.GetTime().AsTime(),
			ID:    folder.GetId(),
			Label: folder.GetLabel(),
		}
	}
	return out
}

func toLegacyObservedDevices(devices []intconfig.ObservedDevice) []ObservedDevice {
	if len(devices) == 0 {
		return nil
	}
	out := make([]ObservedDevice, len(devices))
	for i, device := range devices {
		out[i] = ObservedDevice{
			Time:    device.GetTime().AsTime(),
			ID:      device.GetDeviceId(),
			Name:    device.GetName(),
			Address: device.GetAddress(),
		}
	}
	return out
}

func toLegacyGUI(g *syncthingv2.GUIConfiguration) GUIConfiguration {
	return GUIConfiguration{
		Enabled:                   g.GetEnabled(),
		RawAddress:                g.GetAddress(),
		RawUnixSocketPermissions:  g.GetUnixSocketPermissions(),
		User:                      g.GetUser(),
		Password:                  g.GetPassword(),
		AuthMode:                  AuthMode(g.GetAuthMode()),
		MetricsWithoutAuth:        g.GetMetricsWithoutAuth(),
		RawUseTLS:                 g.GetUseTls(),
		APIKey:                    g.GetApiKey(),
		InsecureAdminAccess:       g.GetInsecureAdminAccess(),
		Theme:                     g.GetTheme(),
		InsecureSkipHostCheck:     g.GetInsecureSkipHostcheck(),
		InsecureAllowFrameLoading: g.GetInsecureAllowFrameLoading(),
		SendBasicAuthPrompt:       g.GetSendBasicAuthPrompt(),
		SessionCookieDurationS:    int(g.GetSessionCookieDurationS()),
		SessionCookiePath:         g.GetSessionCookiePath(),
	}
}

func toLegacyLDAP(l *syncthingv2.LDAPConfiguration) LDAPConfiguration {
	return LDAPConfiguration{
		Address:            l.GetAddress(),
		BindDN:             l.GetBindDn(),
		Transport:          LDAPTransport(l.GetTransport()),
		InsecureSkipVerify: l.GetInsecureSkipVerify(),
		SearchBaseDN:       l.GetSearchBaseDn(),
		SearchFilter:       l.GetSearchFilter(),
	}
}

func toLegacyOptions(o intconfig.OptionsConfiguration) OptionsConfiguration {
	return OptionsConfiguration{
		RawListenAddresses:                 slices.Clone(o.ListenAddresses()),
		RawGlobalAnnServers:                slices.Clone(o.GlobalAnnounceServers()),
		GlobalAnnEnabled:                   o.GetGlobalAnnounceEnabled(),
		LocalAnnEnabled:                    o.GetLocalAnnounceEnabled(),
		LocalAnnPort:                       int(o.GetLocalAnnouncePort()),
		LocalAnnMCAddr:                     o.GetLocalAnnounceMcAddr(),
		MaxSendKbps:                        int(o.GetMaxSendKbps()),
		MaxRecvKbps:                        int(o.GetMaxRecvKbps()),
		ReconnectIntervalS:                 int(o.GetReconnectionIntervalS()),
		RelaysEnabled:                      o.GetRelaysEnabled(),
		RelayReconnectIntervalM:            int(o.GetRelayReconnectIntervalM()),
		StartBrowser:                       o.GetStartBrowser(),
		NATEnabled:                         o.GetNatEnabled(),
		NATLeaseM:                          int(o.GetNatLeaseMinutes()),
		NATRenewalM:                        int(o.GetNatRenewalMinutes()),
		NATTimeoutS:                        int(o.GetNatTimeoutSeconds()),
		URAccepted:                         int(o.GetUrAccepted()),
		URSeen:                             int(o.GetUrSeen()),
		URUniqueID:                         o.GetUrUniqueId(),
		URURL:                              o.GetUrUrl(),
		URPostInsecurely:                   o.GetUrPostInsecurely(),
		URInitialDelayS:                    int(o.GetUrInitialDelayS()),
		AutoUpgradeIntervalH:               int(o.GetAutoUpgradeIntervalH()),
		UpgradeToPreReleases:               o.GetUpgradeToPreReleases(),
		KeepTemporariesH:                   int(o.GetKeepTemporariesH()),
		ProgressUpdateIntervalS:            int(o.GetProgressUpdateIntervalS()),
		LimitBandwidthInLan:                o.GetLimitBandwidthInLan(),
		MinHomeDiskFree:                    toLegacySize(o.MinHomeDiskFree()),
		ReleasesURL:                        o.GetReleasesUrl(),
		AlwaysLocalNets:                    slices.Clone(o.GetAlwaysLocalNets()),
		OverwriteRemoteDevNames:            o.GetOverwriteRemoteDeviceNamesOnConnect(),
		TempIndexMinBlocks:                 int(o.GetTempIndexMinBlocks()),
		UnackedNotificationIDs:             slices.Clone(o.GetUnackedNotificationIds()),
		TrafficClass:                       int(o.GetTrafficClass()),
		SetLowPriority:                     o.GetSetLowPriority(),
		RawMaxFolderConcurrency:            int(o.GetMaxFolderConcurrency()),
		CRURL:                              o.GetCrUrl(),
		CREnabled:                          o.GetCrashReportingEnabled(),
		StunKeepaliveStartS:                int(o.GetStunKeepaliveStartS()),
		StunKeepaliveMinS:                  int(o.GetStunKeepaliveMinS()),
		RawStunServers:                     slices.Clone(o.StunServers()),
		RawMaxCIRequestKiB:                 int(o.GetMaxConcurrentIncomingRequestKiB()),
		AnnounceLANAddresses:               o.GetAnnounceLanAddresses(),
		SendFullIndexOnUpgrade:             o.GetSendFullIndexOnUpgrade(),
		FeatureFlags:                       slices.Clone(o.GetFeatureFlags()),
		AuditEnabled:                       o.GetAuditEnabled(),
		AuditFile:                          o.GetAuditFile(),
		ConnectionLimitEnough:              int(o.GetConnectionLimitEnough()),
		ConnectionLimitMax:                 int(o.GetConnectionLimitMax()),
		ConnectionPriorityTCPLAN:           int(o.GetConnectionPriorityTcpLan()),
		ConnectionPriorityQUICLAN:          int(o.GetConnectionPriorityQuicLan()),
		ConnectionPriorityTCPWAN:           int(o.GetConnectionPriorityTcpWan()),
		ConnectionPriorityQUICWAN:          int(o.GetConnectionPriorityQuicWan()),
		ConnectionPriorityRelay:            int(o.GetConnectionPriorityRelay()),
		ConnectionPriorityUpgradeThreshold: int(o.GetConnectionPriorityUpgradeThreshold()),
	}
}

func toLegacyDefaults(d intconfig.Defaults) Defaults {
	return Defaults{
		Folder:  toLegacyFolder(d.GetFolder()),
		Device:  toLegacyDevice(d.GetDevice()),
		Ignores: toLegacyIgnores(d.GetIgnores()),
	}
}

func toLegacyIgnores(i *syncthingv2.Ignores) Ignores {
	return Ignores{Lines: slices.Clone(i.GetLines())}
}

func filesystemTypeToLegacy(t syncthingv2.FilesystemType) FilesystemType {
	if t == syncthingv2.FilesystemType_FILESYSTEM_TYPE_FAKE {
		return FilesystemTypeFake
	}
	return FilesystemTypeBasic
}
