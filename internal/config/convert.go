package config

import (
	"cmp"
	"slices"

	configv1 "github.com/syncthing/syncthing/internal/config/v1"
	configv2 "github.com/syncthing/syncthing/internal/config/v2"
)

func FromV1(v1 *configv1.Configuration) *configv2.Configuration {
	folders := make([]*configv2.FolderConfiguration, len(v1.Folders))
	for i, fld := range v1.Folders {
		folders[i] = v2folder(fld)
	}

	devices := make([]*configv2.DeviceConfiguration, len(v1.Devices))
	for i, dev := range v1.Devices {
		devices[i] = v2Device(dev)
	}

	return configv2.Configuration_builder{
		Folders: folders,
		Devices: devices,
		Options: v2Options(v1.Options),
	}.Build()
}

func v2folder(fld configv1.FolderConfiguration) *configv2.FolderConfiguration {
	var (
		zf   = (*configv2.FolderConfiguration)(nil)
		zfso = (*configv2.FolderConfiguration_FilesystemOptions)(nil)
		zsc  = (*configv2.FolderConfiguration_Scanning)(nil)
		zpu  = (*configv2.FolderConfiguration_Pulling)(nil)
		zwa  = (*configv2.FolderConfiguration_Scanning_Watcher)(nil)
	)

	bld := configv2.FolderConfiguration_builder{
		Type:       v2FolderType(fld.Type),
		FolderId:   ptr(fld.ID),
		Enabled:    ptr(!fld.Paused, zf.GetEnabled()),
		Label:      ptr(fld.Label, ""),
		MarkerName: ptr(fld.MarkerName, zf.GetMarkerName()),

		Filesystem: configv2.FolderConfiguration_Filesystem_builder{
			Path: ptr(fld.Path),
			Type: v2FSType(fld.FilesystemType),
		}.Build(),

		FilesystemOptions: configv2.FolderConfiguration_FilesystemOptions_builder{
			CaseSensitive:   ptr(fld.CaseSensitiveFS, zfso.GetCaseSensitive()),
			JunctionsAsDirs: ptr(fld.JunctionsAsDirs, zfso.GetJunctionsAsDirs()),
			MinDiskFree:     v2Size(fld.MinDiskFree),
			AutoNormalize:   ptr(fld.AutoNormalize, zfso.GetAutoNormalize()),
		}.Build(),

		Scanning: configv2.FolderConfiguration_Scanning_builder{
			Watcher: configv2.FolderConfiguration_Scanning_Watcher_builder{
				Enabled:  ptr(fld.FSWatcherEnabled, zwa.GetEnabled()),
				DelayS:   ptr(fld.FSWatcherDelayS, zwa.GetDelayS()),
				TimeoutS: ptr(fld.FSWatcherTimeoutS, zwa.GetTimeoutS()),
			}.Build(),
			RescanIntervalS:   ptr(int32(fld.RescanIntervalS), zsc.GetRescanIntervalS()),
			NumHashers:        ptr(int32(fld.Hashers), zsc.GetNumHashers()),
			ModTimeWindowS:    ptr(int32(fld.RawModTimeWindowS), zsc.GetModTimeWindowS()),
			ProgressIntervalS: ptr(int32(fld.ScanProgressIntervalS), zsc.GetProgressIntervalS()),
			ScanOwnership:     ptr(fld.SendOwnership, zsc.GetScanOwnership()),
			ScanXattrs:        ptr(fld.SendXattrs, zsc.GetScanXattrs()),
			XattrFilter:       configv2.XattrFilter_builder{}.Build(),
		}.Build(),

		Pulling: configv2.FolderConfiguration_Pulling_builder{
			Order:                   v2PullOrder(fld.Order),
			BlockOrder:              v2BlockPullOrder(fld.BlockPullOrder),
			NumCopiers:              ptr(int32(fld.Copiers), zpu.GetNumCopiers(), 0),
			MaxPendingKib:           ptr(int32(fld.PullerMaxPendingKiB), zpu.GetMaxPendingKib(), 0),
			MaxConflicts:            ptr(int32(fld.MaxConflicts), zpu.GetMaxConflicts()),
			MaxConcurrentWrites:     ptr(int32(fld.MaxConcurrentWrites), zpu.GetMaxConcurrentWrites(), 0),
			FailedPauseS:            ptr(int32(fld.PullerPauseS), zpu.GetFailedPauseS(), 0),
			ChangeDelayS:            ptr(fld.PullerDelayS, zpu.GetChangeDelayS()),
			SparseFiles:             ptr(!fld.DisableSparseFiles, zpu.GetSparseFiles()),
			IgnorePermissions:       ptr(fld.IgnorePerms, zpu.GetIgnorePermissions()),
			TemporaryIndexes:        ptr(!fld.DisableTempIndexes, zpu.GetTemporaryIndexes()),
			Fsync:                   ptr(!fld.DisableFsync, zpu.GetFsync()),
			CopyRangeMethod:         v2CopyRangeMethod(fld.CopyRangeMethod),
			SyncOwnership:           ptr(fld.SyncOwnership, zpu.GetSyncOwnership()),
			SyncXattrs:              ptr(fld.SyncXattrs, zpu.GetSyncXattrs()),
			IgnoreDelete:            ptr(fld.IgnoreDelete, zpu.GetIgnoreDelete()),
			CopyOwnershipFromParent: ptr(fld.CopyOwnershipFromParent, zpu.GetCopyOwnershipFromParent()),
		}.Build(),
	}

	for _, dev := range fld.Devices {
		bld.SharedWith = append(bld.SharedWith, configv2.FolderConfiguration_Sharing_builder{
			DeviceId:           ptr(dev.DeviceID.String()),
			IntroducedBy:       ptr(dev.IntroducedBy.String(), ""),
			EncryptionPassword: ptr(dev.EncryptionPassword, ""),
		}.Build())
	}

	if vt := v2VersioningType(fld.Versioning.Type); vt != nil {
		zv := (*configv2.FolderConfiguration_Versioning)(nil)
		vb := configv2.FolderConfiguration_Versioning_builder{
			Type: vt,
			Filesystem: configv2.FolderConfiguration_Filesystem_builder{
				Path: ptr(fld.Versioning.FSPath, ""),
				Type: v2FSType(fld.Versioning.FSType),
			}.Build(),
			CleanupIntervalS: ptr(int32(fld.Versioning.CleanupIntervalS), zv.GetCleanupIntervalS()),
		}
		for k, v := range fld.Versioning.Params {
			vb.Params = append(vb.Params, configv2.KV_builder{
				Key:   ptr(k),
				Value: ptr(v),
			}.Build())
		}
		slices.SortFunc(vb.Params, func(a, b *configv2.KV) int { return cmp.Compare(a.GetKey(), b.GetKey()) })
		bld.Versioning = vb.Build()
	}

	return bld.Build()
}

