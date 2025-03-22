// Copyright (C) 2025 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build !cgo && !wazero

package sqlite

import (
	"context"
	"database/sql/driver"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/syncthing/syncthing/lib/build"
	"modernc.org/sqlite"
)

const (
	dbDriver      = "sqlite"
	commonOptions = "_pragma=foreign_keys(1)&_pragma=recursive_triggers(1)&_pragma=cache_size(-65536)&_pragma=synchronous(1)"
)

func init() {
	build.AddTag("modernc-sqlite")
	sqlite.RegisterConnectionHook(func(conn sqlite.ExecQuerierContext, dsn string) error {
		uri, err := url.Parse(filepath.ToSlash(dsn))
		if err != nil {
			return err
		}
		mainFile := uri.Path
		mainExt := filepath.Ext(mainFile)
		withoutExt := strings.TrimSuffix(mainFile, mainExt)
		blocksFile := filepath.FromSlash(withoutExt + "-blocks" + mainExt)
		_, err = conn.ExecContext(context.Background(), `ATTACH ? AS blocks`, []driver.NamedValue{{Ordinal: 1, Value: blocksFile}})
		return err
	})
}
