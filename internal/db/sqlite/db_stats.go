// Copyright (C) 2025 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package sqlite

import (
	"cmp"
	"slices"
)

type DatabaseStatistics struct {
	Name     string               `json:"name"`
	FolderID string               `json:"folderID,omitempty"`
	Tables   []TableStatistics    `json:"tables"`
	Total    TableStatistics      `json:"total"`
	Children []DatabaseStatistics `json:"children,omitempty"`
}

type TableStatistics struct {
	Name   string `json:"name,omitempty"`
	Size   int64  `json:"size"`
	Unused int64  `json:"unused"`
}

func (s *DB) Statistics() (*DatabaseStatistics, error) {
	ts, total, err := s.tableStats()
	if err != nil {
		return nil, wrap(err)
	}
	ds := DatabaseStatistics{
		Name:   s.baseName,
		Tables: ts,
		Total:  total,
	}

	err = s.forEachFolder(func(fdb *folderDB) error {
		tables, total, err := fdb.tableStats()
		if err != nil {
			return wrap(err)
		}
		ds.Children = append(ds.Children, DatabaseStatistics{
			Name:     fdb.baseName,
			FolderID: fdb.folderID,
			Tables:   tables,
			Total:    total,
		})
		return nil
	})
	if err != nil {
		return nil, wrap(err)
	}

	return &ds, nil
}

func (s *baseDB) tableStats() ([]TableStatistics, TableStatistics, error) {
	var stats []TableStatistics
	if err := s.stmt(`
		SELECT name, pgsize AS size, unused FROM dbstat
		WHERE aggregate=true
		ORDER BY name
	`).Select(&stats); err != nil {
		return nil, TableStatistics{}, wrap(err)
	}
	var total TableStatistics
	for _, s := range stats {
		total.Size += s.Size
		total.Unused += s.Unused
	}
	return stats, total, nil
}

type DuplicateFile struct {
	Folder string
	Size   int
	Count  int
	Names  []string
}

func (s *DB) DuplicateFiles() ([]DuplicateFile, error) {
	var res []DuplicateFile
	err := s.forEachFolder(func(fdb *folderDB) error {
		var blocklists []struct {
			Hash  []byte `db:"blocklist_hash"`
			Count int
		}
		err := fdb.stmt(`
			SELECT blocklist_hash, count(*) as count FROM files
			WHERE device_idx = {{.LocalDeviceIdx}} AND NOT deleted AND blocklist_hash IS NOT NULL
			GROUP BY blocklist_hash
			HAVING count > 1
		`).Select(&blocklists)
		if err != nil {
			return err
		}
		for _, bl := range blocklists {
			var files []struct {
				Name string
				Size int
			}
			err := fdb.stmt(`
				SELECT name, size FROM files
				WHERE device_idx = {{.LocalDeviceIdx}} AND blocklist_hash = ?
			`).Select(&files, bl.Hash)
			if err != nil {
				return err
			}

			d := DuplicateFile{
				Folder: fdb.folderID,
				Size:   files[0].Size,
				Count:  bl.Count,
			}
			for _, f := range files {
				d.Names = append(d.Names, f.Name)
			}
			res = append(res, d)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	slices.SortFunc(res, func(a, b DuplicateFile) int {
		return cmp.Compare(a.Count*a.Size, b.Count*b.Size)
	})
	return res, nil
}
