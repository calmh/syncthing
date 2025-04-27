package config

import (
	"context"
	"io"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/proto"

	configv2 "github.com/syncthing/syncthing/internal/config/v2"
)

type (
	ChangeFunc func(oldCfg, newCfg *configv2.Configuration)
	ModifyFunc func(cur *configv2.Configuration) error
)

type Manager struct {
	current   *configv2.Configuration
	tasks     chan (func())
	listeners []ChangeFunc
}

func NewManager(cfg *configv2.Configuration) *Manager {
	return &Manager{
		current: cfg,
		tasks:   make(chan func(), 1),
	}
}

func (m *Manager) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case fn := <-m.tasks:
			fn()
		}
	}
}

func (m *Manager) Listen(l ChangeFunc) {
	m.tasks <- func() {
		m.listeners = append(m.listeners, l)
	}
}

func (m *Manager) Current() *configv2.Configuration {
	curc := make(chan *configv2.Configuration, 1)
	m.tasks <- func() {
		curc <- proto.Clone(m.current).(*configv2.Configuration)
	}
	cur := <-curc
	return cur
}

func (m *Manager) Modify(fn ModifyFunc) error {
	errC := make(chan error, 1)
	m.tasks <- func() {
		oldCfg := m.current
		newCfg := proto.Clone(m.current).(*configv2.Configuration)
		if err := fn(newCfg); err != nil {
			errC <- err
			return
		}
		m.current = newCfg
		for _, cfn := range m.listeners {
			cfn(oldCfg, newCfg)
		}
		errC <- nil
	}
	return <-errC
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
	done := make(chan struct{})
	m.tasks <- func() {
		m.current = &cfg
		close(done)
	}
	<-done
	return nil
}

func (m *Manager) WriteYAML(w io.Writer) error {
	errC := make(chan error, 1)
	m.tasks <- func() {
		bs, err := protoyaml.MarshalOptions{Indent: 4}.Marshal(m.current)
		if err != nil {
			errC <- err
			return
		}
		if _, err := w.Write(bs); err != nil {
			errC <- err
			return
		}
		errC <- nil
	}
	return <-errC
}
