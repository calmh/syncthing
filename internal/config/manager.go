// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/syncthing/syncthing/internal/slogutil"
	"github.com/syncthing/syncthing/lib/events"
	"github.com/syncthing/syncthing/lib/osutil"
	"github.com/syncthing/syncthing/lib/protocol"
	"github.com/syncthing/syncthing/lib/sliceutil"
)

const (
	maxModifications = 1000
	minSaveInterval  = 5 * time.Second
)

var errTooManyModifications = errors.New("too many concurrent config modifications")

// The Committer and Verifier interfaces are implemented by objects
// that need to know about or have a say in configuration changes.
//
// When the configuration is about to be changed, VerifyConfiguration() is
// called for each subscribing object that implements it, with copies of the
// old and new configuration. A nil error is returned if the new configuration
// is acceptable (i.e., does not contain any errors that would prevent it from
// being a valid config). Otherwise an error describing the problem is returned.
//
// If any subscriber returns an error from VerifyConfiguration(), the
// configuration change is not committed and an error is returned to whoever
// tried to commit the broken config.
//
// If all verification calls returns nil, CommitConfiguration() is called for
// each subscribing object. The callee returns true if the new configuration
// has been successfully applied, otherwise false. Any Commit() call returning
// false will result in a "restart needed" response to the API/user. Note that
// the new configuration will still have been applied by those who were
// capable of doing so.
//
// A Committer must take care not to hold any locks while changing the
// configuration (e.g. calling Manager.SetFolder), that are also acquired in any
// methods of the Committer interface.
type Committer interface {
	CommitConfiguration(from, to Configuration) (handled bool)
	String() string
}

// A Verifier can determine if a new configuration is acceptable.
// See the description for Committer, above.
type Verifier interface {
	VerifyConfiguration(from, to Configuration) error
}

// Waiter allows to wait for the given config operation to complete.
type Waiter interface {
	Wait()
}

type noopWaiter struct{}

func (noopWaiter) Wait() {}

// ModifyFunction gets a pointer to a copy of the currently active configuration
// for modification.
type ModifyFunction func(*Configuration)

type Manager struct {
	cfg      Configuration
	path     string
	evLogger events.Logger
	myID     protocol.DeviceID
	queue    chan modifyEntry

	waiter Waiter // Latest ongoing config change
	subs   []Committer
	mut    sync.Mutex

	requiresRestart atomic.Bool
}

// Manage wraps an existing Configuration structure and ties it to a file on
// disk.
// The returned Manager is a suture.Service, thus needs to be started (added to
// a supervisor).
func Manage(path string, cfg Configuration, myID protocol.DeviceID, evLogger events.Logger) *Manager {
	w := &Manager{
		cfg:      cfg,
		path:     path,
		evLogger: evLogger,
		myID:     myID,
		queue:    make(chan modifyEntry, maxModifications),
		waiter:   noopWaiter{}, // Noop until first config change
	}
	return w
}

// Load loads an existing file on disk and returns a new configuration
// manager. The file must contain a configuration in YAML format; it is
// validated before the manager is created.
// The returned Manager is a suture.Service, thus needs to be started (added to
// a supervisor).
func Load(path string, myID protocol.DeviceID, evLogger events.Logger) (*Manager, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg, err := Unmarshal(bs)
	if err != nil {
		return nil, err
	}
	return Manage(path, cfg, myID, evLogger), nil
}

func (w *Manager) ConfigPath() string {
	return w.path
}

func (w *Manager) MyID() protocol.DeviceID {
	return w.myID
}

// Subscribe registers the given handler to be called on any future
// configuration changes. It returns the config that is in effect while
// subscribing, which can be used for initial setup.
func (w *Manager) Subscribe(c Committer) Configuration {
	w.mut.Lock()
	defer w.mut.Unlock()
	w.subs = append(w.subs, c)
	return w.cfg.Copy()
}

// Unsubscribe de-registers the given handler from any future calls to
// configuration changes and only returns after a potential ongoing config
// change is done.
func (w *Manager) Unsubscribe(c Committer) {
	w.mut.Lock()
	for i := range w.subs {
		if w.subs[i] == c {
			w.subs = sliceutil.RemoveAndZero(w.subs, i)
			break
		}
	}
	waiter := w.waiter
	w.mut.Unlock()
	// Waiting mustn't be done under lock, as the goroutines in notifyListener
	// may dead-lock when trying to access lock on config read operations.
	waiter.Wait()
}

// RawCopy returns a copy of the currently managed Configuration object.
func (w *Manager) RawCopy() Configuration {
	w.mut.Lock()
	defer w.mut.Unlock()
	return w.cfg.Copy()
}

func (w *Manager) Modify(fn ModifyFunction) (Waiter, error) {
	return w.modifyQueued(fn)
}

