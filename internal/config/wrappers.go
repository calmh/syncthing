// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	configpb "github.com/syncthing/syncthing/internal/gen/syncthing/v2/config"
	"github.com/syncthing/syncthing/lib/protocol"
)

// The configuration tree is exposed through wrapper types around the
// generated messages, providing device IDs as the native type instead of
// strings. Types that do not contain device IDs are exposed as aliases of
// the generated types.

// Configuration is the Syncthing configuration.
type Configuration struct {
	*configpb.Configuration
}

// GetFolders returns the folder configurations.
func (c Configuration) GetFolders() []FolderConfiguration {
	pbs := c.Configuration.GetFolders()
	out := make([]FolderConfiguration, len(pbs))
	for i, pb := range pbs {
		out[i] = FolderConfiguration{pb}
	}
	return out
}

// SetFolders sets the folder configurations.
func (c Configuration) SetFolders(folders []FolderConfiguration) {
	pbs := make([]*configpb.FolderConfiguration, len(folders))
	for i, folder := range folders {
		pbs[i] = folder.FolderConfiguration
	}
	c.Configuration.SetFolders(pbs)
}

// GetDevices returns the device configurations.
func (c Configuration) GetDevices() []DeviceConfiguration {
	pbs := c.Configuration.GetDevices()
	out := make([]DeviceConfiguration, len(pbs))
	for i, pb := range pbs {
		out[i] = DeviceConfiguration{pb}
	}
	return out
}

// SetDevices sets the device configurations.
func (c Configuration) SetDevices(devices []DeviceConfiguration) {
	pbs := make([]*configpb.DeviceConfiguration, len(devices))
	for i, device := range devices {
		pbs[i] = device.DeviceConfiguration
	}
	c.Configuration.SetDevices(pbs)
}

// GetRemoteIgnoredDevices returns the devices that have been ignored after
// being seen remotely.
func (c Configuration) GetRemoteIgnoredDevices() []ObservedDevice {
	pbs := c.Configuration.GetRemoteIgnoredDevices()
	out := make([]ObservedDevice, len(pbs))
	for i, pb := range pbs {
		out[i] = ObservedDevice{pb}
	}
	return out
}

// SetRemoteIgnoredDevices sets the devices that have been ignored after
// being seen remotely.
func (c Configuration) SetRemoteIgnoredDevices(devices []ObservedDevice) {
	pbs := make([]*configpb.ObservedDevice, len(devices))
	for i, device := range devices {
		pbs[i] = device.ObservedDevice
	}
	c.Configuration.SetRemoteIgnoredDevices(pbs)
}

// GetDefaults returns the default values for new folders and devices.
func (c Configuration) GetDefaults() Defaults {
	return Defaults{c.Configuration.GetDefaults()}
}

// SetDefaults sets the default values for new folders and devices.
func (c Configuration) SetDefaults(defaults Defaults) {
	c.Configuration.SetDefaults(defaults.Defaults)
}

// Defaults are the default values for new folders and devices.
type Defaults struct {
	*configpb.Defaults
}

// GetFolder returns the default folder configuration.
func (d Defaults) GetFolder() FolderConfiguration {
	return FolderConfiguration{d.Defaults.GetFolder()}
}

// SetFolder sets the default folder configuration.
func (d Defaults) SetFolder(folder FolderConfiguration) {
	d.Defaults.SetFolder(folder.FolderConfiguration)
}

// GetDevice returns the default device configuration.
func (d Defaults) GetDevice() DeviceConfiguration {
	return DeviceConfiguration{d.Defaults.GetDevice()}
}

// SetDevice sets the default device configuration.
func (d Defaults) SetDevice(device DeviceConfiguration) {
	d.Defaults.SetDevice(device.DeviceConfiguration)
}

// FolderConfiguration is the configuration for a shared folder.
type FolderConfiguration struct {
	*configpb.FolderConfiguration
}

// GetDevices returns the devices the folder is shared with.
func (f FolderConfiguration) GetDevices() []FolderDeviceConfiguration {
	pbs := f.FolderConfiguration.GetDevices()
	out := make([]FolderDeviceConfiguration, len(pbs))
	for i, pb := range pbs {
		out[i] = FolderDeviceConfiguration{pb}
	}
	return out
}

