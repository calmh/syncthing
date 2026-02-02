// Copyright (C) 2017 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package syncthing

import (
	"log/slog"
	"syscall"

	"github.com/syncthing/syncthing/internal/slogutil"
)

// https://docs.microsoft.com/windows/win32/secauthz/well-known-sids
const securityLocalSystemRID = "S-1-5-18"

func isSuperUser() bool {
	tok, err := syscall.OpenCurrentProcessToken()
	if err != nil {
		slog.Debug("OpenCurrentProcessToken failed", slogutil.Error(err))
		return false
	}
	defer tok.Close()

	user, err := tok.GetTokenUser()
	if err != nil {
		slog.Debug("GetTokenUser failed", slogutil.Error(err))
		return false
	}

	if user.User.Sid == nil {
		slog.Debug("Sid is nil")
		return false
	}

	sid, err := user.User.Sid.String()
	if err != nil {
		slog.Debug("Sid.String failed", slogutil.Error(err))
		return false
	}

	slog.Debug("Got SID", slog.String("sid", sid))
	return sid == securityLocalSystemRID
}
