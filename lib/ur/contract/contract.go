// Copyright (C) 2020 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package contract

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/syncthing/syncthing/lib/structutil"
)

type Report struct {
	// Generated
	Received time.Time `json:"-"` // Only from DB
	Date     string    `json:"date,omitempty"`
	Address  string    `json:"address,omitempty"`

	// v1 fields

	UniqueID       string  `json:"uniqueID,omitempty" metric:"-" since:"1"`
	Version        string  `json:"version,omitempty" metric:"reports_total,gaugeVec:version" since:"1"`
	LongVersion    string  `json:"longVersion,omitempty" metric:"-" since:"1"`
	Platform       string  `json:"platform,omitempty" metric:"-" since:"1"`
	NumFolders     int     `json:"numFolders,omitempty" metric:"num_folders,summary" since:"1"`
	NumDevices     int     `json:"numDevices,omitempty" metric:"num_devices,summary" since:"1"`
	TotFiles       int     `json:"totFiles,omitempty" metric:"total_files,summary" since:"1"`
	FolderMaxFiles int     `json:"folderMaxFiles,omitempty" metric:"folder_max_files,summary" since:"1"`
	TotMiB         int     `json:"totMiB,omitempty" metric:"total_data_mib,summary" since:"1"`
	FolderMaxMiB   int     `json:"folderMaxMiB,omitempty" metric:"folder_max_data_mib,summary" since:"1"`
	MemoryUsageMiB int     `json:"memoryUsageMiB,omitempty" metric:"memory_usage_mib,summary" since:"1"`
	SHA256Perf     float64 `json:"sha256Perf,omitempty" metric:"sha256_perf_mibps,summary" since:"1"`
	HashPerf       float64 `json:"hashPerf,omitempty" metric:"hash_perf_mibps,summary" since:"1"`
	MemorySize     int     `json:"memorySize,omitempty" metric:"memory_size_mib,summary" since:"1"`

	// v2 fields

	URVersion  int `json:"urVersion,omitempty" metric:"reports_total,gaugeVec:ur_version" since:"2"`
	NumCPU     int `json:"numCPU,omitempty" metric:"num_cpu,summary" since:"2"`
	FolderUses struct {
		SendOnly            int `json:"sendonly,omitempty" metric:"folder_mode{mode=sendonly},summary" since:"2"`
		SendReceive         int `json:"sendreceive,omitempty" metric:"folder_mode{mode=sendreceive},summary" since:"2"`
		ReceiveOnly         int `json:"receiveonly,omitempty" metric:"folder_mode{mode=receiveonly},summary" since:"2"`
		IgnorePerms         int `json:"ignorePerms,omitempty" metric:"folder_ignore_perms,summary" since:"2"`
		IgnoreDelete        int `json:"ignoreDelete,omitempty" metric:"folder_ignore_delete,summary" since:"2"`
		AutoNormalize       int `json:"autoNormalize,omitempty" metric:"folder_auto_normalize,summary" since:"2"`
		SimpleVersioning    int `json:"simpleVersioning,omitempty" metric:"folder_versioning{kind=simple},summary" since:"2"`
		ExternalVersioning  int `json:"externalVersioning,omitempty" metric:"folder_versioning{kind=external},summary" since:"2"`
		StaggeredVersioning int `json:"staggeredVersioning,omitempty" metric:"folder_versioning{kind=staggered},summary" since:"2"`
		TrashcanVersioning  int `json:"trashcanVersioning,omitempty" metric:"folder_versioning{kind=trashcan},summary" since:"2"`
	} `json:"folderUses,omitempty" since:"2"`

	DeviceUses struct {
		Introducer       int `json:"introducer,omitempty" metric:"device_introducer,summary" since:"2"`
		CustomCertName   int `json:"customCertName,omitempty" metric:"device_custom_cert_name,summary" since:"2"`
		CompressAlways   int `json:"compressAlways,omitempty" metric:"device_compress{when=always},summary" since:"2"`
		CompressMetadata int `json:"compressMetadata,omitempty" metric:"device_compress{when=metadata},summary" since:"2"`
		CompressNever    int `json:"compressNever,omitempty" metric:"device_compress{when=never},summary" since:"2"`
		DynamicAddr      int `json:"dynamicAddr,omitempty" metric:"device_address{kind=dynamic},summary" since:"2"`
		StaticAddr       int `json:"staticAddr,omitempty" metric:"device_address{kind=static},summary" since:"2"`
	} `json:"deviceUses,omitempty" since:"2"`

	Announce struct {
		GlobalEnabled     bool `json:"globalEnabled,omitempty" metric:"discovery_enabled{kind=global},gauge" since:"2"`
		LocalEnabled      bool `json:"localEnabled,omitempty" metric:"discovery_enabled{kind=local},gauge" since:"2"`
		DefaultServersDNS int  `json:"defaultServersDNS,omitempty" metric:"discovery_default_servers_dns,summary" since:"2"`
		DefaultServersIP  int  `json:"defaultServersIP,omitempty" metric:"-" since:"2"` // Deprecated and not provided client-side anymore
		OtherServers      int  `json:"otherServers,omitempty" metric:"discovery_other_servers,summary" since:"2"`
	} `json:"announce,omitempty" since:"2"`

	Relays struct {
		Enabled        bool `json:"enabled,omitempty" metric:"relay_enabled,gauge" since:"2"`
		DefaultServers int  `json:"defaultServers,omitempty" metric:"relay_default_servers,summary" since:"2"`
		OtherServers   int  `json:"otherServers,omitempty" metric:"relay_other_servers,summary" since:"2"`
	} `json:"relays,omitempty" since:"2"`

	UsesRateLimit        bool `json:"usesRateLimit,omitempty" metric:"rate_limit_enabled,gauge" since:"2"`
	UpgradeAllowedManual bool `json:"upgradeAllowedManual,omitempty" metric:"upgrade_allowed_manual,gauge" since:"2"`
	UpgradeAllowedAuto   bool `json:"upgradeAllowedAuto,omitempty" metric:"upgrade_allowed_auto,gauge" since:"2"`

	// V2.5 fields (fields that were in v2 but never added to the database
	UpgradeAllowedPre bool  `json:"upgradeAllowedPre,omitempty" metric:"upgrade_allowed_pre,gauge" since:"2"`
	RescanIntvs       []int `json:"rescanIntvs,omitempty" metric:"rescan_intervals,summary" since:"2"`

	// v3 fields

	Uptime                     int    `json:"uptime,omitempty" metric:"uptime_seconds,summary" since:"3"`
	NATType                    string `json:"natType,omitempty" metric:"nat_detection,gaugeVec:type" since:"3"`
	AlwaysLocalNets            bool   `json:"alwaysLocalNets,omitempty" metric:"always_local_nets,gauge" since:"3"`
	CacheIgnoredFiles          bool   `json:"cacheIgnoredFiles,omitempty" metric:"cache_ignored_files,gauge" since:"3"`
	OverwriteRemoteDeviceNames bool   `json:"overwriteRemoteDeviceNames,omitempty" metric:"overwrite_remote_device_names,gauge" since:"3"`
	ProgressEmitterEnabled     bool   `json:"progressEmitterEnabled,omitempty" metric:"progress_emitter_enabled,gauge" since:"3"`
	CustomDefaultFolderPath    bool   `json:"customDefaultFolderPath,omitempty" metric:"custom_default_folder_path,gauge" since:"3"`
	WeakHashSelection          string `json:"weakHashSelection,omitempty" metric:"-" since:"3"` // Deprecated and not provided client-side anymore
	CustomTrafficClass         bool   `json:"customTrafficClass,omitempty" metric:"custom_traffic_class,gauge" since:"3"`
	CustomTempIndexMinBlocks   bool   `json:"customTempIndexMinBlocks,omitempty" metric:"custom_temp_index_min_blocks,gauge" since:"3"`
	TemporariesDisabled        bool   `json:"temporariesDisabled,omitempty" metric:"temporaries_disabled,gauge" since:"3"`
	TemporariesCustom          bool   `json:"temporariesCustom,omitempty" metric:"temporaries_custom,gauge" since:"3"`
	LimitBandwidthInLan        bool   `json:"limitBandwidthInLan,omitempty" metric:"limit_bandwidth_in_lan,gauge" since:"3"`
	CustomReleaseURL           bool   `json:"customReleaseURL,omitempty" metric:"custom_release_url,gauge" since:"3"`
	RestartOnWakeup            bool   `json:"restartOnWakeup,omitempty" metric:"restart_on_wakeup,gauge" since:"3"`
	CustomStunServers          bool   `json:"customStunServers,omitempty" metric:"custom_stun_servers,gauge" since:"3"`

	FolderUsesV3 struct {
		ScanProgressDisabled    int            `json:"scanProgressDisabled,omitempty" metric:"folder_scan_progress_disabled,summary" since:"3"`
		ConflictsDisabled       int            `json:"conflictsDisabled,omitempty" metric:"folder_conflicts_disabled,summary" since:"3"`
		ConflictsUnlimited      int            `json:"conflictsUnlimited,omitempty" metric:"folder_conflicts_unlimited,summary" since:"3"`
		ConflictsOther          int            `json:"conflictsOther,omitempty" metric:"folder_conflicts_other,summary" since:"3"`
		DisableSparseFiles      int            `json:"disableSparseFiles,omitempty" metric:"folder_disable_sparse_files,summary" since:"3"`
		DisableTempIndexes      int            `json:"disableTempIndexes,omitempty" metric:"folder_disable_temp_indexes,summary" since:"3"`
		AlwaysWeakHash          int            `json:"alwaysWeakHash,omitempty" metric:"folder_always_weakhash,summary" since:"3"`
		CustomWeakHashThreshold int            `json:"customWeakHashThreshold,omitempty" metric:"folder_custom_weakhash_threshold,summary" since:"3"`
		FsWatcherEnabled        int            `json:"fsWatcherEnabled,omitempty" metric:"folder_fswatcher_enabled,summary" since:"3"`
		PullOrder               map[string]int `json:"pullOrder,omitempty" metric:"folder_pull_order,summaryVec:order" since:"3"`
		FilesystemType          map[string]int `json:"filesystemType,omitempty" metric:"folder_file_system_type,summaryVec:type" since:"3"`
		FsWatcherDelays         []int          `json:"fsWatcherDelays,omitempty" metric:"folder_fswatcher_delays,summary" since:"3"`
		CustomMarkerName        int            `json:"customMarkerName,omitempty" metric:"folder_custom_markername,summary" since:"3"`
		CopyOwnershipFromParent int            `json:"copyOwnershipFromParent,omitempty" metric:"folder_copy_parent_ownership,summary" since:"3"`
		ModTimeWindowS          []int          `json:"modTimeWindowS,omitempty" metric:"folder_modtime_window_s,summary" since:"3"`
		MaxConcurrentWrites     []int          `json:"maxConcurrentWrites,omitempty" metric:"folder_max_concurrent_writes,summary" since:"3"`
		DisableFsync            int            `json:"disableFsync,omitempty" metric:"folder_disable_fsync,summary" since:"3"`
		BlockPullOrder          map[string]int `json:"blockPullOrder,omitempty" metric:"folder_block_pull_order:summaryVec:order" since:"3"`
		CopyRangeMethod         map[string]int `json:"copyRangeMethod,omitempty" metric:"folder_copy_range_method:summaryVec:method" since:"3"`
		CaseSensitiveFS         int            `json:"caseSensitiveFS,omitempty" metric:"folder_case_sensitive_fs,summary" since:"3"`
		ReceiveEncrypted        int            `json:"receiveencrypted,omitempty" metric:"folder_receive_encrypted,summary" since:"3"`
	} `json:"folderUsesV3,omitempty" since:"3"`

	DeviceUsesV3 struct {
		Untrusted int `json:"untrusted,omitempty" metric:"device_untrusted,summary" since:"3"`
	} `json:"deviceUsesV3,omitempty" since:"3"`

	GUIStats struct {
		Enabled                   int            `json:"enabled,omitempty" metric:"gui_enabled,summary" since:"3"`
		UseTLS                    int            `json:"useTLS,omitempty" metric:"gui_use_tls,summary" since:"3"`
		UseAuth                   int            `json:"useAuth,omitempty" metric:"gui_use_auth,summary" since:"3"`
		InsecureAdminAccess       int            `json:"insecureAdminAccess,omitempty" metric:"gui_insecure_admin_access,summary" since:"3"`
		Debugging                 int            `json:"debugging,omitempty" metric:"gui_debugging,summary" since:"3"`
		InsecureSkipHostCheck     int            `json:"insecureSkipHostCheck,omitempty" metric:"gui_insecure_skip_host_check,summary" since:"3"`
		InsecureAllowFrameLoading int            `json:"insecureAllowFrameLoading,omitempty" metric:"gui_insecure_allow_frame_loading,summary" since:"3"`
		ListenLocal               int            `json:"listenLocal,omitempty" metric:"gui_listen_local,summary" since:"3"`
		ListenUnspecified         int            `json:"listenUnspecified,omitempty" metric:"gui_listen_unspecified,summary" since:"3"`
		Theme                     map[string]int `json:"theme,omitempty" metric:"gui_theme,summaryVec:theme" since:"3"`
	} `json:"guiStats,omitempty" since:"3"`

	BlockStats struct {
		Total             int `json:"total,omitempty" metric:"blocks_processed_total,gauge" since:"3"`
		Renamed           int `json:"renamed,omitempty" metric:"blocks_processed{source=renamed},gauge" since:"3"`
		Reused            int `json:"reused,omitempty" metric:"blocks_processed{source=reused},gauge" since:"3"`
		Pulled            int `json:"pulled,omitempty" metric:"blocks_processed{source=pulled},gauge" since:"3"`
		CopyOrigin        int `json:"copyOrigin,omitempty" metric:"blocks_processed{source=copy_origin},gauge" since:"3"`
		CopyOriginShifted int `json:"copyOriginShifted,omitempty" metric:"blocks_processed{source=copy_origin_shifted},gauge" since:"3"`
		CopyElsewhere     int `json:"copyElsewhere,omitempty" metric:"blocks_processed{source=copy_elsewhere},gauge" since:"3"`
	} `json:"blockStats,omitempty" since:"3"`

	TransportStats map[string]int `json:"transportStats,omitempty" since:"3"`

	IgnoreStats struct {
		Lines           int `json:"lines,omitempty" metric:"folder_ignore_lines_total,summary" since:"3"`
		Inverts         int `json:"inverts,omitempty" metric:"folder_ignore_lines{kind=inverts},summary" since:"3"`
		Folded          int `json:"folded,omitempty" metric:"folder_ignore_lines{kind=folded},summary" since:"3"`
		Deletable       int `json:"deletable,omitempty" metric:"folder_ignore_lines{kind=deletable},summary" since:"3"`
		Rooted          int `json:"rooted,omitempty" metric:"folder_ignore_lines{kind=rooted},summary" since:"3"`
		Includes        int `json:"includes,omitempty" metric:"folder_ignore_lines{kind=includes},summary" since:"3"`
		EscapedIncludes int `json:"escapedIncludes,omitempty" metric:"folder_ignore_lines{kind=escapedIncludes},summary" since:"3"`
		DoubleStars     int `json:"doubleStars,omitempty" metric:"folder_ignore_lines{kind=doubleStars},summary" since:"3"`
		Stars           int `json:"stars,omitempty" metric:"folder_ignore_lines{kind=stars},summary" since:"3"`
	} `json:"ignoreStats,omitempty" since:"3"`

	// V3 fields added late in the RC
	WeakHashEnabled bool `json:"weakHashEnabled,omitempty" metric:"-" since:"3"` // Deprecated and not provided client-side anymore

	// Added in post processing
	OS           string `json:"-" metric:"reports_total,gaugeVec:os"`
	Arch         string `json:"-" metric:"reports_total,gaugeVec:arch"`
	Compiler     string `json:"-" metric:"compiler,gaugeVec:compiler"`
	Builder      string `json:"-" metric:"compiler,gaugeVec:builder"`
	Distribution string `json:"-" metric:"reports_total,gaugeVec:distribution"`
	City         string `json:"-" metric:"location,gaugeVec:city"`
	Country      string `json:"-" metric:"location,gaugeVec:country"`
}

