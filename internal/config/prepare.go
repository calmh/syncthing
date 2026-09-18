// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/protocol"
	"github.com/syncthing/syncthing/lib/rand"
	"github.com/syncthing/syncthing/lib/sliceutil"
	"github.com/syncthing/syncthing/lib/stringutil"
)

var (
	errFolderIDEmpty     = errors.New("folder has empty ID")
	errFolderIDDuplicate = errors.New("folder has duplicate ID")
	errFolderPathEmpty   = errors.New("folder has empty path")
)

const maxConcurrentWritesLimit = 256

// Prepare normalises the configuration in preparation for use: the local
// device is added to the device list and all folders, device and folder
// lists are deduplicated and sorted, references between folders and
// devices are made consistent, and values are clamped to sane ranges.
// Defaults are not materialised; they are declared in the schema and
// exposed through getters.
//
// An error is returned if the configuration is inconsistent in a way
// that cannot be resolved, e.g. duplicate folder IDs.
func Prepare(cfg *Configuration, myID protocol.DeviceID) error {
	if cfg.Configuration == nil {
		return nil
	}

	ensureMyDevice(cfg, myID)

	existingDevices := prepareDeviceList(cfg)
	sharedFolders, err := prepareFolders(cfg, myID, existingDevices)
	if err != nil {
		return err
	}
	prepareDevices(cfg, sharedFolders)

	prepareGUI(cfg)

	guiPWIsSet := cfg.GetGui().GetUser() != "" && cfg.GetGui().GetPassword() != ""
	prepareOptions(cfg.GetOptions(), guiPWIsSet)

	prepareIgnoredDevices(cfg, existingDevices)
	prepareDefaults(cfg, myID, existingDevices)

	return nil
}

// ensureMyDevice adds the local device to the device list, if it is not
// already present.
func ensureMyDevice(cfg *Configuration, myID protocol.DeviceID) {
	if myID == protocol.EmptyDeviceID {
		return
	}
	for _, device := range cfg.GetDevices() {
		if device.GetDeviceId() == myID {
			return
		}
	}
	myName, _ := os.Hostname()
	cfg.SetDevices(append(cfg.GetDevices(), DeviceConfiguration{
		DeviceConfiguration: syncthingv2.DeviceConfiguration_builder{
			DeviceId: new(myID.String()),
			Name:     new(myName),
		}.Build(),
	}))
}

// prepareDeviceList removes devices with empty or duplicate IDs from the
// device list, sorts it by device ID, and returns the remaining devices.
func prepareDeviceList(cfg *Configuration) map[protocol.DeviceID]DeviceConfiguration {
	devices := cfg.GetDevices()
	seen := make(map[protocol.DeviceID]bool, len(devices))
	filtered := make([]DeviceConfiguration, 0, len(devices))
	for _, device := range devices {
		id := device.GetDeviceId()
		if id == protocol.EmptyDeviceID || seen[id] {
			continue
		}
		seen[id] = true
		filtered = append(filtered, device)
	}
	slices.SortFunc(filtered, func(a, b DeviceConfiguration) int {
		return a.GetDeviceId().Compare(b.GetDeviceId())
	})
	cfg.SetDevices(filtered)

	existing := make(map[protocol.DeviceID]DeviceConfiguration, len(filtered))
	for _, device := range filtered {
		existing[device.GetDeviceId()] = device
	}
	return existing
}

// prepareFolders prepares each folder and sorts the folder list by ID. It
// returns the folders shared with each device. An error is returned for
// folders with empty or duplicate IDs or an empty path.
func prepareFolders(cfg *Configuration, myID protocol.DeviceID, existingDevices map[protocol.DeviceID]DeviceConfiguration) (map[protocol.DeviceID][]string, error) {
	folders := cfg.GetFolders()
	sharedFolders := make(map[protocol.DeviceID][]string, len(folders))
	existing := make(map[string]bool, len(folders))
	for _, folder := range folders {
		if folder.GetId() == "" {
			return nil, errFolderIDEmpty
		}
		if folder.GetPath() == "" {
			return nil, fmt.Errorf("folder %q: %w", folder.GetId(), errFolderPathEmpty)
		}
		if existing[folder.GetId()] {
			return nil, fmt.Errorf("folder %q: %w", folder.GetId(), errFolderIDDuplicate)
		}
		existing[folder.GetId()] = true

		prepareFolder(folder, myID, existingDevices)

		for _, dev := range folder.GetDevices() {
			sharedFolders[dev.GetDeviceId()] = append(sharedFolders[dev.GetDeviceId()], folder.GetId())
		}
	}
	slices.SortFunc(folders, func(a, b FolderConfiguration) int {
		return strings.Compare(a.GetId(), b.GetId())
	})
	cfg.SetFolders(folders)
	return sharedFolders, nil
}

