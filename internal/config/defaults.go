// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"net"
	"strconv"

	"github.com/syncthing/syncthing/lib/netutil"
)

const (
	// OldestHandledVersion is the oldest config version that we can
	// read. Configs older than this get a best-effort attempt at
	// conversion.
	OldestHandledVersion = 10

	// CurrentVersion is the current version of the configuration
	// format. It increments whenever a change is made that requires
	// migration from previous formats.
	CurrentVersion = 52

	// MaxRescanIntervalS is the maximum rescan interval, in seconds.
	MaxRescanIntervalS = 365 * 24 * 60 * 60

	// DefaultTCPPort defines the default TCP port used if the URI does
	// not specify one, for example tcp://0.0.0.0
	DefaultTCPPort = 22000

	// DefaultQUICPort defines the default QUIC port used if the URI does
	// not specify one, for example quic://0.0.0.0
	DefaultQUICPort = 22000

	// DefaultTheme is the default and fallback theme for the web UI.
	DefaultTheme = "default"

	// DefaultMarkerName is the default folder marker name.
	DefaultMarkerName = ".stfolder"

	// EncryptionTokenName is the file name used to store the folder
	// encryption password token.
	EncryptionTokenName = "syncthing-encryption_password_token" //nolint: gosec
)

var (
	// DefaultListenAddresses should be substituted when the configuration
	// contains <listenAddress>default</listenAddress>. This is done by the
	// "consumer" of the configuration as we don't want these saved to the
	// config.
	DefaultListenAddresses = []string{
		netutil.AddressURL("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(DefaultTCPPort))),
		"dynamic+https://relays.syncthing.net/endpoint",
		netutil.AddressURL("quic", net.JoinHostPort("0.0.0.0", strconv.Itoa(DefaultQUICPort))),
	}

	// DefaultDiscoveryServersV4 should be substituted when the
	// configuration contains <globalAnnounceServer>default-v4</globalAnnounceServer>.
	DefaultDiscoveryServersV4 = []string{
		"https://discovery-lookup.syncthing.net/v2/?noannounce",
		"https://discovery-announce-v4.syncthing.net/v2/?nolookup",
	}

	// DefaultDiscoveryServersV6 should be substituted when the
	// configuration contains <globalAnnounceServer>default-v6</globalAnnounceServer>.
	DefaultDiscoveryServersV6 = []string{
		"https://discovery-lookup.syncthing.net/v2/?noannounce",
		"https://discovery-announce-v6.syncthing.net/v2/?nolookup",
	}

	// DefaultDiscoveryServers should be substituted when the
	// configuration contains <globalAnnounceServer>default</globalAnnounceServer>.
	DefaultDiscoveryServers = append(DefaultDiscoveryServersV4, DefaultDiscoveryServersV6...)

	// Default fallback STUN servers, used if the primary servers can't be
	// resolved or are down.
	DefaultFallbackStunServers = []string{
		"stun.counterpath.com:3478",
		"stun.hitv.com:3478",
		"stun.internetcalls.com:3478",
		"stun.miwifi.com:3478",
		"stun.schlund.de:3478",
		"stun.sipgate.net:3478",
		"stun.voip.aebc.com:3478",
		"stun.voipbuster.com:3478",
		"stun.voipstunt.com:3478",
	}
)
