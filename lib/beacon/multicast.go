// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package beacon

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"time"

	"github.com/syncthing/syncthing/internal/slogutil"
	"github.com/syncthing/syncthing/lib/build"
	"github.com/syncthing/syncthing/lib/netutil"

	"golang.org/x/net/ipv6"
)

func NewMulticast(addr string) Interface {
	c := newCast("multicastBeacon")
	c.addReader(func(ctx context.Context) error {
		return readMulticasts(ctx, c.outbox, addr)
	})
	c.addWriter(func(ctx context.Context) error {
		return writeMulticasts(ctx, c.inbox, addr)
	})
	return c
}

func writeMulticasts(ctx context.Context, inbox <-chan []byte, addr string) error {
	gaddr, err := net.ResolveUDPAddr("udp6", addr)
	if err != nil {
		slog.DebugContext(ctx, "Failed to resolve UDP address", slogutil.Error(err))
		return err
	}

	conn, err := net.ListenPacket("udp6", ":0")
	if err != nil {
		slog.DebugContext(ctx, "Failed to listen on UDP", slogutil.Error(err))
		return err
	}
	doneCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-doneCtx.Done()
		conn.Close()
	}()

	pconn := ipv6.NewPacketConn(conn)

	wcm := &ipv6.ControlMessage{
		HopLimit: 1,
	}

	for {
		var bs []byte
		select {
		case bs = <-inbox:
		case <-doneCtx.Done():
			return doneCtx.Err()
		}

		intfs, err := netutil.Interfaces()
		if err != nil {
			slog.DebugContext(ctx, "Failed to list interfaces", slogutil.Error(err))
			return err
		}

		success := 0
		for _, intf := range intfs {
			if intf.Flags&net.FlagRunning == 0 || intf.Flags&net.FlagMulticast == 0 {
				continue
			}

			if build.IsAndroid && intf.Flags&net.FlagPointToPoint != 0 {
				// skip  cellular interfaces
				continue
			}

			wcm.IfIndex = intf.Index
			pconn.SetWriteDeadline(time.Now().Add(time.Second))
			_, err = pconn.WriteTo(bs, wcm, gaddr)
			pconn.SetWriteDeadline(time.Time{})

			if err != nil {
				slog.DebugContext(ctx, "Write error", slogutil.Address(gaddr), slog.String("intf", intf.Name), slogutil.Error(err))
				continue
			}

			slog.DebugContext(ctx, "Sent multicast", slog.Int("bytes", len(bs)), slogutil.Address(gaddr), slog.String("intf", intf.Name))

			success++

			select {
			case <-doneCtx.Done():
				return doneCtx.Err()
			default:
			}
		}

		if success == 0 {
			return err
		}
	}
}

func readMulticasts(ctx context.Context, outbox chan<- recv, addr string) error {
	gaddr, err := net.ResolveUDPAddr("udp6", addr)
	if err != nil {
		slog.DebugContext(ctx, "Failed to resolve UDP address", slogutil.Error(err))
		return err
	}

	conn, err := net.ListenPacket("udp6", addr)
	if err != nil {
		slog.DebugContext(ctx, "Failed to listen on UDP", slogutil.Error(err))
		return err
	}
	doneCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-doneCtx.Done()
		conn.Close()
	}()

	intfs, err := netutil.Interfaces()
	if err != nil {
		slog.DebugContext(ctx, "Failed to list interfaces", slogutil.Error(err))
		return err
	}

	pconn := ipv6.NewPacketConn(conn)
	joined := 0
	for _, intf := range intfs {
		if intf.Flags&net.FlagRunning == 0 || intf.Flags&net.FlagMulticast == 0 {
			continue
		}

		if build.IsAndroid && intf.Flags&net.FlagPointToPoint != 0 {
			// skip  cellular interfaces
			continue
		}

		err := pconn.JoinGroup(&intf, &net.UDPAddr{IP: gaddr.IP})
		if err != nil {
			slog.DebugContext(ctx, "IPv6 join failed", slog.String("intf", intf.Name), slogutil.Error(err))
		} else {
			slog.DebugContext(ctx, "IPv6 join success", slog.String("intf", intf.Name))
		}
		joined++
	}

	if joined == 0 {
		slog.DebugContext(ctx, "No multicast interfaces available")
		return errors.New("no multicast interfaces available")
	}

	bs := make([]byte, 65536)
	for {
		select {
		case <-doneCtx.Done():
			return doneCtx.Err()
		default:
		}
		n, _, addr, err := pconn.ReadFrom(bs)
		if err != nil {
			slog.DebugContext(ctx, "Read error", slogutil.Error(err))
			return err
		}
		slog.DebugContext(ctx, "Received multicast", slog.Int("bytes", n), slogutil.Address(addr))

		c := make([]byte, n)
		copy(c, bs)
		select {
		case outbox <- recv{c, addr}:
		default:
			slog.DebugContext(ctx, "Dropping message")
		}
	}
}