// prepareFolder normalises the folder configuration. It ensures that the
// folder is only shared with known and non-duplicate devices, that the
// local device is present, and that values are within sane ranges.
func prepareFolder(f FolderConfiguration, myID protocol.DeviceID, existingDevices map[protocol.DeviceID]DeviceConfiguration) {
	// Ensure that
	// - any loose devices are not present in the wrong places
	// - there are no duplicate devices
	// - we are part of the devices
	// - folder is not shared in trusted mode with an untrusted device
	f.SetDevices(ensureExistingDevices(f.GetDevices(), existingDevices))
	f.SetDevices(ensureNoDuplicateFolderDevices(f.GetDevices()))
	f.SetDevices(ensureDevicePresent(f.GetDevices(), myID))
	f.SetDevices(ensureNoUntrustedTrustingSharing(f, f.GetDevices(), existingDevices))

	slices.SortFunc(f.GetDevices(), func(a, b FolderDeviceConfiguration) int {
		return a.GetDeviceId().Compare(b.GetDeviceId())
	})

	if f.HasRescanIntervalS() {
		if f.GetRescanIntervalS() > MaxRescanIntervalS {
			f.SetRescanIntervalS(MaxRescanIntervalS)
		} else if f.GetRescanIntervalS() < 0 {
			// A negative interval is meaningless; zero disables the
			// periodic rescan.
			f.SetRescanIntervalS(0)
		}
	}

	if f.HasFsWatcherDelayS() {
		if f.GetFsWatcherDelayS() <= 0 {
			// A zero or negative delay disables the watcher.
			f.SetFsWatcherEnabled(false)
		} else if f.GetFsWatcherDelayS() < 0.01 {
			f.SetFsWatcherDelayS(0.01)
		}
	}

	if v := f.GetVersioning(); v != nil && v.HasCleanupIntervalS() {
		if v.GetCleanupIntervalS() > MaxRescanIntervalS {
			v.SetCleanupIntervalS(MaxRescanIntervalS)
		} else if v.GetCleanupIntervalS() < 0 {
			v.SetCleanupIntervalS(0)
		}
	}

	if f.HasMaxConcurrentWrites() {
		if f.GetMaxConcurrentWrites() <= 0 {
			// Zero or negative means the default.
			f.ClearMaxConcurrentWrites()
		} else if f.GetMaxConcurrentWrites() > maxConcurrentWritesLimit {
			f.SetMaxConcurrentWrites(maxConcurrentWritesLimit)
		}
	}

	if f.GetType() == syncthingv2.FolderType_FOLDER_TYPE_RECEIVE_ENCRYPTED {
		f.SetIgnorePerms(true)
	}
}

func ensureExistingDevices(devices []FolderDeviceConfiguration, existingDevices map[protocol.DeviceID]DeviceConfiguration) []FolderDeviceConfiguration {
	count := len(devices)
	i := 0
loop:
	for i < count {
		if _, ok := existingDevices[devices[i].GetDeviceId()]; !ok {
			devices[i] = devices[count-1]
			count--
			continue loop
		}
		i++
	}
	return devices[0:count]
}

func ensureNoDuplicateFolderDevices(devices []FolderDeviceConfiguration) []FolderDeviceConfiguration {
	count := len(devices)
	i := 0
	seenDevices := make(map[protocol.DeviceID]bool)
loop:
	for i < count {
		id := devices[i].GetDeviceId()
		if _, ok := seenDevices[id]; ok {
			devices[i] = devices[count-1]
			count--
			continue loop
		}
		seenDevices[id] = true
		i++
	}
	return devices[0:count]
}

