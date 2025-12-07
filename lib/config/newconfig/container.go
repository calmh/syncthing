package newconfig

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"slices"
	"strconv"
)

type simpleValue interface {
	// expanding the set of supported types requires implementing
	// marshal/unmarshal and corresponding tests
	~int | ~string | ~float64 | ~bool
}

// A Setting represents the default and current value of a given setting.
type Setting[T simpleValue] struct {
	value  T
	defval T
	valid  bool
}

// Get returns the current value or the default if no value is set.
func (c Setting[T]) Get() T {
	if c.valid {
		return c.value
	}
	return c.defval
}

// Set sets the current value.
func (c *Setting[T]) Set(v T) {
	c.value = v
	c.valid = true
}

// SetOrDefault sets the value if it is different from the default,
// otherwise marks the setting as using the default value.
func (c *Setting[T]) SetOrDefault(v T) {
	if v == c.defval {
		var zero T
		c.value = zero
		c.valid = false
		return
	}
	c.Set(v)
}

func (c *Setting[T]) ParseDefault(s string) error {
	return c.unmarshalValueInto(s, &c.defval)
}

func (c Setting[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var zero T
	if !c.valid {
		start.Attr = []xml.Attr{{Name: xml.Name{Local: "usesDefault"}, Value: "yes"}}
	} else if c.value == zero && len(fmt.Sprint(c.value)) == 0 {
		start.Attr = []xml.Attr{{Name: xml.Name{Local: "usesDefault"}, Value: "no"}}
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if c.valid {
		if err := e.EncodeToken(xml.CharData(fmt.Sprint(c.value))); err != nil {
			return err
		}
	} else {
		def := fmt.Sprint(c.defval)
		if def == "" {
			def = "(empty)"
		}
		if err := e.EncodeToken(xml.Comment(fmt.Sprintf(" default: %v ", def))); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}

func (c *Setting[T]) unmarshalValueInto(s string, ptr *T) error {
	v := reflect.ValueOf(ptr).Elem()
	switch v.Kind() {
	case reflect.Int:
		iv, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(iv)
	case reflect.String:
		v.SetString(s)
	case reflect.Float64:
		fv, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(fv)
	case reflect.Bool:
		bv, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(bv)
	default:
		return fmt.Errorf("unknown type %v", v.Kind())
	}
	return nil
}

func (c *Setting[T]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var zero T
	c.value = zero

	c.valid = !slices.ContainsFunc(start.Attr, func(a xml.Attr) bool {
		return a.Name.Local == "usesDefault" && a.Value == "yes"
	})

	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch tok := tok.(type) {
		case xml.EndElement:
			return nil
		case xml.CharData:
			if err := c.unmarshalValueInto(string(tok), &c.value); err != nil {
				return err
			}
		case xml.Comment:
			continue
		default:
			fmt.Printf("wat %+v", tok)
		}
	}
}

type OptionsConfiguration struct {
	RawListenAddresses      []string        `json:"listenAddresses" xml:"listenAddress" default:"default"`
	RawGlobalAnnServers     []string        `json:"globalAnnounceServers" xml:"globalAnnounceServer" default:"default"`
	GlobalAnnEnabled        Setting[bool]   `json:"globalAnnounceEnabled" xml:"globalAnnounceEnabled" default:"true"`
	LocalAnnEnabled         Setting[bool]   `json:"localAnnounceEnabled" xml:"localAnnounceEnabled" default:"true"`
	LocalAnnPort            Setting[int]    `json:"localAnnouncePort" xml:"localAnnouncePort" default:"21027"`
	LocalAnnMCAddr          Setting[string] `json:"localAnnounceMCAddr" xml:"localAnnounceMCAddr" default:"[ff12::8384]:21027"`
	MaxSendKbps             Setting[int]    `json:"maxSendKbps" xml:"maxSendKbps"`
	MaxRecvKbps             Setting[int]    `json:"maxRecvKbps" xml:"maxRecvKbps"`
	ReconnectIntervalS      Setting[int]    `json:"reconnectionIntervalS" xml:"reconnectionIntervalS" default:"60"`
	RelaysEnabled           Setting[bool]   `json:"relaysEnabled" xml:"relaysEnabled" default:"true"`
	RelayReconnectIntervalM Setting[int]    `json:"relayReconnectIntervalM" xml:"relayReconnectIntervalM" default:"10"`
	StartBrowser            Setting[bool]   `json:"startBrowser" xml:"startBrowser" default:"true"`
	NATEnabled              Setting[bool]   `json:"natEnabled" xml:"natEnabled" default:"true"`
	NATLeaseM               Setting[int]    `json:"natLeaseMinutes" xml:"natLeaseMinutes" default:"60"`
	NATRenewalM             Setting[int]    `json:"natRenewalMinutes" xml:"natRenewalMinutes" default:"30"`
	NATTimeoutS             Setting[int]    `json:"natTimeoutSeconds" xml:"natTimeoutSeconds" default:"10"`
	URAccepted              Setting[int]    `json:"urAccepted" xml:"urAccepted"`
	URSeen                  Setting[int]    `json:"urSeen" xml:"urSeen"`
	URUniqueID              Setting[string] `json:"urUniqueId" xml:"urUniqueID"`
	URURL                   Setting[string] `json:"urURL" xml:"urURL" default:"https://data.syncthing.net/newdata"`
	URPostInsecurely        Setting[bool]   `json:"urPostInsecurely" xml:"urPostInsecurely" default:"false"`
	URInitialDelayS         Setting[int]    `json:"urInitialDelayS" xml:"urInitialDelayS" default:"1800"`
	AutoUpgradeIntervalH    Setting[int]    `json:"autoUpgradeIntervalH" xml:"autoUpgradeIntervalH" default:"12"`
	UpgradeToPreReleases    Setting[bool]   `json:"upgradeToPreReleases" xml:"upgradeToPreReleases"`
	KeepTemporariesH        Setting[int]    `json:"keepTemporariesH" xml:"keepTemporariesH" default:"24"`
	CacheIgnoredFiles       Setting[bool]   `json:"cacheIgnoredFiles" xml:"cacheIgnoredFiles" default:"false"`
	ProgressUpdateIntervalS Setting[int]    `json:"progressUpdateIntervalS" xml:"progressUpdateIntervalS" default:"5"`
	LimitBandwidthInLan     Setting[bool]   `json:"limitBandwidthInLan" xml:"limitBandwidthInLan" default:"false"`
	// MinHomeDiskFree             Size     `json:"minHomeDiskFree" xml:"minHomeDiskFree" default:"1 %"`
	ReleasesURL                 Setting[string] `json:"releasesURL" xml:"releasesURL" default:"https://upgrades.syncthing.net/meta.json"`
	AlwaysLocalNets             []string        `json:"alwaysLocalNets" xml:"alwaysLocalNet"`
	OverwriteRemoteDevNames     Setting[bool]   `json:"overwriteRemoteDeviceNamesOnConnect" xml:"overwriteRemoteDeviceNamesOnConnect" default:"false"`
	TempIndexMinBlocks          Setting[int]    `json:"tempIndexMinBlocks" xml:"tempIndexMinBlocks" default:"10"`
	UnackedNotificationIDs      []string        `json:"unackedNotificationIDs" xml:"unackedNotificationID"`
	TrafficClass                Setting[int]    `json:"trafficClass" xml:"trafficClass"`
	DeprecatedDefaultFolderPath Setting[string] `json:"-" xml:"defaultFolderPath,omitempty"` // Deprecated: Do not use.
	SetLowPriority              Setting[bool]   `json:"setLowPriority" xml:"setLowPriority" default:"true"`
	RawMaxFolderConcurrency     Setting[int]    `json:"maxFolderConcurrency" xml:"maxFolderConcurrency"`
	CRURL                       Setting[string] `json:"crURL" xml:"crashReportingURL" default:"https://crash.syncthing.net/newcrash"`
	CREnabled                   Setting[bool]   `json:"crashReportingEnabled" xml:"crashReportingEnabled" default:"true"`
	StunKeepaliveStartS         Setting[int]    `json:"stunKeepaliveStartS" xml:"stunKeepaliveStartS" default:"180"`
	StunKeepaliveMinS           Setting[int]    `json:"stunKeepaliveMinS" xml:"stunKeepaliveMinS" default:"20"`
	RawStunServers              []string        `json:"stunServers" xml:"stunServer" default:"default"`
	RawMaxCIRequestKiB          Setting[int]    `json:"maxConcurrentIncomingRequestKiB" xml:"maxConcurrentIncomingRequestKiB"`
	AnnounceLANAddresses        Setting[bool]   `json:"announceLANAddresses" xml:"announceLANAddresses" default:"true"`
	SendFullIndexOnUpgrade      Setting[bool]   `json:"sendFullIndexOnUpgrade" xml:"sendFullIndexOnUpgrade"`
	FeatureFlags                []string        `json:"featureFlags" xml:"featureFlag"`
	AuditEnabled                Setting[bool]   `json:"auditEnabled" xml:"auditEnabled" default:"false" restart:"true"`
	AuditFile                   Setting[string] `json:"auditFile" xml:"auditFile" restart:"true"`
	// The number of connections at which we stop trying to connect to more
	// devices, zero meaning no limit. Does not affect incoming connections.
	ConnectionLimitEnough Setting[int] `json:"connectionLimitEnough" xml:"connectionLimitEnough"`
	// The maximum number of connections which we will allow in total, zero
	// meaning no limit. Affects incoming connections and prevents
	// attempting outgoing connections.
	ConnectionLimitMax                 Setting[int] `json:"connectionLimitMax" xml:"connectionLimitMax"`
	ConnectionPriorityTCPLAN           Setting[int] `json:"connectionPriorityTcpLan" xml:"connectionPriorityTcpLan" default:"10"`
	ConnectionPriorityQUICLAN          Setting[int] `json:"connectionPriorityQuicLan" xml:"connectionPriorityQuicLan" default:"20"`
	ConnectionPriorityTCPWAN           Setting[int] `json:"connectionPriorityTcpWan" xml:"connectionPriorityTcpWan" default:"30"`
	ConnectionPriorityQUICWAN          Setting[int] `json:"connectionPriorityQuicWan" xml:"connectionPriorityQuicWan" default:"40"`
	ConnectionPriorityRelay            Setting[int] `json:"connectionPriorityRelay" xml:"connectionPriorityRelay" default:"50"`
	ConnectionPriorityUpgradeThreshold Setting[int] `json:"connectionPriorityUpgradeThreshold" xml:"connectionPriorityUpgradeThreshold" default:"0"`
}