func v2FolderType(v1 configv1.FolderType) *configv2.FolderType {
	switch v1 {
	case configv1.FolderTypeSendReceive:
		return configv2.FolderType_FOLDER_TYPE_SEND_RECEIVE.Enum()
	case configv1.FolderTypeSendOnly:
		return configv2.FolderType_FOLDER_TYPE_SEND_ONLY.Enum()
	case configv1.FolderTypeReceiveOnly:
		return configv2.FolderType_FOLDER_TYPE_RECEIVE_ONLY.Enum()
	case configv1.FolderTypeReceiveEncrypted:
		return configv2.FolderType_FOLDER_TYPE_RECEIVE_ENCRYPTED.Enum()
	default:
		return nil
	}
}

func v2FSType(v1 configv1.FilesystemType) *configv2.FilesystemType {
	switch v1 {
	case configv1.FilesystemTypeBasic:
		return nil
	case configv1.FilesystemTypeFake:
		return configv2.FilesystemType_FILESYSTEM_TYPE_FAKE.Enum()
	default:
		return nil
	}
}

func v2Size(v1 configv1.Size) *configv2.Size {
	var bld configv2.Size_builder
	if v1.Unit == "%" {
		bld.Percent = &v1.Value
	} else {
		bld.Bytes = ptr(int64(v1.BaseValue()))
	}
	return bld.Build()
}

