// Copyright (C) 2023 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package build

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func getOSVersion() string {
	// n.b. The returned major and minor version numbers are not enormously
	// intuitive. There's at table at
	// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/wdm/ns-wdm-_osversioninfoexw.
	v := windows.RtlGetVersion()
	return fmt.Sprintf("windows%dp%db%d", v.MajorVersion, v.MinorVersion) // windows10p0
}
