// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"fmt"
)

// defaultNumConnections is the number of connections to use by default.
const defaultNumConnections = 3

// NumConnections returns the number of connections to use for this
// device: zero means the default, and negative values mean a single
// connection.
func (d DeviceConfiguration) NumConnections() int {
	switch {
	case d.GetNumConnections() == 0:
		return defaultNumConnections
	case d.GetNumConnections() < 0:
		return 1
	default:
		return int(d.GetNumConnections())
	}
}

// IgnoredFolder returns whether the given folder is in the list of
// folders ignored from this device.
func (d DeviceConfiguration) IgnoredFolder(folder string) bool {
	for _, ignoredFolder := range d.GetIgnoredFolders() {
		if ignoredFolder.GetId() == folder {
			return true
		}
	}
	return false
}

// Description returns a human readable representation of the device.
func (d DeviceConfiguration) Description() string {
	if d.GetName() == "" {
		return d.GetDeviceId().Short().String()
	}
	return fmt.Sprintf("%s (%s)", d.GetName(), d.GetDeviceId().Short())
}