func v2PullOrder(v1 configv1.PullOrder) *configv2.PullOrder {
	switch v1 {
	case configv1.PullOrderRandom:
		return nil
	case configv1.PullOrderAlphabetic:
		return configv2.PullOrder_PULL_ORDER_ALPHABETIC.Enum()
	case configv1.PullOrderSmallestFirst:
		return configv2.PullOrder_PULL_ORDER_SMALLEST_FIRST.Enum()
	case configv1.PullOrderLargestFirst:
		return configv2.PullOrder_PULL_ORDER_LARGEST_FIRST.Enum()
	case configv1.PullOrderOldestFirst:
		return configv2.PullOrder_PULL_ORDER_OLDEST_FIRST.Enum()
	case configv1.PullOrderNewestFirst:
		return configv2.PullOrder_PULL_ORDER_NEWEST_FIRST.Enum()
	default:
		return nil
	}
}

func v2BlockPullOrder(v1 configv1.BlockPullOrder) *configv2.BlockPullOrder {
	switch v1 {
	case configv1.BlockPullOrderStandard:
		return nil
	case configv1.BlockPullOrderRandom:
		return configv2.BlockPullOrder_BLOCK_PULL_ORDER_RANDOM.Enum()
	case configv1.BlockPullOrderInOrder:
		return configv2.BlockPullOrder_BLOCK_PULL_ORDER_IN_ORDER.Enum()
	default:
		return nil
	}
}

func v2CopyRangeMethod(v1 configv1.CopyRangeMethod) *configv2.CopyRangeMethod {
	switch v1 {
	case configv1.CopyRangeMethodStandard:
		return nil
	case configv1.CopyRangeMethodIoctl:
		return configv2.CopyRangeMethod_COPY_RANGE_METHOD_IOCTL.Enum()
	case configv1.CopyRangeMethodCopyFileRange:
		return configv2.CopyRangeMethod_COPY_RANGE_METHOD_COPY_FILE_RANGE.Enum()
	case configv1.CopyRangeMethodSendFile:
		return configv2.CopyRangeMethod_COPY_RANGE_METHOD_SEND_FILE.Enum()
	case configv1.CopyRangeMethodDuplicateExtents:
		return configv2.CopyRangeMethod_COPY_RANGE_METHOD_DUPLICATE_EXTENTS.Enum()
	case configv1.CopyRangeMethodAllWithFallback:
		return configv2.CopyRangeMethod_COPY_RANGE_METHOD_ALL_WITH_FALLBACK.Enum()
	default:
		return nil
	}
}

func v2VersioningType(v1 string) *configv2.VersioningType {
	switch v1 {
	case "simple":
		return configv2.VersioningType_VERSIONING_TYPE_SIMPLE.Enum()
	case "trashcan":
		return configv2.VersioningType_VERSIONING_TYPE_TRASHCAN.Enum()
	case "staggered":
		return configv2.VersioningType_VERSIONING_TYPE_STAGGERED.Enum()
	case "external":
		return configv2.VersioningType_VERSIONING_TYPE_EXTERNAL.Enum()
	default:
		return nil
	}
}

func v2Device(v1 configv1.DeviceConfiguration) *configv2.DeviceConfiguration {
	zd := (*configv2.DeviceConfiguration)(nil)
	zr := (*configv2.RateLimits)(nil)

	bld := configv2.DeviceConfiguration_builder{
		DeviceId:                 ptr(v1.DeviceID.String()),
		Enabled:                  ptr(!v1.Paused, zd.GetEnabled()),
		Name:                     ptr(v1.Name, zd.GetName()),
		Compression:              v2Compression(v1.Compression),
		CertificateCommonName:    ptr(v1.CertName, zd.GetCertificateCommonName(), ""),
		Introducer:               ptr(v1.Introducer, zd.GetIntroducer()),
		SkipIntroductionRemovals: ptr(v1.SkipIntroductionRemovals, zd.GetSkipIntroductionRemovals()),
		IntroducedBy:             ptr(v1.IntroducedBy.String(), zd.GetIntroducedBy()),
		AutoAcceptFolders:        ptr(v1.AutoAcceptFolders, zd.GetAutoAcceptFolders()),
		MaxRequestKib:            ptr(int32(v1.MaxRequestKiB), zd.GetMaxRequestKib(), 0),
		Untrusted:                ptr(v1.Untrusted, zd.GetUntrusted()),
		RemoteGuiPort:            ptr(int32(v1.RemoteGUIPort), zd.GetRemoteGuiPort()),
		NumConnections:           ptr(int32(v1.NumConnections()), zd.GetNumConnections(), 0),
	}
	if v1.MaxSendKbps != 0 || v1.MaxRecvKbps != 0 {
		bld.RateLimits = configv2.RateLimits_builder{
			MaxSendKbps: ptr(int32(v1.MaxSendKbps), zr.GetMaxSendKbps()),
			MaxRecvKbps: ptr(int32(v1.MaxRecvKbps), zr.GetMaxRecvKbps()),
		}.Build()
	}

	return bld.Build()
}

