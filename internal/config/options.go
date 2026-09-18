// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"fmt"
	"log/slog"
	"net"
	"runtime"
	"slices"
	"strings"

	"github.com/syncthing/syncthing/internal/slogutil"
	"github.com/syncthing/syncthing/lib/protocol"
	"github.com/syncthing/syncthing/lib/rand"
	"github.com/syncthing/syncthing/lib/stringutil"
)

// EffectiveListenAddresses returns the listen addresses with the
// "default" entries expanded to the default set of addresses.
func (o OptionsConfiguration) EffectiveListenAddresses() []string {
	addresses := expandDefaults(o.GetListenAddresses(), func() []string {
		return DefaultListenAddresses
	})
	return stringutil.UniqueTrimmedStrings(addresses)
}

// EffectiveGlobalAnnounceServers returns the global announce (discovery)
// servers with the "default", "default-v4" and "default-v6" entries
// expanded to the default servers.
func (o OptionsConfiguration) EffectiveGlobalAnnounceServers() []string {
	var servers []string
	entries := o.GetGlobalAnnounceServers()
	if len(entries) == 0 {
		entries = []string{"default"}
	}
	for _, srv := range entries {
		switch srv {
		case "default":
			servers = append(servers, DefaultDiscoveryServers...)
		case "default-v4":
			servers = append(servers, DefaultDiscoveryServersV4...)
		case "default-v6":
			servers = append(servers, DefaultDiscoveryServersV6...)
		default:
			servers = append(servers, srv)
		}
	}
	return stringutil.UniqueTrimmedStrings(servers)
}

// EffectiveStunServers returns the servers to use for STUN, with the
// "default" entries expanded to the public STUN servers, preferring
// those hosted by the Syncthing project and resolved via DNS.
func (o OptionsConfiguration) EffectiveStunServers() []string {
	var addresses []string
	entries := o.GetStunServers()
	if len(entries) == 0 {
		entries = []string{"default"}
	}
	for _, addr := range entries {
		switch addr {
		case "default":
			_, records, err := net.LookupSRV("stun", "udp", "syncthing.net")
			if err != nil {
				slog.Debug("Unable to resolve primary STUN servers via DNS", slogutil.Error(err))
			}

			for _, record := range records {
				priority := record.Priority
				target := strings.TrimSuffix(record.Target, ".")
				address := fmt.Sprintf("%s:%d", target, record.Port)
				slog.Debug("Resolved primary STUN server", "address", address, "priority", priority)
				addresses = append(addresses, address)
			}

			fallbackAddresses := slices.Clone(DefaultFallbackStunServers)
			rand.Shuffle(fallbackAddresses)
			addresses = append(addresses, fallbackAddresses...)
		default:
			addresses = append(addresses, addr)
		}
	}

	return stringutil.UniqueTrimmedStrings(addresses)
}

// expandDefaults expands "default" entries in the list, returning the
// given default list if the input is empty.
func expandDefaults(entries []string, defaults func() []string) []string {
	if len(entries) == 0 {
		entries = []string{"default"}
	}
	var out []string
	for _, entry := range entries {
		if entry == "default" {
			out = append(out, defaults()...)
			continue
		}
		out = append(out, entry)
	}
	return out
}

// IsStunDisabled returns whether contacting STUN servers is disabled.
func (o OptionsConfiguration) IsStunDisabled() bool {
	return o.GetStunKeepaliveMinS() < 1 || o.GetStunKeepaliveStartS() < 1 || !o.GetNatEnabled()
}

// MaxFolderConcurrency returns how many folders may concurrently be in
// I/O-intensive operations: a set value if positive, unlimited if
// negative, and otherwise the number of CPUs.
func (o OptionsConfiguration) MaxFolderConcurrency() int {
	// If a value is set, trust that.
	if v := o.GetMaxFolderConcurrency(); v > 0 {
		return int(v)
	}
	if o.GetMaxFolderConcurrency() < 0 {
		// -1 etc means unlimited, which in the implementation means zero
		return 0
	}
	// Otherwise default to the number of CPU cores in the system as a
	// rough approximation of system powerfulness.
	if n := runtime.GOMAXPROCS(-1); n > 0 {
		return n
	}
	// We should never get here to begin with, but since we're here
	// let's use some sort of reasonable compromise between the old "no
	// limit" and getting nothing done... (Median number of folders out
	// there at time of writing is two, 95-percentile at 12 folders.)
	return 4 // https://xkcd.com/221/
}

// MaxConcurrentIncomingRequestKiB returns how much data to allow in
// flight in incoming requests: a set value if positive, disabled if
// negative, the default if unset, and never less than the minimum
// required to make progress.
func (o OptionsConfiguration) MaxConcurrentIncomingRequestKiB() int {
	// Negative is disabled, which in limiter land is spelled zero
	if o.GetMaxConcurrentIncomingRequestKiB() < 0 {
		return 0
	}

	if o.GetMaxConcurrentIncomingRequestKiB() == 0 {
		// The default is 256 MiB
		return 256 * 1024 // KiB
	}

	// We can't really do less than a couple of concurrent blocks or
	// we'll pretty much stall completely. Check that an explicit value
	// is large enough.
	const minAllowed = 2 * protocol.MaxBlockSize / 1024
	if o.GetMaxConcurrentIncomingRequestKiB() < minAllowed {
		return minAllowed
	}

	// Roll with it.
	return int(o.GetMaxConcurrentIncomingRequestKiB())
}

// AutoUpgradeEnabled returns whether automatic upgrades are enabled.
func (o OptionsConfiguration) AutoUpgradeEnabled() bool {
	return o.GetAutoUpgradeIntervalH() > 0
}

// FeatureFlag returns whether the named feature flag is set.
func (o OptionsConfiguration) FeatureFlag(name string) bool {
	return slices.Contains(o.GetFeatureFlags(), name)
}

// LowestConnectionLimit is the lower of ConnectionLimitEnough or
// ConnectionLimitMax, or whichever of them is actually set if only one
// of them is set. It's the point where we should stop dialing.
func (o OptionsConfiguration) LowestConnectionLimit() int {
	limit := o.GetConnectionLimitEnough()
	if limit == 0 || (o.GetConnectionLimitMax() != 0 && o.GetConnectionLimitMax() < limit) {
		// It doesn't really make sense to set Max lower than Enough but
		// someone might do it while experimenting and it's easy for us
		// to do the right thing.
		limit = o.GetConnectionLimitMax()
	}
	return int(limit)
}
