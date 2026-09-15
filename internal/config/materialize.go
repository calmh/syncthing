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

// MaterializeDefaults returns a copy of the configuration with all
// defaults filled in: every singular field that has a schema-declared
// default and is unset is set to that default, and every field that has
// a getter-declared default (for fields that cannot carry one in the
// schema, such as repeated and message fields) is set to its effective
// value. Set fields are never changed, in particular not explicitly set
// zero values.
func MaterializeDefaults(cfg Configuration) Configuration {
	out := proto.Clone(cfg).(*syncthingv2.Configuration)
	materializeDefaults(out.ProtoReflect())
	return Configuration{out}
}

// getterDefaulters materializes the defaults declared by getter
// functions on the wrapper types, keyed by message name. Adding a
// getter-declared default means adding it to the wrapper's getter and
// its materializeDefaults method; the walk picks it up from here.
var getterDefaulters = map[protoreflect.FullName]func(protoreflect.Message){
	"syncthing.v2.OptionsConfiguration": func(m protoreflect.Message) {
		if opts, ok := m.Interface().(*syncthingv2.OptionsConfiguration); ok {
			(OptionsConfiguration{opts}).materializeDefaults()
		}
	},
	"syncthing.v2.FolderConfiguration": func(m protoreflect.Message) {
		if folder, ok := m.Interface().(*syncthingv2.FolderConfiguration); ok {
			(FolderConfiguration{folder}).materializeDefaults()
		}
	},
	"syncthing.v2.DeviceConfiguration": func(m protoreflect.Message) {
		if device, ok := m.Interface().(*syncthingv2.DeviceConfiguration); ok {
			(DeviceConfiguration{device}).materializeDefaults()
		}
	},
}

func materializeDefaults(m protoreflect.Message) {
	if defaulter, ok := getterDefaulters[m.Descriptor().FullName()]; ok {
		defaulter(m)
	}
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
