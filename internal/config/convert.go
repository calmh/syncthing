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
		devices[i] = configv2.DeviceConfiguration_builder{
			Id: ptr(dev.DeviceID.String()),
		}.Build()
	}

	return configv2.Configuration_builder{
		Folders: folders,
		Devices: devices,
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
		Enabled:    nondefPtr(!fld.Paused, zf.GetEnabled()),
		Label:      nondefPtr(fld.Label, ""),
		MarkerName: nondefPtr(fld.MarkerName, zf.GetMarkerName()),

		Filesystem: configv2.FolderConfiguration_Filesystem_builder{
			Path: ptr(fld.Path),
			Type: v2FSType(fld.FilesystemType),
		}.Build(),

		FilesystemOptions: configv2.FolderConfiguration_FilesystemOptions_builder{
			CaseSensitive:   nondefPtr(fld.CaseSensitiveFS, zfso.GetCaseSensitive()),
			JunctionsAsDirs: nondefPtr(fld.JunctionsAsDirs, zfso.GetJunctionsAsDirs()),
			MinDiskFree:     v2Size(fld.MinDiskFree),
			AutoNormalize:   nondefPtr(fld.AutoNormalize, zfso.GetAutoNormalize()),
		}.Build(),

		Scanning: configv2.FolderConfiguration_Scanning_builder{
			Watcher: configv2.FolderConfiguration_Scanning_Watcher_builder{
				Enabled:  nondefPtr(fld.FSWatcherEnabled, zwa.GetEnabled()),
				DelayS:   nondefPtr(fld.FSWatcherDelayS, zwa.GetDelayS()),
				TimeoutS: nondefPtr(fld.FSWatcherTimeoutS, zwa.GetTimeoutS()),
			}.Build(),
			RescanIntervalS:   nondefPtr(int32(fld.RescanIntervalS), zsc.GetRescanIntervalS()),
			NumHashers:        nondefPtr(int32(fld.Hashers), zsc.GetNumHashers()),
			ModTimeWindowS:    nondefPtr(int32(fld.RawModTimeWindowS), zsc.GetModTimeWindowS()),
			ProgressIntervalS: nondefPtr(int32(fld.ScanProgressIntervalS), zsc.GetProgressIntervalS()),
			ScanOwnership:     nondefPtr(fld.SendOwnership, zsc.GetScanOwnership()),
			ScanXattrs:        nondefPtr(fld.SendXattrs, zsc.GetScanXattrs()),
			XattrFilter:       configv2.XattrFilter_builder{}.Build(),
		}.Build(),

		Pulling: configv2.FolderConfiguration_Pulling_builder{
			Order:                   v2PullOrder(fld.Order),
			BlockOrder:              v2BlockPullOrder(fld.BlockPullOrder),
			NumCopiers:              nondefPtr(int32(fld.Copiers), zpu.GetNumCopiers(), 0),
			MaxPendingKib:           nondefPtr(int32(fld.PullerMaxPendingKiB), zpu.GetMaxPendingKib(), 0),
			MaxConflicts:            nondefPtr(int32(fld.MaxConflicts), zpu.GetMaxConflicts()),
			MaxConcurrentWrites:     nondefPtr(int32(fld.MaxConcurrentWrites), zpu.GetMaxConcurrentWrites(), 0),
			FailedPauseS:            nondefPtr(int32(fld.PullerPauseS), zpu.GetFailedPauseS(), 0),
			ChangeDelayS:            nondefPtr(fld.PullerDelayS, zpu.GetChangeDelayS()),
			SparseFiles:             nondefPtr(!fld.DisableSparseFiles, zpu.GetSparseFiles()),
			IgnorePermissions:       nondefPtr(fld.IgnorePerms, zpu.GetIgnorePermissions()),
			TemporaryIndexes:        nondefPtr(!fld.DisableTempIndexes, zpu.GetTemporaryIndexes()),
			Fsync:                   nondefPtr(!fld.DisableFsync, zpu.GetFsync()),
			CopyRangeMethod:         v2CopyRangeMethod(fld.CopyRangeMethod),
			SyncOwnership:           nondefPtr(fld.SyncOwnership, zpu.GetSyncOwnership()),
			SyncXattrs:              nondefPtr(fld.SyncXattrs, zpu.GetSyncXattrs()),
			IgnoreDelete:            nondefPtr(fld.IgnoreDelete, zpu.GetIgnoreDelete()),
			CopyOwnershipFromParent: nondefPtr(fld.CopyOwnershipFromParent, zpu.GetCopyOwnershipFromParent()),
		}.Build(),
	}

	for _, dev := range fld.Devices {
		bld.SharedWith = append(bld.SharedWith, configv2.FolderConfiguration_Sharing_builder{
			DeviceId:           ptr(dev.DeviceID.String()),
			IntroducedBy:       nondefPtr(dev.IntroducedBy.String(), ""),
			EncryptionPassword: nondefPtr(dev.EncryptionPassword, ""),
		}.Build())
	}

	if vt := v2VersioningType(fld.Versioning.Type); vt != nil {
		zv := (*configv2.FolderConfiguration_Versioning)(nil)
		vb := configv2.FolderConfiguration_Versioning_builder{
			Type: vt,
			Filesystem: configv2.FolderConfiguration_Filesystem_builder{
				Path: nondefPtr(fld.Versioning.FSPath, ""),
				Type: v2FSType(fld.Versioning.FSType),
			}.Build(),
			CleanupIntervalS: nondefPtr(int32(fld.Versioning.CleanupIntervalS), zv.GetCleanupIntervalS()),
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

func ptr[T any](v T) *T {
	return &v
}

func nondefPtr[T comparable](v T, defs ...T) *T {
	for _, def := range defs {
		if v == def {
			return nil
		}
	}
	return &v
}
