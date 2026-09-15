// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package api

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	newconfig "github.com/syncthing/syncthing/internal/config"
	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
	"github.com/syncthing/syncthing/internal/gen/syncthing/v2/syncthingv2connect"
	"github.com/syncthing/syncthing/lib/config"
)

// configSource is the subset of the configuration wrapper used by the
// config service.
type configSource interface {
	RawCopy() config.Configuration
}

// configService implements the ConfigurationService ConnectRPC service,
// serving the configuration in the new format.
type configService struct {
	cfg configSource
}

var _ syncthingv2connect.ConfigurationServiceHandler = (*configService)(nil)

// GetConfiguration returns the current configuration. The configuration
// is converted from the legacy format, which remains the source of truth
// for now.
func (s *configService) GetConfiguration(_ context.Context, req *connect.Request[syncthingv2.GetConfigRequest]) (*connect.Response[syncthingv2.GetConfigResponse], error) {
	cfg := config.FromLegacy(s.cfg.RawCopy())
	if req.Msg.GetMaterializeDefaults() {
		cfg = newconfig.MaterializeDefaults(cfg)
	}
	return connect.NewResponse(syncthingv2.GetConfigResponse_builder{
		Configuration: cfg.Configuration,
	}.Build()), nil
}

// withConnectGetDefaults fills in the query parameters that the Connect
// protocol requires on GET requests, so that a bare GET works: the
// encoding defaults to JSON and the message to the encoding of an empty
// request. Requests that set the parameters are passed through unchanged.
func withConnectGetDefaults(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			query := r.URL.Query()
			if !query.Has("encoding") || !query.Has("message") {
				if !query.Has("encoding") {
					query.Set("encoding", "json")
				}
				if !query.Has("message") {
					if query.Get("encoding") == "json" {
						query.Set("message", "{}")
					} else {
						query.Set("message", "")
					}
				}
				modified := r.Clone(r.Context())
				modified.URL.RawQuery = query.Encode()
				r = modified
			}
		}
		handler.ServeHTTP(w, r)
	})
}
