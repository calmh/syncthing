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

func getOSVersion() (*OSVersion, error) {
	v, err := windows.RtlGetVersion()
	if err != nil {
		return runtime.goOS
	}
	return fmt.Sprintf("windows%dp%d", v.MajorVersion, v.MinorVersion),
}
