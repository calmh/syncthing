// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"fmt"
	"net"
	"strconv"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/lib/netutil"
	"github.com/syncthing/syncthing/lib/protocol"
)

// New returns a new, empty configuration for the given device, prepared
// for use.
func New(myID protocol.DeviceID) Configuration {
	cfg := Configuration{syncthingv2.Configuration_builder{
		Version: new(int32(CurrentVersion)),
		Options: syncthingv2.OptionsConfiguration_builder{
			UnackedNotificationIds: []string{"authenticationUserAndPassword"},
		}.Build(),
	}.Build()}

	// Can't happen.
	if err := Prepare(&cfg, myID); err != nil {
		panic("bug: error in preparing new configuration: " + err.Error())
	}

	return cfg
}

// ProbeFreePorts probes for free ports to use for the GUI and sync
// listeners, setting them in the configuration if the defaults are
// unavailable.
func (c Configuration) ProbeFreePorts() error {
	if c.GetGui().Network() == "tcp" {
		guiHost, guiPort, err := net.SplitHostPort(c.GetGui().Address())
		if err != nil {
			return fmt.Errorf("get default port (GUI): %w", err)
		}
		port, err := strconv.Atoi(guiPort)
		if err != nil {
			return fmt.Errorf("convert default port (GUI): %w", err)
		}
		port, err = getFreePort(guiHost, port)
		if err != nil {
			return fmt.Errorf("get free port (GUI): %w", err)
		}
		gui := c.GetGui()
		if gui.GUIConfiguration == nil {
			gui = GUIConfiguration{syncthingv2.GUIConfiguration_builder{}.Build()}
			c.SetGui(gui)
		}
		gui.SetAddress(net.JoinHostPort(guiHost, strconv.Itoa(port)))
	}

	port, err := getFreePort("0.0.0.0", DefaultTCPPort)
	if err != nil {
		return fmt.Errorf("get free port (BEP): %w", err)
	}
	opts := c.GetOptions()
	if opts.OptionsConfiguration == nil {
		opts = OptionsConfiguration{syncthingv2.OptionsConfiguration_builder{}.Build()}
		c.SetOptions(opts)
	}
	if port == DefaultTCPPort {
		opts.SetListenAddresses([]string{"default"})
	} else {
		opts.SetListenAddresses([]string{
			netutil.AddressURL("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(port))),
			"dynamic+https://relays.syncthing.net/endpoint",
			netutil.AddressURL("quic", net.JoinHostPort("0.0.0.0", strconv.Itoa(port))),
		})
	}

	return nil
}

// getFreePort tries to listen on the given ports in order, returning
// the first that succeeds. If none succeed, a random high port is
// returned.
func getFreePort(host string, ports ...int) (int, error) {
	for _, port := range ports {
		c, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
		if err == nil {
			c.Close()
			return port, nil
		}
	}

	c, err := net.Listen("tcp", host+":0")
	if err != nil {
		return 0, err
	}
	addr := c.Addr().(*net.TCPAddr)
	c.Close()
	return addr.Port, nil
}
