// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"slices"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/protocol"
	"google.golang.org/protobuf/proto"
)

// The configuration tree is exposed through wrapper types around the
// generated messages, providing device IDs as the native type instead of
// strings. Types that do not contain device IDs are exposed as aliases of
// the generated types.

// Configuration is the Syncthing configuration.
type Configuration struct {
	*syncthingv2.Configuration
}

func (c Configuration) Copy() Configuration {
	cp := proto.Clone(c.Configuration).(*syncthingv2.Configuration)
	return Configuration{cp}
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
	pbs := make([]*syncthingv2.FolderConfiguration, len(folders))
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
	pbs := make([]*syncthingv2.DeviceConfiguration, len(devices))
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
	pbs := make([]*syncthingv2.ObservedDevice, len(devices))
	for i, device := range devices {
		pbs[i] = device.ObservedDevice
	}
	c.Configuration.SetRemoteIgnoredDevices(pbs)
}

// GetOptions returns the global configuration options.
func (c Configuration) GetOptions() OptionsConfiguration {
	return OptionsConfiguration{c.Configuration.GetOptions()}
}

// SetOptions sets the global configuration options.
func (c Configuration) SetOptions(options OptionsConfiguration) {
	c.Configuration.SetOptions(options.OptionsConfiguration)
}

// Device returns the configuration for the given device, if present, and
// its index in the device list.
func (c Configuration) Device(id protocol.DeviceID) (DeviceConfiguration, int, bool) {
	for i, device := range c.GetDevices() {
		if device.GetDeviceId() == id {
			return device, i, true
		}
	}
	return DeviceConfiguration{}, 0, false
}

// Folder returns the configuration for the given folder, if present, and
// its index in the folder list.
func (c Configuration) Folder(id string) (FolderConfiguration, int, bool) {
	for i, folder := range c.GetFolders() {
		if folder.GetId() == id {
			return folder, i, true
		}
	}
	return FolderConfiguration{}, 0, false
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
	*syncthingv2.Defaults
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
	*syncthingv2.FolderConfiguration
}

func (f FolderConfiguration) Copy() FolderConfiguration {
	cp := proto.Clone(f.FolderConfiguration).(*syncthingv2.FolderConfiguration)
	return FolderConfiguration{cp}
}

// Device returns the configuration for the given device among the
// devices the folder is shared with, if present.
func (f FolderConfiguration) Device(id protocol.DeviceID) (FolderDeviceConfiguration, bool) {
	for _, device := range f.GetDevices() {
		if device.GetDeviceId() == id {
			return device, true
		}
	}
	return FolderDeviceConfiguration{}, false
}

// MinDiskFree returns the minimum required free space on the disk the
// folder resides on. If unset, the default is one percent.
func (f FolderConfiguration) MinDiskFree() *Size {
	if size := f.GetMinDiskFree(); size != nil {
		return size
	}
	return syncthingv2.Size_builder{Percent: new(1.0)}.Build()
}

// materializeDefaults sets the fields that have getter-declared defaults
// to their effective values, when unset. It is used when materialising
// the effective configuration.
func (f FolderConfiguration) materializeDefaults() {
	if f.FolderConfiguration == nil {
		return
	}
	if size := f.MinDiskFree(); !proto.Equal(size, f.GetMinDiskFree()) {
		f.SetMinDiskFree(size)
	}
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
	pbs := make([]*syncthingv2.FolderDeviceConfiguration, len(devices))
	for i, device := range devices {
		pbs[i] = device.FolderDeviceConfiguration
	}
	f.FolderConfiguration.SetDevices(pbs)
}

// FolderDeviceConfiguration is a device that a folder is shared with.
type FolderDeviceConfiguration struct {
	*syncthingv2.FolderDeviceConfiguration
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
	*syncthingv2.DeviceConfiguration
}