func ensureDevicePresent(devices []FolderDeviceConfiguration, myID protocol.DeviceID) []FolderDeviceConfiguration {
	if myID == protocol.EmptyDeviceID {
		return devices
	}
	for _, device := range devices {
		if device.GetDeviceId() == myID {
			return devices
		}
	}
	return append(devices, FolderDeviceConfiguration{
		FolderDeviceConfiguration: syncthingv2.FolderDeviceConfiguration_builder{
			DeviceId: new(myID.String()),
		}.Build(),
	})
}

func ensureNoUntrustedTrustingSharing(f FolderConfiguration, devices []FolderDeviceConfiguration, existingDevices map[protocol.DeviceID]DeviceConfiguration) []FolderDeviceConfiguration {
	for i := 0; i < len(devices); i++ {
		dev := devices[i]
		if dev.GetEncryptionPassword() != "" || f.GetType() == syncthingv2.FolderType_FOLDER_TYPE_RECEIVE_ENCRYPTED {
			// There's a password set or the folder is received
			// encrypted, no check required
			continue
		}
		if devCfg, ok := existingDevices[dev.GetDeviceId()]; ok && devCfg.GetUntrusted() {
			slog.Error("Folder is shared in trusted mode with untrusted device; unsharing",
				"device", dev.GetDeviceId(), "folder", f.GetId())
			devices = sliceutil.RemoveAndZero(devices, i)
			i--
		}
	}
	return devices
}

// prepareDevices prepares each device configuration.
func prepareDevices(cfg *Configuration, sharedFolders map[protocol.DeviceID][]string) {
	for _, device := range cfg.GetDevices() {
		prepareDevice(device, sharedFolders[device.GetDeviceId()])
	}
}

// prepareDevice normalises the device configuration: the list of ignored
// folders is deduplicated and folders that are shared are removed from
// it, and an untrusted device is neither an introducer nor
// auto-accepting folders.
func prepareDevice(device DeviceConfiguration, sharedFolders []string) {
	ignoredFolders := deduplicateObservedFolders(device.GetIgnoredFolders())
	kept := make([]*syncthingv2.ObservedFolder, 0, len(ignoredFolders))
	for _, folder := range ignoredFolders {
		if !slices.Contains(sharedFolders, folder.GetId()) {
			kept = append(kept, folder)
		}
	}
	device.SetIgnoredFolders(kept)

	// A device cannot be simultaneously untrusted and an introducer, nor
	// auto accept folders.
	if device.GetUntrusted() {
		if device.GetIntroducer() {
			slog.Warn("Device is both untrusted and an introducer, removing introducer flag", "device", device.GetDeviceId())
			device.SetIntroducer(false)
		}
		if device.GetAutoAcceptFolders() {
			slog.Warn("Device is both untrusted and auto-accepting folders, removing auto-accept flag", "device", device.GetDeviceId())
			device.SetAutoAcceptFolders(false)
		}
	}
}

// prepareGUI ensures that the GUI configuration exists and has an API
// key, and normalises the session cookie path.
func prepareGUI(cfg *Configuration) {
	gui := cfg.GetGui()
	if gui.GUIConfiguration == nil {
		gui = GUIConfiguration{syncthingv2.GUIConfiguration_builder{}.Build()}
		cfg.SetGui(gui)
	}
	if gui.GetApiKey() == "" {
		gui.SetApiKey(rand.String(32))
	}
	if path := strings.TrimSpace(gui.GetSessionCookiePath()); path != "" {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		gui.SetSessionCookiePath(path)
	}
}

