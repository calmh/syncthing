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
	"github.com/syndtr/goleveldb/leveldb/util"
)

func dups(ldb *db.Lowlevel) {
	var localDevice uint32
	it := ldb.NewIterator(util.BytesPrefix([]byte{db.KeyTypeDeviceIdx}), nil)
	for it.Next() {
		dkey := binary.BigEndian.Uint32(it.Key()[1:])
		val := it.Value()
		if len(val) == 0 {
			continue
		}
		fmt.Println(dkey, protocol.DeviceIDFromBytes(val))
		if protocol.DeviceIDFromBytes(val) == protocol.LocalDeviceID {
			localDevice = dkey
			break
		}
	}
	it.Release()

	type blockdata struct {
		hash  string
		size  int
		count int
	}

	blocks := make(map[string]blockdata)

	it = ldb.NewIterator(util.BytesPrefix([]byte{db.KeyTypeDevice}), nil)
	for it.Next() {
		key := it.Key()
		device := binary.BigEndian.Uint32(key[1+4:])
		if device != localDevice {
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
	it.Release()

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
