// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package syncthing

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	newdb "github.com/syncthing/syncthing/internal/db"
	"github.com/syncthing/syncthing/internal/db/dbext"
	"github.com/syncthing/syncthing/internal/db/sqlite"
	"github.com/syncthing/syncthing/lib/build"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/db"
	"github.com/syncthing/syncthing/lib/db/backend"
	"github.com/syncthing/syncthing/lib/events"
	"github.com/syncthing/syncthing/lib/fs"
	"github.com/syncthing/syncthing/lib/locations"
	"github.com/syncthing/syncthing/lib/protocol"
	"github.com/syncthing/syncthing/lib/tlsutil"
)

func EnsureDir(dir string, mode fs.FileMode) error {
	fs := fs.NewFilesystem(fs.FilesystemTypeBasic, dir)
	err := fs.MkdirAll(".", mode)
	if err != nil {
		return err
	}

	if fi, err := fs.Stat("."); err == nil {
		// Apparently the stat may fail even though the mkdirall passed. If it
		// does, we'll just assume things are in order and let other things
		// fail (like loading or creating the config...).
		currentMode := fi.Mode() & 0o777
		if currentMode != mode {
			err := fs.Chmod(".", mode)
			// This can fail on crappy filesystems, nothing we can do about it.
			if err != nil {
				l.Warnln(err)
			}
		}
	}
	return nil
}

func LoadOrGenerateCertificate(certFile, keyFile string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return GenerateCertificate(certFile, keyFile)
	}
	return cert, nil
}

func GenerateCertificate(certFile, keyFile string) (tls.Certificate, error) {
	l.Infof("Generating ECDSA key and certificate for %s...", tlsDefaultCommonName)
	return tlsutil.NewCertificate(certFile, keyFile, tlsDefaultCommonName, deviceCertLifetimeDays)
}

func DefaultConfig(path string, myID protocol.DeviceID, evLogger events.Logger, noDefaultFolder, skipPortProbing bool) (config.Wrapper, error) {
	newCfg := config.New(myID)

	if skipPortProbing {
		l.Infoln("Using default network port numbers instead of probing for free ports")
		// Record address override initially
		newCfg.GUI.RawAddress = newCfg.GUI.Address()
	} else if err := newCfg.ProbeFreePorts(); err != nil {
		return nil, err
	}

	if noDefaultFolder {
		l.Infoln("We will skip creation of a default folder on first start")
		return config.Wrap(path, newCfg, myID, evLogger), nil
	}

	fcfg := newCfg.Defaults.Folder.Copy()
	fcfg.ID = "default"
	fcfg.Label = "Default Folder"
	fcfg.FilesystemType = config.FilesystemTypeBasic
	fcfg.Path = locations.Get(locations.DefFolder)
	newCfg.Folders = append(newCfg.Folders, fcfg)
	l.Infoln("Default folder created and/or linked to new config")
	return config.Wrap(path, newCfg, myID, evLogger), nil
}

// LoadConfigAtStartup loads an existing config. If it doesn't yet exist, it
// creates a default one, without the default folder if noDefaultFolder is true.
// Otherwise it checks the version, and archives and upgrades the config if
// necessary or returns an error, if the version isn't compatible.
func LoadConfigAtStartup(path string, cert tls.Certificate, evLogger events.Logger, allowNewerConfig, noDefaultFolder, skipPortProbing bool) (config.Wrapper, error) {
	myID := protocol.NewDeviceID(cert.Certificate[0])
	cfg, originalVersion, err := config.Load(path, myID, evLogger)
	if fs.IsNotExist(err) {
		cfg, err = DefaultConfig(path, myID, evLogger, noDefaultFolder, skipPortProbing)
		if err != nil {
			return nil, fmt.Errorf("failed to generate default config: %w", err)
		}
		err = cfg.Save()
		if err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		l.Infof("Default config saved. Edit %s to taste (with Syncthing stopped) or use the GUI", cfg.ConfigPath())
	} else if err == io.EOF {
		return nil, errors.New("failed to load config: unexpected end of file. Truncated or empty configuration?")
	} else if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if originalVersion != config.CurrentVersion {
		if originalVersion == config.CurrentVersion+1101 {
			l.Infof("Now, THAT's what we call a config from the future! Don't worry. As long as you hit that wire with the connecting hook at precisely eighty-eight miles per hour the instant the lightning strikes the tower... everything will be fine.")
		}
		if originalVersion > config.CurrentVersion && !allowNewerConfig {
			return nil, fmt.Errorf("config file version (%d) is newer than supported version (%d). If this is expected, use --allow-newer-config to override.", originalVersion, config.CurrentVersion)
		}
		err = archiveAndSaveConfig(cfg, originalVersion)
		if err != nil {
			return nil, fmt.Errorf("config archive: %w", err)
		}
	}

	return cfg, nil
}

