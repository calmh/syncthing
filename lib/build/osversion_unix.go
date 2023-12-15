//go:build !darwin && !windows

package build

import "runtime"

func getOSVersion() strign {
	return runtime.GOOS
}
