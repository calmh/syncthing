// Copyright (C) 2018 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package connections

import configv1 "github.com/syncthing/syncthing/internal/config/v1"

// invalidListener is never valid
type invalidListener struct {
	listenerFactory
	err error
}

func (i invalidListener) Valid(_ configv1.Configuration) error {
	if i.err == nil {
		// fallback so we don't accidentally return nil
		return errUnsupported
	}
	return i.err
}

// invalidDialer is never valid
type invalidDialer struct {
	dialerFactory
	err error
}

func (i invalidDialer) Valid(_ configv1.Configuration) error {
	if i.err == nil {
		// fallback so we don't accidentally return nil
		return errUnsupported
	}
	return i.err
}

func init() {
	listeners["kcp"] = invalidListener{err: errDeprecated}
	listeners["kcp4"] = invalidListener{err: errDeprecated}
	listeners["kcp6"] = invalidListener{err: errDeprecated}
	dialers["kcp"] = invalidDialer{err: errDeprecated}
	dialers["kcp4"] = invalidDialer{err: errDeprecated}
	dialers["kcp6"] = invalidDialer{err: errDeprecated}
}
