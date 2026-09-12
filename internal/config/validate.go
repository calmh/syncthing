// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"fmt"
	"slices"

	validate "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	configpb "github.com/syncthing/syncthing/internal/gen/config"
	"github.com/syncthing/syncthing/lib/protocol"
)

// Validate checks the configuration against the validation rules in the
// schema. In addition to the protovalidate rules, string fields annotated
// with the (config.device_id) option must parse as a device ID, including
// matching check digits.
func Validate(cfg Configuration) error {
	return validateConfig(cfg.Configuration)
}

func validateConfig(cfg *configpb.Configuration) error {
	if cfg == nil {
		return nil
	}
	if err := protovalidate.GlobalValidator.Validate(cfg); err != nil {
		return err
	}
	var violations []*protovalidate.Violation
	collectDeviceIDViolations(cfg.ProtoReflect(), nil, &violations)
	if len(violations) > 0 {
		return &protovalidate.ValidationError{Violations: violations}
	}
	return nil
}

// collectDeviceIDViolations walks the message tree and adds a violation
// for each populated device ID field that does not parse as a device ID.
// The field path is recorded so that errors can be located in the source
// document.
func collectDeviceIDViolations(m protoreflect.Message, path []*validate.FieldPathElement, violations *[]*protovalidate.Violation) {
	fields := m.Descriptor().Fields()
	for i := range fields.Len() {
		fd := fields.Get(i)
		if !m.Has(fd) {
			continue
		}
		fieldPath := withField(path, fd)
		switch {
		case fd.IsList():
			list := m.Get(fd).List()
			if fd.Kind() == protoreflect.MessageKind {
				for j := range list.Len() {
					collectDeviceIDViolations(list.Get(j).Message(), withIndex(fieldPath, j), violations)
				}
			} else if isDeviceIDField(fd) {
				for j := range list.Len() {
					checkDeviceID(list.Get(j).String(), withIndex(fieldPath, j), violations)
				}
			}
		case fd.IsMap():
			// Map keys are not part of the recorded field path; the schema
			// currently contains no maps of messages or device IDs.
			kind := fd.MapValue().Kind()
			m.Get(fd).Map().Range(func(_ protoreflect.MapKey, value protoreflect.Value) bool {
				switch {
				case kind == protoreflect.MessageKind:
					collectDeviceIDViolations(value.Message(), fieldPath, violations)
				case kind == protoreflect.StringKind && isDeviceIDField(fd):
					checkDeviceID(value.String(), fieldPath, violations)
				}
				return true
			})
		case fd.Kind() == protoreflect.MessageKind:
			collectDeviceIDViolations(m.Get(fd).Message(), fieldPath, violations)
		case fd.Kind() == protoreflect.StringKind && isDeviceIDField(fd):
			checkDeviceID(m.Get(fd).String(), fieldPath, violations)
		}
	}
}

// checkDeviceID adds a violation if the value is not a parseable device
// ID. An empty value is allowed and means "no device".
func checkDeviceID(value string, path []*validate.FieldPathElement, violations *[]*protovalidate.Violation) {
	if value == "" {
		return
	}
	if _, err := protocol.DeviceIDFromString(value); err != nil {
		*violations = append(*violations, &protovalidate.Violation{
			Proto: &validate.Violation{
				Field:   &validate.FieldPath{Elements: path},
				Message: new(fmt.Sprintf("invalid device ID %q: %v", value, err)),
			},
		})
	}
}

// isDeviceIDField reports whether the field is annotated with the
// (config.device_id) option.
func isDeviceIDField(fd protoreflect.FieldDescriptor) bool {
	value, _ := proto.GetExtension(fd.Options(), configpb.E_DeviceId).(bool)
	return value
}

// withField returns the field path with the field appended.
func withField(path []*validate.FieldPathElement, fd protoreflect.FieldDescriptor) []*validate.FieldPathElement {
	out := slices.Clone(path)
	return append(out, &validate.FieldPathElement{FieldName: new(string(fd.Name()))})
}

// withIndex returns the field path with a list index attached to the last
// element.
func withIndex(path []*validate.FieldPathElement, index int) []*validate.FieldPathElement {
	out := slices.Clone(path)
	last := out[len(out)-1]
	out[len(out)-1] = &validate.FieldPathElement{
		FieldName: last.FieldName,
		Subscript: &validate.FieldPathElement_Index{Index: uint64(index)},
	}
	return out
}
