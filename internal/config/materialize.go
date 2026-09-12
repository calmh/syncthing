// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	syncthingv2 "github.com/syncthing/syncthing/internal/gen/syncthing/v2"
)

// MaterializeDefaults returns a copy of the configuration with the schema
// defaults filled in: every singular field that has a declared default
// and is unset is set to that default. Set fields are never changed, in
// particular not explicitly set zero values. Fields without declared
// defaults, and repeated and message fields, carry no defaults in the
// schema and are left as they are.
func MaterializeDefaults(cfg Configuration) Configuration {
	out := proto.Clone(cfg).(*syncthingv2.Configuration)
	materializeDefaults(out.ProtoReflect())
	return Configuration{out}
}

func materializeDefaults(m protoreflect.Message) {
	fields := m.Descriptor().Fields()
	for i := range fields.Len() {
		fd := fields.Get(i)
		switch {
		case fd.IsList():
			if fd.Kind() == protoreflect.MessageKind {
				list := m.Get(fd).List()
				for j := range list.Len() {
					materializeDefaults(list.Get(j).Message())
				}
			}
		case fd.IsMap():
			if fd.MapValue().Kind() == protoreflect.MessageKind {
				m.Get(fd).Map().Range(func(_ protoreflect.MapKey, value protoreflect.Value) bool {
					materializeDefaults(value.Message())
					return true
				})
			}
		case fd.Kind() == protoreflect.MessageKind:
			if m.Has(fd) {
				materializeDefaults(m.Get(fd).Message())
			}
		default:
			if !m.Has(fd) && fd.HasDefault() {
				m.Set(fd, fd.Default())
			}
		}
	}
}
