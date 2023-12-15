// Copyright (C) 2023 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package build

import (
	"os/exec"
	"runtime"
	"strings"
)

func getOSVersion() string {
	cmd := exec.Command("sw_vers")
	out, err := cmd.Output()
	if err != nil {
		return runtime.GOOS
	}

	var name, version string
	for _, line := range strings.Split(string(out), "\n") {
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		switch k {
		case "ProductName":
			name = strings.TrimSpace(v)
		case "ProductVersion":
			version = strings.ReplaceAll(strings.TrimSpace(v), ".", "p")
		}
	}
	return strings.ToLower(name + version) // macos14p2
}
