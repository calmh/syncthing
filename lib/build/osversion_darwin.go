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
	var parts []string
	for _, line := range strings.Split(string(out), "\n") {
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		switch k {
		case "ProductName":
			parts = append(parts, strings.TrimSpace(v))
		case "ProductVersion":
			parts = append(parts, strings.ReplaceAll(strings.TrimSpace(v), ".", "p"))
		}
	}
	return strings.ToLower(strings.Join(parts, "")) // macos14p2
}