func New() *Report {
	r := &Report{}
	structutil.FillNil(r)
	return r
}

func (r *Report) Validate() error {
	if r.UniqueID == "" || r.Version == "" || r.Platform == "" {
		return errors.New("missing required field")
	}
	if len(r.Date) != 8 {
		return errors.New("date not initialized")
	}

	// Some fields may not be null.
	if r.RescanIntvs == nil {
		r.RescanIntvs = []int{}
	}
	if r.FolderUsesV3.FsWatcherDelays == nil {
		r.FolderUsesV3.FsWatcherDelays = []int{}
	}

	return nil
}

func (r *Report) ClearForVersion(version int) error {
	return clear(r, version)
}

func (r *Report) FieldPointers() []interface{} {
	// All the fields of the Report, in the same order as the database fields.
	return []interface{}{
		&r.Received, &r.UniqueID, &r.Version, &r.LongVersion, &r.Platform,
		&r.NumFolders, &r.NumDevices, &r.TotFiles, &r.FolderMaxFiles,
		&r.TotMiB, &r.FolderMaxMiB, &r.MemoryUsageMiB, &r.SHA256Perf,
		&r.MemorySize, &r.Date,
		// V2
		&r.URVersion, &r.NumCPU, &r.FolderUses.SendOnly, &r.FolderUses.IgnorePerms,
		&r.FolderUses.IgnoreDelete, &r.FolderUses.AutoNormalize, &r.DeviceUses.Introducer,
		&r.DeviceUses.CustomCertName, &r.DeviceUses.CompressAlways,
		&r.DeviceUses.CompressMetadata, &r.DeviceUses.CompressNever,
		&r.DeviceUses.DynamicAddr, &r.DeviceUses.StaticAddr,
		&r.Announce.GlobalEnabled, &r.Announce.LocalEnabled,
		&r.Announce.DefaultServersDNS, &r.Announce.DefaultServersIP,
		&r.Announce.OtherServers, &r.Relays.Enabled, &r.Relays.DefaultServers,
		&r.Relays.OtherServers, &r.UsesRateLimit, &r.UpgradeAllowedManual,
		&r.UpgradeAllowedAuto, &r.FolderUses.SimpleVersioning,
		&r.FolderUses.ExternalVersioning, &r.FolderUses.StaggeredVersioning,
		&r.FolderUses.TrashcanVersioning,

		// V2.5
		&r.UpgradeAllowedPre,

		// V3
		&r.Uptime, &r.NATType, &r.AlwaysLocalNets, &r.CacheIgnoredFiles,
		&r.OverwriteRemoteDeviceNames, &r.ProgressEmitterEnabled, &r.CustomDefaultFolderPath,
		&r.WeakHashSelection, &r.CustomTrafficClass, &r.CustomTempIndexMinBlocks,
		&r.TemporariesDisabled, &r.TemporariesCustom, &r.LimitBandwidthInLan,
		&r.CustomReleaseURL, &r.RestartOnWakeup, &r.CustomStunServers,

		&r.FolderUsesV3.ScanProgressDisabled, &r.FolderUsesV3.ConflictsDisabled,
		&r.FolderUsesV3.ConflictsUnlimited, &r.FolderUsesV3.ConflictsOther,
		&r.FolderUsesV3.DisableSparseFiles, &r.FolderUsesV3.DisableTempIndexes,
		&r.FolderUsesV3.AlwaysWeakHash, &r.FolderUsesV3.CustomWeakHashThreshold,
		&r.FolderUsesV3.FsWatcherEnabled,

		&r.GUIStats.Enabled, &r.GUIStats.UseTLS, &r.GUIStats.UseAuth,
		&r.GUIStats.InsecureAdminAccess,
		&r.GUIStats.Debugging, &r.GUIStats.InsecureSkipHostCheck,
		&r.GUIStats.InsecureAllowFrameLoading, &r.GUIStats.ListenLocal,
		&r.GUIStats.ListenUnspecified,

		&r.BlockStats.Total, &r.BlockStats.Renamed,
		&r.BlockStats.Reused, &r.BlockStats.Pulled, &r.BlockStats.CopyOrigin,
		&r.BlockStats.CopyOriginShifted, &r.BlockStats.CopyElsewhere,

		&r.IgnoreStats.Lines, &r.IgnoreStats.Inverts, &r.IgnoreStats.Folded,
		&r.IgnoreStats.Deletable, &r.IgnoreStats.Rooted, &r.IgnoreStats.Includes,
		&r.IgnoreStats.EscapedIncludes, &r.IgnoreStats.DoubleStars, &r.IgnoreStats.Stars,

		// V3 added late in the RC
		&r.WeakHashEnabled,
		&r.Address,

		// Receive only folders
		&r.FolderUses.ReceiveOnly,
	}
}