// SetDevices sets the devices the folder is shared with.
func (f FolderConfiguration) SetDevices(devices []FolderDeviceConfiguration) {
	pbs := make([]*configpb.FolderDeviceConfiguration, len(devices))
	for i, device := range devices {
		pbs[i] = device.FolderDeviceConfiguration
	}
	f.FolderConfiguration.SetDevices(pbs)
}

// FolderDeviceConfiguration is a device that a folder is shared with.
type FolderDeviceConfiguration struct {
	*configpb.FolderDeviceConfiguration
}

// GetDeviceId returns the device ID. An unset or empty device ID is
// returned as the empty device ID.
func (d FolderDeviceConfiguration) GetDeviceId() protocol.DeviceID {
	return deviceIDFromString(d.FolderDeviceConfiguration.GetDeviceId())
}

// SetDeviceId sets the device ID. The empty device ID clears the value.
func (d FolderDeviceConfiguration) SetDeviceId(id protocol.DeviceID) {
	d.FolderDeviceConfiguration.SetDeviceId(id.String())
}

// GetIntroducedBy returns the device ID of the introducer. An unset or
// empty device ID is returned as the empty device ID.
func (d FolderDeviceConfiguration) GetIntroducedBy() protocol.DeviceID {
	return deviceIDFromString(d.FolderDeviceConfiguration.GetIntroducedBy())
}

// SetIntroducedBy sets the device ID of the introducer. The empty device
// ID clears the value.
func (d FolderDeviceConfiguration) SetIntroducedBy(id protocol.DeviceID) {
	d.FolderDeviceConfiguration.SetIntroducedBy(id.String())
}

// DeviceConfiguration is the configuration for a device.
type DeviceConfiguration struct {
	*configpb.DeviceConfiguration
}

// GetDeviceId returns the device ID. An unset or empty device ID is
// returned as the empty device ID.
func (d DeviceConfiguration) GetDeviceId() protocol.DeviceID {
	return deviceIDFromString(d.DeviceConfiguration.GetDeviceId())
}

// SetDeviceId sets the device ID. The empty device ID clears the value.
func (d DeviceConfiguration) SetDeviceId(id protocol.DeviceID) {
	d.DeviceConfiguration.SetDeviceId(id.String())
}

// GetIntroducedBy returns the device ID of the introducer. An unset or
// empty device ID is returned as the empty device ID.
func (d DeviceConfiguration) GetIntroducedBy() protocol.DeviceID {
	return deviceIDFromString(d.DeviceConfiguration.GetIntroducedBy())
}

// SetIntroducedBy sets the device ID of the introducer. The empty device
// ID clears the value.
func (d DeviceConfiguration) SetIntroducedBy(id protocol.DeviceID) {
	d.DeviceConfiguration.SetIntroducedBy(id.String())
}

// ObservedDevice is a device encountered on the network.
type ObservedDevice struct {
	*configpb.ObservedDevice
}

// GetDeviceId returns the device ID. An unset or empty device ID is
// returned as the empty device ID.
func (d ObservedDevice) GetDeviceId() protocol.DeviceID {
	return deviceIDFromString(d.ObservedDevice.GetDeviceId())
}

// SetDeviceId sets the device ID. The empty device ID clears the value.
func (d ObservedDevice) SetDeviceId(id protocol.DeviceID) {
	d.ObservedDevice.SetDeviceId(id.String())
}

// deviceIDFromString parses a device ID string. An invalid value parses
// as the empty device ID; configurations are expected to have been
// validated, so this should not occur in practice.
func deviceIDFromString(s string) protocol.DeviceID {
	id, err := protocol.DeviceIDFromString(s)
	if err != nil {
		return protocol.EmptyDeviceID
	}
	return id
}

// Types without device IDs are used as-is.

type (
	GUIConfiguration        = configpb.GUIConfiguration
	LDAPConfiguration       = configpb.LDAPConfiguration
	OptionsConfiguration    = configpb.OptionsConfiguration
	VersioningConfiguration = configpb.VersioningConfiguration
	XattrFilter             = configpb.XattrFilter
	XattrFilterEntry        = configpb.XattrFilterEntry
	Size                    = configpb.Size
	Ignores                 = configpb.Ignores
	ObservedFolder          = configpb.ObservedFolder
)
