package config

import (
	"crypto/sha256"
	"fmt"
	"io"
	"sync"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/proto"

	configv2 "github.com/syncthing/syncthing/internal/config/v2"
)

type Manager struct {
	mut     sync.Mutex
	current *configv2.Configuration
}

func NewManager(cfg *configv2.Configuration) *Manager {
	return &Manager{current: cfg}
}

func (m *Manager) Current() (*configv2.Configuration, string) {
	m.mut.Lock()
	cur := proto.Clone(m.current).(*configv2.Configuration)
	m.mut.Unlock()
	return cur, etag(cur)
}

func etag(cfg *configv2.Configuration) string {
	bs, _ := proto.Marshal(cfg)
	return fmt.Sprintf(`"%x"`, sha256.Sum256(bs))
}

func (m *Manager) Modify(fn func(cfg *configv2.Configuration, etag string) error) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	tmp := proto.Clone(m.current).(*configv2.Configuration)
	if err := fn(tmp, etag(tmp)); err != nil {
		return err
	}
	m.current = tmp
	return nil
}

func (m *Manager) ReadYAML(r io.Reader) error {
	bs, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	var cfg configv2.Configuration
	if err := protoyaml.Unmarshal(bs, &cfg); err != nil {
		return err
	}
	m.mut.Lock()
	m.current = &cfg
	m.mut.Unlock()
	return nil
}

func (m *Manager) WriteYAML(w io.Writer) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	bs, err := protoyaml.MarshalOptions{Indent: 2, EmitUnpopulated: true}.Marshal(m.current)
	if err != nil {
		return err
	}
	if _, err := w.Write(bs); err != nil {
		return err
	}
	return nil
}