func archiveAndSaveConfig(cfg config.Wrapper, originalVersion int) error {
	// Copy the existing config to an archive copy
	archivePath := cfg.ConfigPath() + fmt.Sprintf(".v%d", originalVersion)
	l.Infoln("Archiving a copy of old config file format at:", archivePath)
	if err := copyFile(cfg.ConfigPath(), archivePath); err != nil {
		return err
	}

	// Do a regular atomic config sve
	return cfg.Save()
}

func copyFile(src, dst string) error {
	bs, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dst, bs, 0o600); err != nil {
		// Attempt to clean up
		os.Remove(dst)
		return err
	}

	return nil
}

func OpenDBBackend(path string, tuning config.Tuning) (backend.Backend, error) {
	return backend.Open(path, backend.Tuning(tuning))
}

// Opens a database and attempts migrating the legacy database to the new database format.
func OpenDatabase(path string, oldDBDir string, evLogger events.Logger) (newdb.DB, error) {
	sql, err := sqlite.Open(path)
	if err != nil {
		return nil, err
	}

	sdb := newdb.MetricsWrap(sql)

	if be, err := backend.OpenLevelDBRO(oldDBDir); err == nil {
		// We have not migrated. We should do that.
		err := migrateDatabase(be, evLogger, sdb, oldDBDir)
		if err != nil {
			l.Warnln(err.Error())
			os.Exit(0) // prevent automatic restart by the monitor
		}

		miscDB := dbext.NewMiscDB(sdb)
		_ = miscDB.PutTime("migrated-from-leveldb-at", time.Now())
		_ = miscDB.PutString("migrated-from-leveldb-by", build.LongVersion)
	}

	return sdb, nil
}

func migrateDatabase(be backend.Backend, evLogger events.Logger, sdb newdb.DB, oldDBDir string) error {
	l.Infoln("Migrating database from LevelDB to SQLite; this can take quite a while...")

	ll, err := db.NewLowlevel(be, evLogger)
	if err != nil {
		return errors.New("Failed to migrate: " + err.Error())
	}

	for _, folder := range ll.ListFolders() {
		l.Infoln("Migrating folder", folder, "...")
		var batch []protocol.FileInfo
		fs, err := db.NewFileSet(folder, ll)
		if err != nil {
			return errors.New("Failed to migrate FileInfos: " + err.Error())
		}

		snap, err := fs.Snapshot()
		if err != nil {
			return errors.New("Failed to migrate FileInfos: " + err.Error())
		}

		err = nil
		snap.WithHaveSequence(0, func(f protocol.FileInfo) bool {
			batch = append(batch, f)
			if len(batch) == 1000 {
				if err = sdb.Update(folder, protocol.LocalDeviceID, batch); err != nil {
					return false
				}
				batch = batch[:0]
			}
			return true
		})

		if err != nil {
			return errors.New("Failed to migrate FileInfos: " + err.Error())
		}

		if len(batch) > 0 {
			if err := sdb.Update(folder, protocol.LocalDeviceID, batch); err != nil {
				return errors.New("Failed to migrate FileInfos: " + err.Error())
			}
		}
		snap.Release()
	}

	l.Infoln("Migrating virtual mtimes...")
	if err := ll.IterateMtimes(sdb.MtimePut); err != nil {
		l.Warnln("Failed to migrate mtimes:", err)
	}

	l.Infoln("Migration complete")
	be.Close()

	if err := os.Rename(oldDBDir, oldDBDir+"-migrated"); err != nil {
		return errors.New("Failed to rename old, migrated database: " + err.Error() + ". Please manually move or remove " + oldDBDir)
	}

	return nil
}
