// Copyright (C) 2015 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/syncthing/syncthing/lib/db"
	"github.com/syncthing/syncthing/lib/protocol"
)

func dump(base string, ldb *db.Lowlevel) {
	files := make(map[string]*json.Encoder)

	it := ldb.NewIterator(nil, nil)
	for it.Next() {
		out := func(t string, dec map[string]interface{}) {
			enc, ok := files[t]
			if !ok {
				marker := it.Key()[0]
				name := fmt.Sprintf("%s-%02d-%s.jsons", base, marker, t)
				log.Println("Creating", name, "...")
				fd, err := os.Create(name)
				if err != nil {
					log.Fatal("Output:", err)
				}
				enc = json.NewEncoder(fd)
				files[t] = enc
			}
			m := map[string]interface{}{
				".type": t,
				"dec":   dec,
				"key":   it.Key(),
			}
			if dec == nil {
				m["val"] = it.Value()
			}
			enc.Encode(m)
		}
		key := it.Key()
		switch key[0] {
		case db.KeyTypeDevice:
			folder := binary.BigEndian.Uint32(key[1:])
			device := binary.BigEndian.Uint32(key[1+4:])
			name := nulString(key[1+4+4:])

			var f protocol.FileInfo
			if err := f.Unmarshal(it.Value()); err != nil {
				log.Println("Unmarshal file info:", err)
				continue
			}

			out("device", map[string]interface{}{"folder": folder, "device": device, "name": name, "value": f})

		case db.KeyTypeGlobal:
			folder := binary.BigEndian.Uint32(key[1:])
			name := nulString(key[1+4:])
			var flv db.VersionList
			if err := flv.Unmarshal(it.Value()); err != nil {
				log.Println("Unmarshal global:", err)
				continue
			}
			out("global", map[string]interface{}{"folder": folder, "name": name, "value": flv})

		case db.KeyTypeBlock:
			folder := binary.BigEndian.Uint32(key[1:])
			hash := key[1+4 : 1+4+32]
			name := nulString(key[1+4+32:])
			out("block", map[string]interface{}{"folder": folder, "name": name, "hash": hash, "value": binary.BigEndian.Uint32(it.Value())})

		case db.KeyTypeDeviceStatistic:
			out("dstat", nil)

		case db.KeyTypeFolderStatistic:
			out("fstat", nil)

		case db.KeyTypeVirtualMtime:
			folder := binary.BigEndian.Uint32(key[1:])
			name := nulString(key[1+4:])
			val := it.Value()
			var real, virt time.Time
			real.UnmarshalBinary(val[:len(val)/2])
			virt.UnmarshalBinary(val[len(val)/2:])
			out("mtime", map[string]interface{}{"folder": folder, "name": name, "real": real, "virt": virt})

		case db.KeyTypeFolderIdx:
			key := binary.BigEndian.Uint32(it.Key()[1:])
			out("folderidx", map[string]interface{}{"idx": key, "value": string(it.Value())})

		case db.KeyTypeDeviceIdx:
			key := binary.BigEndian.Uint32(it.Key()[1:])
			val := it.Value()
			if len(val) == 0 {
				out("deviceidx", map[string]interface{}{"idx": key, "value": nil})
			} else {
				dev := protocol.DeviceIDFromBytes(val).String()
				out("deviceidx", map[string]interface{}{"idx": key, "value": dev})
			}

		case db.KeyTypeIndexID:
			key := it.Key()
			devIdx := binary.BigEndian.Uint32(key[1:])
			fldIdx := binary.BigEndian.Uint32(key[5:])
			id := binary.BigEndian.Uint64(it.Value())
			out("indexid", map[string]interface{}{"device": devIdx, "folder": fldIdx, "id": id})

		case db.KeyTypeFolderMeta:
			key := it.Key()
			fldIdx := binary.BigEndian.Uint32(key[1:])
			out("foldermeta", map[string]interface{}{"folder": fldIdx, "value": it.Value()})

		case db.KeyTypeMiscData:
			out("miscdata", map[string]interface{}{"key": string(it.Key()), "value": string(it.Value())})

		case db.KeyTypeSequence:
			key := it.Key()
			fldIdx := binary.BigEndian.Uint32(key[1:])
			seq := binary.BigEndian.Uint64(key[5:])

			folder := binary.BigEndian.Uint32(key[1:])
			device := binary.BigEndian.Uint32(key[1+4:])
			name := nulString(key[1+4+4:])

			out("sequence", map[string]interface{}{"folderInKey": fldIdx, "sequence": seq, "device": device, "folderInValue": folder, "name": name})

		case db.KeyTypeNeed:
			key := it.Key()
			fldIdx := binary.BigEndian.Uint32(key[1:])
			name := string(key[5:])

			out("need", map[string]interface{}{"folder": fldIdx, "name": name})

		default:
			out("unknown", map[string]interface{}{"marker": it.Key()[0], "value": it.Value()})
		}
	}
}
