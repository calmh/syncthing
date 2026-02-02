// Copyright (C) 2015 Audrius Butkevicius and Contributors (see the CONTRIBUTORS file).

package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"

	"github.com/syncthing/syncthing/internal/slogutil"
	"github.com/syncthing/syncthing/lib/osutil"
	"github.com/syncthing/syncthing/lib/rand"
	"github.com/syncthing/syncthing/lib/relay/protocol"
)

type dynamicClient struct {
	commonClient

	pooladdr *url.URL
	certs    []tls.Certificate
	timeout  time.Duration

	mut    sync.RWMutex // Protects client.
	client *staticClient
}

func newDynamicClient(uri *url.URL, certs []tls.Certificate, invitations chan protocol.SessionInvitation, timeout time.Duration) *dynamicClient {
	c := &dynamicClient{
		pooladdr: uri,
		certs:    certs,
		timeout:  timeout,
	}
	c.commonClient = newCommonClient(invitations, c.serve, fmt.Sprintf("dynamicClient@%p", c))
	return c
}

func (c *dynamicClient) serve(ctx context.Context) error {
	uri := *c.pooladdr

	// Trim off the `dynamic+` prefix
	uri.Scheme = uri.Scheme[8:]

	slog.Debug("Looking up dynamic relays", slog.Any("client", c))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		slog.Debug("Failed to lookup dynamic relays", slog.Any("client", c), slogutil.Error(err))
		return err
	}
	data, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Debug("Failed to lookup dynamic relays", slog.Any("client", c), slogutil.Error(err))
		return err
	}

	var ann dynamicAnnouncement
	err = json.NewDecoder(data.Body).Decode(&ann)
	data.Body.Close()
	if err != nil {
		slog.Debug("Failed to lookup dynamic relays", slog.Any("client", c), slogutil.Error(err))
		return err
	}

	var addrs []string
	for _, relayAnn := range ann.Relays {
		ruri, err := url.Parse(relayAnn.URL)
		if err != nil {
			slog.Debug("Failed to parse dynamic relay address", slog.Any("client", c), slogutil.URI(relayAnn.URL), slogutil.Error(err))
			continue
		}
		slog.Debug("Found relay", slog.Any("client", c), slogutil.URI(ruri.String()))
		addrs = append(addrs, ruri.String())
	}

	for _, addr := range relayAddressesOrder(ctx, addrs) {
		select {
		case <-ctx.Done():
			slog.Debug("Stopping", slog.Any("client", c))
			return nil
		default:
			ruri, err := url.Parse(addr)
			if err != nil {
				slog.Debug("Skipping relay", slog.Any("client", c), slog.String("addr", addr), slogutil.Error(err))
				continue
			}
			client := newStaticClient(ruri, c.certs, c.invitations, c.timeout)
			c.mut.Lock()
			c.client = client
			c.mut.Unlock()

			err = c.client.Serve(ctx)
			slog.Debug("Disconnected from relay", slog.String("scheme", c.client.URI().Scheme), slog.String("host", c.client.URI().Host), slogutil.Error(err))

			c.mut.Lock()
			c.client = nil
			c.mut.Unlock()
		}
	}
	slog.Debug("Could not find a connectable relay", slog.Any("client", c))
	return errors.New("could not find a connectable relay")
}

func (c *dynamicClient) Error() error {
	c.mut.RLock()
	defer c.mut.RUnlock()
	if c.client == nil {
		return c.commonClient.Error()
	}
	return c.client.Error()
}

func (c *dynamicClient) String() string {
	return fmt.Sprintf("DynamicClient:%p:%s@%s", c, c.URI(), c.pooladdr)
}

func (c *dynamicClient) URI() *url.URL {
	c.mut.RLock()
	defer c.mut.RUnlock()
	if c.client == nil {
		return nil
	}
	return c.client.URI()
}

// This is the announcement received from the relay server;
// {"relays": [{"url": "relay://10.20.30.40:5060"}, ...]}
type dynamicAnnouncement struct {
	Relays []struct {
		URL string
	}
}

// relayAddressesOrder checks the latency to each relay, rounds latency down to
// the closest 50ms, and puts them in buckets of 50ms latency ranges. Then
// shuffles each bucket, and returns all addresses starting with the ones from
// the lowest latency bucket, ending with the highest latency bucket.
func relayAddressesOrder(ctx context.Context, input []string) []string {
	buckets := make(map[int][]string)

	for _, relay := range input {
		latency, err := osutil.GetLatencyForURL(ctx, relay)
		if err != nil {
			latency = time.Hour
		}

		id := int(latency/time.Millisecond) / 50

		buckets[id] = append(buckets[id], relay)

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}

	var ids []int
	for id, bucket := range buckets {
		rand.Shuffle(bucket)
		ids = append(ids, id)
	}

	slices.Sort(ids)

	addresses := make([]string, 0, len(input))
	for _, id := range ids {
		addresses = append(addresses, buckets[id]...)
	}

	return addresses
}
