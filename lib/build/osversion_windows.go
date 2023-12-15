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
