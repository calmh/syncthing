// Copyright (C) 2021 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package model

import (
	"bytes"
	"sort"
	"time"

	"github.com/syncthing/syncthing/lib/db"
	"github.com/syncthing/syncthing/lib/protocol"
	"github.com/syncthing/syncthing/lib/sync"
)

const (
	dbKey        = "allFolderDevices"
	keepDuration = 14 * 24 * time.Hour
)

type AllKnownFolderDevicesStore interface {
	KnownDevicesForFolder(folderID string) []protocol.DeviceID
	KnownFoldersForDevice(deviceID protocol.DeviceID) []string
}

type allKnownFolderDevicesStore struct {
	kv  *db.NamespacedKV
	mut sync.RWMutex
	cur map[string]map[protocol.DeviceID]int64 // folder -> device -> seen
}

func newAllKnownFolderDevicesStore(kv *db.NamespacedKV) *allKnownFolderDevicesStore {
	s := &allKnownFolderDevicesStore{
		kv:  kv,
		mut: sync.NewRWMutex(),
	}
	_ = s.load()
	return s
}

func (s *allKnownFolderDevicesStore) AddKnownFolderDevices(cc protocol.ClusterConfig) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	now := time.Now().Unix()
	for _, ccFolder := range cc.Folders {
		candFolder, ok := s.cur[ccFolder.ID]
		if !ok {
			candFolder = make(map[protocol.DeviceID]int64)
			s.cur[ccFolder.ID] = candFolder
		}
		for _, dev := range ccFolder.Devices {
			candFolder[dev.ID] = now
		}
	}
	return s.save()
}

func (s *allKnownFolderDevicesStore) KnownDevicesForFolder(folderID string) []protocol.DeviceID {
	s.mut.Lock()
	defer s.mut.Unlock()

	devs := s.cur[folderID]
	ids := make([]protocol.DeviceID, 0, len(devs))
	cutoff := time.Now().Add(-keepDuration).Unix()
	for devID, when := range devs {
		if when < cutoff {
			continue
		}
		ids = append(ids, devID)
	}
	sort.Slice(ids, func(a, b int) bool {
		return bytes.Compare(ids[a][:], ids[b][:]) == -1
	})
	return ids
}

func (s *allKnownFolderDevicesStore) KnownFoldersForDevice(deviceID protocol.DeviceID) []string {
	s.mut.RLock()
	defer s.mut.RUnlock()

	var folders []string
	cutoff := time.Now().Add(-keepDuration).Unix()
	for folderID, devices := range s.cur {
		if when := devices[deviceID]; when >= cutoff {
			folders = append(folders, folderID)
		}
	}
	sort.Strings(folders)
	return folders
}

func (s *allKnownFolderDevicesStore) save() error {
	cands := SharingCandidates{
		Folders: make([]SharingFolderEntry, 0, len(s.cur)),
	}
	cutoff := time.Now().Add(-keepDuration).Unix()
	for folder, devices := range s.cur {
		fdl := SharingFolderEntry{
			Folder:  folder,
			Devices: make([]SharingDeviceEntry, 0, len(devices)),
		}
		for dev, when := range devices {
			if when < cutoff {
				continue
			}
			fdl.Devices = append(fdl.Devices, SharingDeviceEntry{
				Device: dev[:],
				When:   when,
			})
		}
		cands.Folders = append(cands.Folders, fdl)
	}
	bs, err := cands.Marshal()
	if err != nil {
		return err
	}
	return s.kv.PutBytes(dbKey, bs)
}

func (s *allKnownFolderDevicesStore) load() error {
	bs, ok, err := s.kv.Bytes(dbKey)
	if err != nil {
		return err
	}

	s.cur = make(map[string]map[protocol.DeviceID]int64)
	if !ok {
		return nil
	}

	var cands SharingCandidates
	if err := cands.Unmarshal(bs); err != nil {
		return err
	}

	for _, folder := range cands.Folders {
		devs := make(map[protocol.DeviceID]int64)
		for _, dev := range folder.Devices {
			devID, err := protocol.DeviceIDFromBytes(dev.Device)
			if err != nil {
				continue
			}
			devs[devID] = dev.When
		}
		s.cur[folder.Folder] = devs
	}

	return nil
}

// // devicesToBytes returns a list of device IDs in byte slice form, given a
// // list of devices in a cluster config. The returned list is sorted.
// func devicesToBytes(devs []protocol.Device) [][]byte {
// 	res := make([][]byte, len(devs))
// 	for i := range devs {
// 		res[i] = devs[i].ID[:]
// 	}
// 	sort.Slice(res, func(a, b int) bool {
// 		return bytes.Compare(res[a], res[b]) == -1
// 	})
// 	return res
// }

// // mergeLists merges src into dst without duplicates, assuming both lists
// // are already sorted.
// func mergeLists(dst, src [][]byte) [][]byte {
// 	di, si := 0, 0
// 	for {
// 		if si == len(src) {
// 			break
// 		}
// 		if di == len(dst) {
// 			// We've reached the end of dst, merge becomes append
// 			dst = append(dst, src[si])
// 			si++
// 			continue
// 		}

// 		switch bytes.Compare(src[si], dst[di]) {
// 		case -1:
// 			// src[si] should come before dst[di], so we insert here
// 			dst = append(dst[:di], append([][]byte{src[si]}, dst[di:]...)...)
// 			si++
// 		case 0:
// 			// src[si] is already present in dst, look for next src element
// 			si++
// 		case 1:
// 			// src[i] should come later, look further in dst
// 			di++
// 		}
// 	}
// 	return dst
// }

// func sortedContains(bs []byte, bss [][]byte) bool {
// 	idx := sort.Search(len(bss), func(n int) bool {
// 		return bytes.Compare(bs, bss[n]) <= 0
// 	})
// 	return idx < len(bss) && bytes.Equal(bs, bss[idx])
// }