func v2Compression(v1 configv1.Compression) *configv2.Compression {
	switch v1 {
	case configv1.CompressionMetadata:
		return nil
	case configv1.CompressionNever:
		return configv2.Compression_COMPRESSION_NEVER.Enum()
	case configv1.CompressionAlways:
		return configv2.Compression_COMPRESSION_ALWAYS.Enum()
	default:
		return nil
	}
}

func v2Options(v1 configv1.OptionsConfiguration) *configv2.OptionsConfiguration {
	ze := (*configv2.OptionsConfiguration_General)(nil)
	zn := (*configv2.OptionsConfiguration_Network)(nil)
	zcp := (*configv2.OptionsConfiguration_Network_ConnectionPriorities)(nil)
	zna := (*configv2.OptionsConfiguration_Network_NAT)(nil)
	zl := (*configv2.OptionsConfiguration_Network_Discovery_Local)(nil)
	zg := (*configv2.OptionsConfiguration_Network_Discovery_Global)(nil)
	zr := (*configv2.OptionsConfiguration_Network_Relays)(nil)
	zu := (*configv2.OptionsConfiguration_UsageReporting)(nil)
	zau := (*configv2.OptionsConfiguration_AutoUpgrade)(nil)

	bld := configv2.OptionsConfiguration_builder{
		General: configv2.OptionsConfiguration_General_builder{
			// listen
			StartBrowser:         ptr(v1.StartBrowser, ze.GetStartBrowser()),
			KeepTemporariesS:     ptr(int32(v1.KeepTemporariesH*3600), ze.GetKeepTemporariesS(), 0),
			CacheIgnoredFiles:    ptr(v1.CacheIgnoredFiles, ze.GetCacheIgnoredFiles()),
			OverwriteRemoteNames: ptr(v1.OverwriteRemoteDevNames, ze.GetOverwriteRemoteNames()),
			UnackedNotifications: v1.UnackedNotificationIDs,
			SetLowPriority:       ptr(v1.SetLowPriority, ze.GetSetLowPriority()),
			MaxFolderConcurrency: ptr(int32(v1.MaxFolderConcurrency()), ze.GetMaxFolderConcurrency()),
		}.Build(),

		Network: configv2.OptionsConfiguration_Network_builder{
			RateLimitLanConnections:         ptr(v1.LimitBandwidthInLan, zn.GetRateLimitLanConnections()),
			ReconnectIntervalS:              ptr(int32(v1.ReconnectIntervalS), zn.GetReconnectIntervalS()),
			AlwaysLocalNetworks:             v1.AlwaysLocalNets,
			TrafficClass:                    ptr(int32(v1.TrafficClass), zn.GetTrafficClass()),
			MaxConcurrentIncomingRequestKib: ptr(int32(v1.MaxConcurrentIncomingRequestKiB()), zn.GetMaxConcurrentIncomingRequestKib(), 0),
			ProgressUpdateIntervalS:         ptr(int32(v1.ProgressUpdateIntervalS), zn.GetProgressUpdateIntervalS()),
			TempIndexMinBlocks:              ptr(int32(v1.TempIndexMinBlocks), zn.GetTempIndexMinBlocks()),

			RateLimits: configv2.RateLimits_builder{
				MaxSendKbps: ptr(int32(v1.MaxSendKbps), 0),
				MaxRecvKbps: ptr(int32(v1.MaxRecvKbps), 0),
			}.Build(),

			ConnectionLimits: configv2.OptionsConfiguration_Network_ConnectionLimits_builder{
				Enough: ptr(int32(v1.ConnectionLimitEnough), 0),
				Max:    ptr(int32(v1.ConnectionLimitMax), 0),
			}.Build(),

			ConnectionPriorities: configv2.OptionsConfiguration_Network_ConnectionPriorities_builder{
				TcpLan:           ptr(int32(v1.ConnectionPriorityTCPLAN), zcp.GetTcpLan()),
				QuicLan:          ptr(int32(v1.ConnectionPriorityQUICLAN), zcp.GetQuicLan()),
				TcpWan:           ptr(int32(v1.ConnectionPriorityTCPWAN), zcp.GetTcpWan()),
				QuicWan:          ptr(int32(v1.ConnectionPriorityQUICWAN), zcp.GetQuicWan()),
				Relays:           ptr(int32(v1.ConnectionPriorityRelay), zcp.GetRelays()),
				UpgradeThreshold: ptr(int32(v1.ConnectionPriorityUpgradeThreshold), zcp.GetUpgradeThreshold()),
			}.Build(),

			Nat: configv2.OptionsConfiguration_Network_NAT_builder{
				Enabled:             ptr(v1.NATEnabled, zna.GetEnabled()),
				LeaseIntervalS:      ptr(int32(v1.NATLeaseM*60), zna.GetLeaseIntervalS()),
				RenewalIntervalS:    ptr(int32(v1.NATRenewalM*60), zna.GetRenewalIntervalS()),
				TimeoutIntervalS:    ptr(int32(v1.NATTimeoutS), zna.GetTimeoutIntervalS()),
				StunKeepaliveStartS: ptr(int32(v1.StunKeepaliveStartS), zna.GetStunKeepaliveStartS()),
				StunKeepaliveMinS:   ptr(int32(v1.StunKeepaliveMinS), zna.GetStunKeepaliveMinS()),
				// servers
			}.Build(),

			Discovery: configv2.OptionsConfiguration_Network_Discovery_builder{
				AnnounceLanAddresses: ptr(v1.AnnounceLANAddresses, false),
				Local: configv2.OptionsConfiguration_Network_Discovery_Local_builder{
					Enabled:              ptr(v1.LocalAnnEnabled, zl.GetEnabled()),
					Ipv4Port:             ptr(int32(v1.LocalAnnPort), zl.GetIpv4Port(), 0),
					Ipv6MulticastAddress: ptr(v1.LocalAnnMCAddr, zl.GetIpv6MulticastAddress(), ""),
				}.Build(),
				Global: configv2.OptionsConfiguration_Network_Discovery_Global_builder{
					Enabled: ptr(v1.GlobalAnnEnabled, zg.GetEnabled()),
					// servers
				}.Build(),
			}.Build(),

			Relays: configv2.OptionsConfiguration_Network_Relays_builder{
				Enabled:            ptr(v1.RelaysEnabled, zr.GetEnabled()),
				ReconnectIntervalS: ptr(int32(v1.RelayReconnectIntervalM*60), zr.GetReconnectIntervalS()),
			}.Build(),
		}.Build(),

		AutoUpgrade: configv2.OptionsConfiguration_AutoUpgrade_builder{
			Enabled:           ptr(v1.AutoUpgradeEnabled(), zau.GetEnabled()),
			CheckIntervalS:    ptr(int32(v1.AutoUpgradeIntervalH*3600), zau.GetCheckIntervalS()),
			ReleaseCandidates: ptr(v1.UpgradeToPreReleases, zau.GetReleaseCandidates()),
			ServerUrl:         ptr(v1.ReleasesURL, zau.GetServerUrl()),
		}.Build(),

		UsageReporting: configv2.OptionsConfiguration_UsageReporting_builder{
			Enabled:       ptr(v1.URAccepted > 0, zu.GetEnabled()),
			UniqueId:      ptr(v1.URUniqueID, ""),
			Url:           ptr(v1.URURL, zu.GetUrl()),
			InitialDelayS: ptr(int32(v1.URInitialDelayS), zu.GetInitialDelayS()),
			TlsInsecure:   ptr(v1.URPostInsecurely, zu.GetTlsInsecure()),
		}.Build(),

		Audit: configv2.OptionsConfiguration_Audit_builder{
			Enabled: ptr(v1.AuditEnabled, false),
			File:    ptr(v1.AuditFile, ""),
		}.Build(),
	}
	return bld.Build()
}

func ptr[T comparable](v T, defs ...T) *T {
	for _, def := range defs {
		if v == def {
			return nil
		}
	}
	return &v
}
