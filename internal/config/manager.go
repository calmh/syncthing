package config

import (
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

func (m *Manager) Current() *configv2.Configuration {
	m.mut.Lock()
	defer m.mut.Unlock()
	return proto.Clone(m.current).(*configv2.Configuration)
}

func (m *Manager) Modify(fn func(*configv2.Configuration) error) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	tmp := proto.Clone(m.current).(*configv2.Configuration)
	if err := fn(tmp); err != nil {
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
	bs, err := protoyaml.Marshal(m.current)
	if err != nil {
		return err
	}
	if _, err := w.Write(bs); err != nil {
		return err
	}
	return nil
}
