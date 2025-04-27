package configv2

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func (c *Configuration) ETag() string {
	bs, _ := proto.MarshalOptions{Deterministic: true}.Marshal(c)
	hash := sha256.Sum256(bs)
	return fmt.Sprintf(`"%s"`, base64.RawStdEncoding.EncodeToString(hash[:]))
}

func (u *URLs) GetURLsDefault(def ...string) []string {
	if u == nil {
		return def
	}
	return u.GetUrls()
}
