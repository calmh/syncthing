package configv2

import (
	"crypto/sha256"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func (c *Configuration) ETag() string {
	bs, _ := proto.Marshal(c)
	return fmt.Sprintf(`"%x"`, sha256.Sum256(bs))
}
