// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"fmt"

	"github.com/syncthing/syncthing/lib/fs"
)

// SizeInBytes returns the size expressed in bytes. Percentages and unset
// sizes have no byte value and return zero.
func SizeInBytes(s *Size) float64 {
	switch {
	case s.HasBytes():
		return s.GetBytes()
	case s.HasMib():
		return s.GetMib() * (1 << 20)
	case s.HasGib():
		return s.GetGib() * (1 << 30)
	default:
		return 0
	}
}

// SizeIsPercentage reports whether the size is a percentage of total
// capacity rather than an absolute size.
func SizeIsPercentage(s *Size) bool {
	return s.HasPercent()
}

// CheckFreeSpace checks that the free space does not fall below the
// minimum required free space. A nil size, or a zero size, means no
// requirement.
func CheckFreeSpace(minFree *Size, usage fs.Usage) error {
	if minFree == nil {
		return nil
	}

	if SizeIsPercentage(minFree) {
		val := minFree.GetPercent()
		if val <= 0 {
			return nil
		}
		freePct := (float64(usage.Free) / float64(usage.Total)) * 100
		if freePct < val {
			return fmt.Errorf("current %.2f %% < required %v", freePct, sizeString(minFree))
		}
		return nil
	}

	val := SizeInBytes(minFree)
	if val <= 0 {
		return nil
	}
	if float64(usage.Free) < val {
		return fmt.Errorf("current %sB < required %v", formatSI(usage.Free), sizeString(minFree))
	}

	return nil
}

// checkAvailableSpace checks that the free space does not fall below the
// minimum required free space, considering additional required space
// for a future operation.
func checkAvailableSpace(req uint64, minFree *Size, usage fs.Usage) error {
	if usage.Free < req {
		return fmt.Errorf("current %sB < required %sB", formatSI(usage.Free), formatSI(req))
	}
	usage.Free -= req
	return CheckFreeSpace(minFree, usage)
}

// sizeString returns a human readable representation of the size.
func sizeString(s *Size) string {
	switch {
	case s.HasPercent():
		return fmt.Sprintf("%v %%", s.GetPercent())
	case s.HasMib():
		return fmt.Sprintf("%v MiB", s.GetMib())
	case s.HasGib():
		return fmt.Sprintf("%v GiB", s.GetGib())
	default:
		return fmt.Sprintf("%v B", s.GetBytes())
	}
}

func formatSI(b uint64) string {
	switch {
	case b < 1000:
		return fmt.Sprintf("%d ", b)
	case b < 1000*1000:
		return fmt.Sprintf("%.1f K", float64(b)/1000)
	case b < 1000*1000*1000:
		return fmt.Sprintf("%.1f M", float64(b)/(1000*1000))
	case b < 1000*1000*1000*1000:
		return fmt.Sprintf("%.1f G", float64(b)/(1000*1000*1000))
	default:
		return fmt.Sprintf("%.1f T", float64(b)/(1000*1000*1000*1000))
	}
}
