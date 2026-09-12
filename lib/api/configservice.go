// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package api

import (
	"context"

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
func (s *configService) GetConfiguration(_ context.Context, _ *connect.Request[syncthingv2.GetConfigRequest]) (*connect.Response[syncthingv2.GetConfigResponse], error) {
	cfg := newconfig.FromLegacy(s.cfg.RawCopy())
	return connect.NewResponse(syncthingv2.GetConfigResponse_builder{
		Configuration: cfg.Configuration,
	}.Build()), nil
}
