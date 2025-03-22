// Copyright (C) 2025 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build cgo

package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"path/filepath"
	"strings"

	"github.com/mattn/go-sqlite3"
)

const (
	dbDriver      = "sqlite3_extended"
	commonOptions = "_fk=true&_rt=true&_cache_size=-65536&_sync=1&_txlock=immediate"
)

func init() {
	sql.Register("sqlite3_extended",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				mainFile := conn.GetFilename("main")
				mainExt := filepath.Ext(mainFile)
				withoutExt := strings.TrimSuffix(mainFile, mainExt)
				blocksFile := withoutExt + "-blocks" + mainExt
				_, err := conn.Exec(`ATTACH ? AS blocks`, []driver.Value{blocksFile})
				return err
			},
		})
}