func (w *Manager) modifyQueued(modifyFunc ModifyFunction) (Waiter, error) {
	e := modifyEntry{
		modifyFunc: modifyFunc,
		res:        make(chan modifyResult),
	}
	select {
	case w.queue <- e:
	default:
		return noopWaiter{}, errTooManyModifications
	}
	res := <-e.res
	return res.w, res.err
}

func (w *Manager) Serve(ctx context.Context) error {
	defer w.serveSave()

	var e modifyEntry
	saveTimer := time.NewTimer(0)
	<-saveTimer.C
	saveTimerRunning := false
	for {
		select {
		case e = <-w.queue:
		case <-saveTimer.C:
			w.serveSave()
			saveTimerRunning = false
			continue
		case <-ctx.Done():
			return ctx.Err()
		}

		var waiter Waiter = noopWaiter{}
		var err error

		// Let the caller modify the config.
		to := w.RawCopy()
		e.modifyFunc(&to)

		// Check if the config was actually changed at all.
		w.mut.Lock()
		if !proto.Equal(w.cfg.Configuration, to.Configuration) {
			waiter, err = w.replaceLocked(to)
			if !saveTimerRunning {
				saveTimer.Reset(minSaveInterval)
				saveTimerRunning = true
			}
		}
		w.mut.Unlock()

		e.res <- modifyResult{
			w:   waiter,
			err: err,
		}

		// Wait for all subscriber to handle the config change before continuing
		// to process the next change.
		done := make(chan struct{})
		go func() {
			waiter.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (w *Manager) serveSave() {
	if w.path == "" {
		return
	}
	if err := w.Save(); err != nil {
		slog.Warn("Failed to save config", slogutil.Error(err))
	}
}

func (w *Manager) replaceLocked(to Configuration) (Waiter, error) {
	from := w.cfg

	if err := Validate(to); err != nil {
		return noopWaiter{}, err
	}

	for _, sub := range w.subs {
		sub, ok := sub.(Verifier)
		if !ok {
			continue
		}
		if err := sub.VerifyConfiguration(from.Copy(), to.Copy()); err != nil {
			return noopWaiter{}, err
		}
	}

	w.cfg = to

	w.waiter = w.notifyListeners(from.Copy(), to.Copy())

	return w.waiter, nil
}

func (w *Manager) notifyListeners(from, to Configuration) Waiter {
	wg := new(sync.WaitGroup)
	for _, sub := range w.subs {
		wg.Go(func() {
			w.notifyListener(sub, from, to)
		})
	}
	return wg
}

func (w *Manager) notifyListener(sub Committer, from, to Configuration) {
	if !sub.CommitConfiguration(from, to) {
		w.requiresRestart.Store(true)
	}
}

// Devices returns a map of devices.
func (w *Manager) Devices() map[protocol.DeviceID]DeviceConfiguration {
	w.mut.Lock()
	defer w.mut.Unlock()
	devices := w.cfg.GetDevices()
	deviceMap := make(map[protocol.DeviceID]DeviceConfiguration, len(devices))
	for _, dev := range devices {
		deviceMap[dev.GetDeviceId()] = dev.Copy()
	}
	return deviceMap
}

// DeviceList returns a slice of devices.
func (w *Manager) DeviceList() []DeviceConfiguration {
	w.mut.Lock()
	defer w.mut.Unlock()
	return w.cfg.Copy().GetDevices()
}

// RemoveDevice removes the device from the configuration
func (w *Manager) RemoveDevice(id protocol.DeviceID) (Waiter, error) {
	return w.modifyQueued(func(cfg *Configuration) {
		if _, i, ok := cfg.Device(id); ok {
			devices := cfg.Configuration.GetDevices()
			cfg.Configuration.SetDevices(slices.Delete(devices, i, i+1))
		}
	})
}

func (w *Manager) DefaultDevice() DeviceConfiguration {
	w.mut.Lock()
	defer w.mut.Unlock()
	return w.cfg.GetDefaults().GetDevice().Copy()
}

// Folders returns a map of folders.
func (w *Manager) Folders() map[string]FolderConfiguration {
	w.mut.Lock()
	defer w.mut.Unlock()
	folders := w.cfg.GetFolders()
	folderMap := make(map[string]FolderConfiguration, len(folders))
	for _, fld := range folders {
		folderMap[fld.GetId()] = fld.Copy()
	}
	return folderMap
}

// FolderList returns a slice of folders.
func (w *Manager) FolderList() []FolderConfiguration {
	w.mut.Lock()
	defer w.mut.Unlock()
	return w.cfg.Copy().GetFolders()
}

// RemoveFolder removes the folder from the configuration
func (w *Manager) RemoveFolder(id string) (Waiter, error) {
	return w.modifyQueued(func(cfg *Configuration) {
		if _, i, ok := cfg.Folder(id); ok {
			folders := cfg.Configuration.GetFolders()
			cfg.Configuration.SetFolders(slices.Delete(folders, i, i+1))
		}
	})
}

// FolderPasswords returns the folder passwords set for this device, for
// folders that have an encryption password set.
func (w *Manager) FolderPasswords(device protocol.DeviceID) map[string]string {
	w.mut.Lock()
	defer w.mut.Unlock()
	res := make(map[string]string)
	for _, folder := range w.cfg.GetFolders() {
		if dev, ok := folder.Device(device); ok {
			if pw := dev.GetEncryptionPassword(); pw != "" {
				res[folder.GetId()] = pw
			}
		}
	}
	return res
}

func (w *Manager) DefaultFolder() FolderConfiguration {
	w.mut.Lock()
	defer w.mut.Unlock()
	return w.cfg.GetDefaults().GetFolder().Copy()
}

// Options returns the current options configuration object, or nil if
// unset. Getters on the result are nil-safe and return the schema
// defaults.
func (w *Manager) Options() *OptionsConfiguration {
	w.mut.Lock()
	defer w.mut.Unlock()
	if opts := w.cfg.GetOptions(); opts != nil {
		return proto.Clone(opts).(*OptionsConfiguration)
	}
	return nil
}

func (w *Manager) LDAP() *LDAPConfiguration {
	w.mut.Lock()
	defer w.mut.Unlock()
	if ldap := w.cfg.GetLdap(); ldap != nil {
		return proto.Clone(ldap).(*LDAPConfiguration)
	}
	return nil
}

// GUI returns the current GUI configuration object, or nil if unset.
// Getters on the result are nil-safe and return the schema defaults.
func (w *Manager) GUI() *GUIConfiguration {
	w.mut.Lock()
	defer w.mut.Unlock()
	if gui := w.cfg.GetGui(); gui != nil {
		return proto.Clone(gui).(*GUIConfiguration)
	}
	return nil
}

// DefaultIgnores returns the list of ignore patterns to be used by default on
// folders, or nil if unset.
func (w *Manager) DefaultIgnores() *Ignores {
	w.mut.Lock()
	defer w.mut.Unlock()
	if ignores := w.cfg.GetDefaults().GetIgnores(); ignores != nil {
		return proto.Clone(ignores).(*Ignores)
	}
	return nil
}

// IgnoredDevice returns whether or not connection attempts from the given
// device should be silently ignored.
func (w *Manager) IgnoredDevice(id protocol.DeviceID) bool {
	w.mut.Lock()
	defer w.mut.Unlock()
	for _, device := range w.cfg.GetRemoteIgnoredDevices() {
		if device.GetDeviceId() == id {
			return true
		}
	}
	return false
}

// IgnoredDevices returns a slice of ignored devices.
func (w *Manager) IgnoredDevices() []ObservedDevice {
	w.mut.Lock()
	defer w.mut.Unlock()
	devices := w.cfg.GetRemoteIgnoredDevices()
	res := make([]ObservedDevice, len(devices))
	for i, device := range devices {
		res[i] = device.Copy()
	}
	return res
}

// IgnoredFolder returns whether or not share attempts for the given
// folder should be silently ignored.
func (w *Manager) IgnoredFolder(device protocol.DeviceID, folder string) bool {
	dev, ok := w.Device(device)
	if !ok {
		return false
	}
	for _, f := range dev.GetIgnoredFolders() {
		if f.GetId() == folder {
			return true
		}
	}
	return false
}

// Device returns the configuration for the given device and an "ok" bool.
func (w *Manager) Device(id protocol.DeviceID) (DeviceConfiguration, bool) {
	w.mut.Lock()
	defer w.mut.Unlock()
	device, _, ok := w.cfg.Device(id)
	if !ok {
		return DeviceConfiguration{}, false
	}
	return device.Copy(), ok
}

// Folder returns the configuration for the given folder and an "ok" bool.
func (w *Manager) Folder(id string) (FolderConfiguration, bool) {
	w.mut.Lock()
	defer w.mut.Unlock()
	fcfg, _, ok := w.cfg.Folder(id)
	if !ok {
		return FolderConfiguration{}, false
	}
	return fcfg.Copy(), ok
}

// Save writes the configuration to disk, and generates a ConfigSaved event.
func (w *Manager) Save() error {
	w.mut.Lock()
	defer w.mut.Unlock()

	fd, err := osutil.CreateAtomic(w.path)
	if err != nil {
		return err
	}

	bs, err := Marshal(w.cfg)
	if err != nil {
		fd.Close()
		return err
	}
	if _, err := fd.Write(bs); err != nil {
		fd.Close()
		return err
	}
	if err := fd.Close(); err != nil {
		return err
	}

	// The event carries the configuration in JSON form, so that it can be
	// embedded as such in the event stream.
	if data, err := protojson.Marshal(w.cfg.Configuration); err == nil {
		w.evLogger.Log(events.ConfigSaved, json.RawMessage(data))
	}
	return nil
}

func (w *Manager) RequiresRestart() bool { return w.requiresRestart.Load() }

type modifyEntry struct {
	modifyFunc ModifyFunction
	res        chan modifyResult
}

type modifyResult struct {
	w   Waiter
	err error
}
