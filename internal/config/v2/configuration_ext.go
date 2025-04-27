package configv2

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/syncthing/syncthing/lib/protocol"
	"google.golang.org/protobuf/proto"
)

func (c *Configuration) ETag() string {
	bs, _ := proto.MarshalOptions{Deterministic: true}.Marshal(c)
	hash := sha256.Sum256(bs)
	return fmt.Sprintf(`"%s"`, base64.RawStdEncoding.EncodeToString(hash[:]))
}

func (c *Configuration) GetDevice(deviceID protocol.DeviceID) (*DeviceConfiguration, bool) {
	devStr := deviceID.String()
	for _, dev := range c.GetDevices() {
		if dev.GetDeviceId() == devStr {
			return dev, true
		}
	}
	return nil, false
}

func (c *Configuration) FolderPasswords(deviceID protocol.DeviceID) map[string]string {
	res := make(map[string]string)
	devStr := deviceID.String()
nextFolder:
	for _, fld := range c.GetFolders() {
		for _, dev := range fld.GetSharedWith() {
			if dev.GetDeviceId() == devStr {
				if pwd := dev.GetEncryptionPassword(); pwd != "" {
					res[fld.GetFolderId()] = pwd
				}
				break nextFolder
			}
		}
	}
}

func (u *URLs) GetURLsDefault(def ...string) []string {
	if u == nil {
		return def
	}
	return u.GetUrls()
}

func (c Compression) ToProtocol() protocol.Compression {
	switch c {
	case Compression_COMPRESSION_METADATA:
		return protocol.CompressionMetadata
	case Compression_COMPRESSION_ALWAYS:
		return protocol.CompressionAlways
	case Compression_COMPRESSION_NEVER:
		return protocol.CompressionNever
	default:
		return protocol.CompressionMetadata
	}
}
