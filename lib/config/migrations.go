// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"cmp"
	"slices"
	"sync"
)

// migrations is the set of config migration functions, with their target
// config version. The conversion function can be nil in which case we just
// update the config version. The order of migrations doesn't matter here,
// put the newest on top for readability.
var (
	migrations = migrationSet{
		{51, migrateToConfigV51},
		{50, migrateToConfigV50},
	}
	migrationsMut = sync.Mutex{}
)

type migrationSet []migration

// apply applies all the migrations in the set, as required by the current
// version and target version, in the correct order.
func (ms migrationSet) apply(cfg *Configuration) {
	// Make sure we apply the migrations in target version order regardless
	// of how it was defined.
	slices.SortFunc(ms, func(a, b migration) int {
		return cmp.Compare(a.targetVersion, b.targetVersion)
	})

	// Apply all migrations.
	for _, m := range ms {
		m.apply(cfg)
	}
}

// A migration is a target config version and a function to do the needful
// to reach that version. The function does not need to change the actual
// cfg.Version field.
type migration struct {
	targetVersion int
	convert       func(cfg *Configuration)
}

// apply applies the conversion function if the current version is below the
// target version and the function is not nil, and updates the current
// version.
func (m migration) apply(cfg *Configuration) {
	if cfg.Version >= m.targetVersion {
		return
	}
	if m.convert != nil {
		m.convert(cfg)
	}
	cfg.Version = m.targetVersion
}

func migrateToConfigV51(cfg *Configuration) {
	oldDefault := 2
	for i, fcfg := range cfg.Folders {
		if fcfg.MaxConcurrentWrites == oldDefault {
			cfg.Folders[i].MaxConcurrentWrites = maxConcurrentWritesDefault
		}
	}
	if cfg.Defaults.Folder.MaxConcurrentWrites == oldDefault {
		cfg.Defaults.Folder.MaxConcurrentWrites = maxConcurrentWritesDefault
	}
}

func migrateToConfigV50(cfg *Configuration) {
	// v50 is Syncthing 2.0
}