func (d DeviceConfiguration) Copy() DeviceConfiguration {
	cp := proto.Clone(d.DeviceConfiguration).(*syncthingv2.DeviceConfiguration)
	return DeviceConfiguration{cp}
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

// Addresses returns the addresses or host names to use when attempting to
// connect to this device. If unset, the default is to use local and
// global discovery to find the device.
func (d DeviceConfiguration) Addresses() []string {
	if addrs := d.GetAddresses(); len(addrs) > 0 && !(len(addrs) == 1 && addrs[0] == "") {
		return addrs
	}
	return []string{"dynamic"}
}

// materializeDefaults sets the fields that have getter-declared defaults
// to their effective values, when unset. It is used when materialising
// the effective configuration.
func (d DeviceConfiguration) materializeDefaults() {
	if d.DeviceConfiguration == nil {
		return
	}
	if addrs := d.Addresses(); !slices.Equal(addrs, d.GetAddresses()) {
		d.SetAddresses(addrs)
	}
}

// SetIntroducedBy sets the device ID of the introducer. The empty device
// ID clears the value.
func (d DeviceConfiguration) SetIntroducedBy(id protocol.DeviceID) {
	d.DeviceConfiguration.SetIntroducedBy(id.String())
}

// ObservedDevice is a device encountered on the network.
type ObservedDevice struct {
	*syncthingv2.ObservedDevice
}

func (d ObservedDevice) Copy() ObservedDevice {
	cp := proto.Clone(d.ObservedDevice).(*syncthingv2.ObservedDevice)
	return ObservedDevice{cp}
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

// Types without device IDs or defaults are used as-is.

type (
	GUIConfiguration        = syncthingv2.GUIConfiguration
	LDAPConfiguration       = syncthingv2.LDAPConfiguration
	VersioningConfiguration = syncthingv2.VersioningConfiguration
	XattrFilter             = syncthingv2.XattrFilter
	XattrFilterEntry        = syncthingv2.XattrFilterEntry
	Size                    = syncthingv2.Size
	Ignores                 = syncthingv2.Ignores
	ObservedFolder          = syncthingv2.ObservedFolder
)

// OptionsConfiguration is the global configuration options. It wraps the
// generated type to provide defaults for fields that cannot have them in
// the schema, such as repeated fields.
type OptionsConfiguration struct {
	*syncthingv2.OptionsConfiguration
}

func (o OptionsConfiguration) Copy() OptionsConfiguration {
	cp := proto.Clone(o.OptionsConfiguration).(*syncthingv2.OptionsConfiguration)
	return OptionsConfiguration{cp}
}

// ListenAddresses returns the listen addresses for incoming sync
// connections. If unset, the default is the "default" address set, which
// the consumer of the configuration expands to the actual addresses.
func (o OptionsConfiguration) ListenAddresses() []string {
	if addrs := o.GetListenAddresses(); len(addrs) > 0 {
		return addrs
	}
	return []string{"default"}
}

// GlobalAnnounceServers returns the global announce (discovery) servers.
// If unset, the default is the "default" server set, which the consumer
// of the configuration expands to the actual servers.
func (o OptionsConfiguration) GlobalAnnounceServers() []string {
	if servers := o.GetGlobalAnnounceServers(); len(servers) > 0 {
		return servers
	}
	return []string{"default"}
}

// StunServers returns the servers to use for STUN. If unset, the default
// is the "default" server set, which the consumer of the configuration
// expands to the actual servers.
func (o OptionsConfiguration) StunServers() []string {
	if servers := o.GetStunServers(); len(servers) > 0 {
		return servers
	}
	return []string{"default"}
}

// MinHomeDiskFree returns the minimum required free space on the partition
// holding the configuration and index. If unset, the default is one
// percent.
func (o OptionsConfiguration) MinHomeDiskFree() *Size {
	if size := o.GetMinHomeDiskFree(); size != nil {
		return size
	}
	return syncthingv2.Size_builder{Percent: new(1.0)}.Build()
}

// materializeDefaults sets the fields that have getter-declared defaults
// to their effective values, when unset. It is used when materialising
// the effective configuration.
func (o OptionsConfiguration) materializeDefaults() {
	if o.OptionsConfiguration == nil {
		return
	}
	if addrs := o.ListenAddresses(); !slices.Equal(addrs, o.GetListenAddresses()) {
		o.SetListenAddresses(addrs)
	}
	if servers := o.GlobalAnnounceServers(); !slices.Equal(servers, o.GetGlobalAnnounceServers()) {
		o.SetGlobalAnnounceServers(servers)
	}
	if servers := o.StunServers(); !slices.Equal(servers, o.GetStunServers()) {
		o.SetStunServers(servers)
	}
	if size := o.MinHomeDiskFree(); !proto.Equal(size, o.GetMinHomeDiskFree()) {
		o.SetMinHomeDiskFree(size)
	}
}
