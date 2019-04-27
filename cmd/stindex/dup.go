// Copyright (C) 2018 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/syncthing/syncthing/lib/db"
	"github.com/syncthing/syncthing/lib/protocol"
)

func dups(ldb *db.Lowlevel) {
	it := ldb.NewIterator(nil, nil)
	defer it.Release()

	type blockdata struct {
		hash  string
		size  int
		count int
	}

	blocks := make(map[string]blockdata)

	for it.Next() {
		key := it.Key()

		switch key[0] {
		case db.KeyTypeDevice:
			device := binary.BigEndian.Uint32(key[1+4:])
			if device != 1 {
				continue
			}

			var f protocol.FileInfo
			err := f.Unmarshal(it.Value())
			if err != nil {
				fmt.Println("Unable to unmarshal FileInfo:", err)
				continue
			}
			for _, b := range f.Blocks {
				bd := blocks[string(b.Hash)]
				bd.hash = string(b.Hash)
				bd.size = int(b.Size)
				bd.count++
				blocks[string(b.Hash)] = bd
			}
		}
	}

	var tot int64
	var dup int64
	var dups []blockdata
	for _, bd := range blocks {
		if bd.count > 1 {
			dup += int64(bd.count * bd.size)
			dups = append(dups, bd)
		}
		tot += int64(bd.size)
	}

	sort.Slice(dups, func(a, b int) bool { return dups[a].count < dups[b].count })
	for _, bd := range dups {
		fmt.Printf("Block %x (size %d) duplicated %d times\n", bd.hash, bd.size, bd.count)
	}

	fmt.Printf("Total data is %d MiB\n", tot/1024/1024)
	fmt.Printf("Duplicated data is %d MiB (%.01f%%)\n", dup/1024/1024, float64(dup)/float64(tot)*100)
}
