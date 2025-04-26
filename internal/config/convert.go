package config

import (
	configv1 "github.com/syncthing/syncthing/internal/config/v1"
	configv2 "github.com/syncthing/syncthing/internal/config/v2"
)

func FromV1(v1 *configv1.Configuration) *configv2.Configuration {
	folders := make([]*configv2.FolderConfiguration, len(v1.Folders))
	for i, fld := range v1.Folders {
		bld := configv2.FolderConfiguration_builder{
			Type:   v1FolderType(fld.Type),
			Id:     ptr(fld.ID),
			Label:  nondefPtr(fld.Label, ""),
			Paused: nondefPtr(fld.Paused, false),
			Filesystem: configv2.FolderConfiguration_Filesystem_builder{
				Path:            ptr(fld.Path),
				FsType:          v1FSType(fld.FilesystemType),
				CaseSensitiveFs: nondefPtr(fld.CaseSensitiveFS, false),
				JunctionsAsDirs: nondefPtr(fld.JunctionsAsDirs, false),
				MinDiskFree:     v1Size(fld.MinDiskFree),
			}.Build(),
			Scanning: configv2.FolderConfiguration_Scanning_builder{
				Watcher: configv2.FolderConfiguration_Scanning_Watcher_builder{
					Enabled:  nondefPtr(fld.FSWatcherEnabled, true),
					DelayS:   nondefPtr(fld.FSWatcherDelayS, 0),
					TimeoutS: nondefPtr(fld.FSWatcherTimeoutS, 0),
				}.Build(),
				AutoNormalize:     nondefPtr(fld.AutoNormalize, true),
				Ownership:         nondefPtr(fld.SendOwnership, false),
				Xattrs:            nondefPtr(fld.SendXattrs, false),
				Hashers:           nondefPtr(int32(fld.Hashers), 0),
				ModTimeWindowS:    nondefPtr(int32(fld.RawModTimeWindowS), 0),
				ProgressIntervalS: nondefPtr(int32(fld.ScanProgressIntervalS), 0),
				RescanIntervalS:   nondefPtr(int32(fld.RescanIntervalS), 0),
				MarkerName:        nondefPtr(fld.MarkerName, ".stfolder"),
				XattrFilter:       configv2.XattrFilter_builder{}.Build(),
			}.Build(),
		}
		folders[i] = bld.Build()
	}

	devices := make([]*configv2.DeviceConfiguration, len(v1.Devices))
	for i, dev := range v1.Devices {
		devices[i] = configv2.DeviceConfiguration_builder{
			Id: ptr(dev.DeviceID.String()),
		}.Build()
	}

	return configv2.Configuration_builder{
		Devices: devices,
	}.Build()
}

func v1FolderType(v1 configv1.FolderType) *configv2.FolderType {
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

func v1FSType(v1 configv1.FilesystemType) *configv2.FilesystemType {
	switch v1 {
	case configv1.FilesystemTypeBasic:
		return configv2.FilesystemType_FILESYSTEM_TYPE_BASIC.Enum()
	case configv1.FilesystemTypeFake:
		return configv2.FilesystemType_FILESYSTEM_TYPE_FAKE.Enum()
	default:
		return nil
	}
}

func v1Size(v1 configv1.Size) *configv2.Size {
	bld := configv2.Size_builder{
		Value: ptr(v1.BaseValue()),
	}
	if v1.Unit == "%" {
		bld.Unit = configv2.SizeUnit_SIZE_UNIT_PERCENT.Enum()
	} else {
		bld.Unit = configv2.SizeUnit_SIZE_UNIT_BYTES.Enum()
	}
	return bld.Build()
}

func ptr[T any](v T) *T {
	return &v
}

func nondefPtr[T comparable](v, def T) *T {
	if v == def {
		return nil
	}
	return &v
}
