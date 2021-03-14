package model

import (
	"bytes"
	"sort"

	"github.com/syncthing/syncthing/lib/db"
	"github.com/syncthing/syncthing/lib/protocol"
)

const dbKey = "allFolderDevices"

type allFolderDevicesStore struct {
	kv  *db.NamespacedKV
	cur SharingCandidates
}

func newSharingCandidatesTracker(kv *db.NamespacedKV) *allFolderDevicesStore {
	s := &allFolderDevicesStore{kv: kv}
	_ = s.load()
	return s
}

func (s *allFolderDevicesStore) AddFolderDevices(cc *protocol.ClusterConfig) error {
	for _, ccFolder := range cc.Folders {
		candFolder := s.cur.Folders[ccFolder.ID]
		ccDevices := devicesToBytes(ccFolder.Devices)
		candFolder.Devices = mergeLists(candFolder.Devices, ccDevices)
	}
	return s.save()
}

func (s *allFolderDevicesStore) DevicesForFolder(folderID string) ([]protocol.DeviceID, error) {
	devs := s.cur.Folders[folderID].Devices
	ids := make([]protocol.DeviceID, 0, len(devs))
	for i := range devs {
		id, err := protocol.DeviceIDFromBytes(devs[i])
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *allFolderDevicesStore) save() error {
	bs, err := s.cur.Marshal()
	if err != nil {
		return err
	}
	return s.kv.PutBytes(dbKey, bs)
}

func (s *allFolderDevicesStore) load() error {
	bs, ok, err := s.kv.Bytes(dbKey)
	if err != nil {
		return err
	}
	var cur SharingCandidates
	if ok {
		if err := cur.Unmarshal(bs); err != nil {
			return err
		}
	}
	s.cur = cur
	return nil
}

// devicesToBytes returns a list of device IDs in byte slice form, given a
// list of devices in a cluster config. The returned list is sorted.
func devicesToBytes(devs []protocol.Device) [][]byte {
	res := make([][]byte, len(devs))
	for i := range devs {
		res[i] = devs[i].ID[:]
	}
	sort.Slice(res, func(a, b int) bool {
		return bytes.Compare(res[a], res[b]) == -1
	})
	return res
}

// mergeLists merges src into dst without duplicates, assuming both lists
// are already sorted.
func mergeLists(dst, src [][]byte) [][]byte {
	di, si := 0, 0
	for {
		if si == len(src) {
			break
		}
		if di == len(dst) {
			// We've reached the end of dst, merge becomes append
			dst = append(dst, src[si])
			si++
			continue
		}

		switch bytes.Compare(src[si], dst[di]) {
		case -1:
			// src[si] should come before dst[di], so we insert here
			dst = append(dst[:di], append([][]byte{src[si]}, dst[di:]...)...)
			si++
		case 0:
			// src[si] is already present in dst, look for next src element
			si++
		case 1:
			// src[i] should come later, look further in dst
			di++
		}
	}
	return dst
}