func (*Report) FieldNames() []string {
	// The database fields that back this struct in PostgreSQL
	return []string{
		// V1
		"Received",
		"UniqueID",
		"Version",
		"LongVersion",
		"Platform",
		"NumFolders",
		"NumDevices",
		"TotFiles",
		"FolderMaxFiles",
		"TotMiB",
		"FolderMaxMiB",
		"MemoryUsageMiB",
		"SHA256Perf",
		"MemorySize",
		"Date",
		// V2
		"ReportVersion",
		"NumCPU",
		"FolderRO",
		"FolderIgnorePerms",
		"FolderIgnoreDelete",
		"FolderAutoNormalize",
		"DeviceIntroducer",
		"DeviceCustomCertName",
		"DeviceCompressionAlways",
		"DeviceCompressionMetadata",
		"DeviceCompressionNever",
		"DeviceDynamicAddr",
		"DeviceStaticAddr",
		"AnnounceGlobalEnabled",
		"AnnounceLocalEnabled",
		"AnnounceDefaultServersDNS",
		"AnnounceDefaultServersIP",
		"AnnounceOtherServers",
		"RelayEnabled",
		"RelayDefaultServers",
		"RelayOtherServers",
		"RateLimitEnabled",
		"UpgradeAllowedManual",
		"UpgradeAllowedAuto",
		// v0.12.19+
		"FolderSimpleVersioning",
		"FolderExternalVersioning",
		"FolderStaggeredVersioning",
		"FolderTrashcanVersioning",
		// V2.5
		"UpgradeAllowedPre",
		// V3
		"Uptime",
		"NATType",
		"AlwaysLocalNets",
		"CacheIgnoredFiles",
		"OverwriteRemoteDeviceNames",
		"ProgressEmitterEnabled",
		"CustomDefaultFolderPath",
		"WeakHashSelection",
		"CustomTrafficClass",
		"CustomTempIndexMinBlocks",
		"TemporariesDisabled",
		"TemporariesCustom",
		"LimitBandwidthInLan",
		"CustomReleaseURL",
		"RestartOnWakeup",
		"CustomStunServers",

		"FolderScanProgressDisabled",
		"FolderConflictsDisabled",
		"FolderConflictsUnlimited",
		"FolderConflictsOther",
		"FolderDisableSparseFiles",
		"FolderDisableTempIndexes",
		"FolderAlwaysWeakHash",
		"FolderCustomWeakHashThreshold",
		"FolderFsWatcherEnabled",

		"GUIEnabled",
		"GUIUseTLS",
		"GUIUseAuth",
		"GUIInsecureAdminAccess",
		"GUIDebugging",
		"GUIInsecureSkipHostCheck",
		"GUIInsecureAllowFrameLoading",
		"GUIListenLocal",
		"GUIListenUnspecified",

		"BlocksTotal",
		"BlocksRenamed",
		"BlocksReused",
		"BlocksPulled",
		"BlocksCopyOrigin",
		"BlocksCopyOriginShifted",
		"BlocksCopyElsewhere",

		"IgnoreLines",
		"IgnoreInverts",
		"IgnoreFolded",
		"IgnoreDeletable",
		"IgnoreRooted",
		"IgnoreIncludes",
		"IgnoreEscapedIncludes",
		"IgnoreDoubleStars",
		"IgnoreStars",

		// V3 added late in the RC
		"WeakHashEnabled",
		"Address",

		// Receive only folders
		"FolderRecvOnly",
	}
}

func (r Report) Value() (driver.Value, error) {
	// This needs to be string, yet we read back bytes..
	bs, err := json.Marshal(r)
	return string(bs), err
}

func (r *Report) Scan(value interface{}) error {
	// Zero out the previous value
	// JSON un-marshaller does not touch fields that are not in the payload, so we carry over values from a previous
	// scan.
	*r = Report{}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &r)
}

func clear(v interface{}, since int) error {
	s := reflect.ValueOf(v).Elem()
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		tag := t.Field(i).Tag

		v := tag.Get("since")
		if v == "" {
			f.Set(reflect.Zero(f.Type()))
			continue
		}

		vn, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		if vn > since {
			f.Set(reflect.Zero(f.Type()))
			continue
		}

		// Dive deeper
		if f.Kind() == reflect.Ptr {
			f = f.Elem()
		}

		if f.Kind() == reflect.Struct {
			if err := clear(f.Addr().Interface(), since); err != nil {
				return err
			}
		}
	}
	return nil
}
