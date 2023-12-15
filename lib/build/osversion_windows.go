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
	v := windows.RtlGetVersion()
	return fmt.Sprintf("windows%dp%db%d", v.MajorVersion, v.MinorVersion, v.BuildNumber) // windows10p0b22000
}