// prepareOptions normalises the global options: address lists are
// deduplicated, intervals are clamped to sane ranges, and a unique ID is
// generated for usage reporting when it is enabled.
func prepareOptions(opts OptionsConfiguration, guiPWIsSet bool) {
	if opts.OptionsConfiguration == nil {
		return
	}

	if addrs := opts.GetListenAddresses(); len(addrs) > 0 {
		opts.SetListenAddresses(stringutil.UniqueTrimmedStrings(addrs))
	}
	if servers := opts.GetGlobalAnnounceServers(); len(servers) > 0 {
		opts.SetGlobalAnnounceServers(stringutil.UniqueTrimmedStrings(servers))
	}
	if nets := opts.GetAlwaysLocalNets(); len(nets) > 0 {
		opts.SetAlwaysLocalNets(stringutil.UniqueTrimmedStrings(nets))
	}
	if servers := opts.GetStunServers(); len(servers) > 0 {
		opts.SetStunServers(stringutil.UniqueTrimmedStrings(servers))
	}
	if ids := opts.GetUnackedNotificationIds(); len(ids) > 0 {
		opts.SetUnackedNotificationIds(stringutil.UniqueTrimmedStrings(ids))
	}

	// Very short reconnection intervals are annoying.
	if opts.HasReconnectionIntervalS() && opts.GetReconnectionIntervalS() < 5 {
		opts.SetReconnectionIntervalS(5)
	}

	if guiPWIsSet && len(opts.GetUnackedNotificationIds()) > 0 {
		for i, key := range opts.GetUnackedNotificationIds() {
			if key == "authenticationUserAndPassword" {
				ids := slices.Delete(opts.GetUnackedNotificationIds(), i, i+1)
				opts.SetUnackedNotificationIds(ids)
				break
			}
		}
	}

	// Negative limits are meaningless, zero meaning no limit.
	if opts.HasConnectionLimitEnough() && opts.GetConnectionLimitEnough() < 0 {
		opts.ClearConnectionLimitEnough()
	}
	if opts.HasConnectionLimitMax() && opts.GetConnectionLimitMax() < 0 {
		opts.ClearConnectionLimitMax()
	}

	if opts.GetConnectionPriorityQuicWan() <= opts.GetConnectionPriorityQuicLan() {
		opts.SetConnectionPriorityQuicWan(opts.GetConnectionPriorityQuicLan() + 1)
	}
	if opts.GetConnectionPriorityTcpWan() <= opts.GetConnectionPriorityTcpLan() {
		opts.SetConnectionPriorityTcpWan(opts.GetConnectionPriorityTcpLan() + 1)
	}

	// If usage reporting is enabled we must have a unique ID.
	if opts.GetUrAccepted() > 0 && opts.GetUrUniqueId() == "" {
		opts.SetUrUniqueId(rand.String(8))
	}
}

// prepareIgnoredDevices removes the ignored devices that are also
// present in the device list.
func prepareIgnoredDevices(cfg *Configuration, existingDevices map[protocol.DeviceID]DeviceConfiguration) {
	ignored := cfg.GetRemoteIgnoredDevices()
	kept := make([]ObservedDevice, 0, len(ignored))
	for _, device := range ignored {
		if _, ok := existingDevices[device.GetDeviceId()]; !ok {
			kept = append(kept, device)
		}
	}
	cfg.SetRemoteIgnoredDevices(kept)
}

// prepareDefaults prepares the folder and device templates, clearing the
// fields that are meaningless in a template.
func prepareDefaults(cfg *Configuration, myID protocol.DeviceID, existingDevices map[protocol.DeviceID]DeviceConfiguration) {
	defaults := cfg.GetDefaults()
	if defaults.Defaults == nil {
		return
	}
	if folder := defaults.GetFolder(); folder.FolderConfiguration != nil {
		folder.ClearId()
		prepareFolder(folder, myID, existingDevices)
	}
	if device := defaults.GetDevice(); device.DeviceConfiguration != nil {
		device.ClearDeviceId()
		device.ClearIntroducedBy()
		prepareDevice(device, nil)
	}
}

// deduplicateObservedFolders returns the observed folders, deduplicated
// by folder ID keeping the most recently seen entry, sorted by time.
func deduplicateObservedFolders(input []*syncthingv2.ObservedFolder) []*syncthingv2.ObservedFolder {
	byID := make(map[string]*syncthingv2.ObservedFolder, len(input))
	for _, folder := range input {
		if existing, ok := byID[folder.GetId()]; !ok || existing.GetTime().AsTime().Before(folder.GetTime().AsTime()) {
			byID[folder.GetId()] = folder
		}
	}
	output := make([]*syncthingv2.ObservedFolder, 0, len(byID))
	for _, folder := range byID {
		output = append(output, folder)
	}
	slices.SortFunc(output, func(a, b *syncthingv2.ObservedFolder) int {
		return a.GetTime().AsTime().Compare(b.GetTime().AsTime())
	})
	return output
}
