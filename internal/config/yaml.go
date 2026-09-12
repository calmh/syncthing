// Copyright (C) 2026 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package config

import (
	"fmt"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/proto"

	configpb "github.com/syncthing/syncthing/internal/gen/config"
)

// Marshal returns the configuration in YAML format.
func Marshal(cfg *configpb.Configuration) ([]byte, error) {
	return (protoyaml.MarshalOptions{Indent: 2}).Marshal(cfg)
}

// Unmarshal parses a configuration in YAML format. An error is returned if
// the data is not valid YAML, does not match the schema, or violates the
// validation rules in the schema.
func Unmarshal(data []byte) (*configpb.Configuration, error) {
	var cfg configpb.Configuration
	opts := protoyaml.UnmarshalOptions{
		Validator: validator{},
	}
	if err := opts.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// validator adapts Validate to the protoyaml Validator interface.
type validator struct{}

func (validator) Validate(message proto.Message) error {
	cfg, ok := message.(*configpb.Configuration)
	if !ok {
		return fmt.Errorf("validation: unexpected message type %T", message)
	}
	return Validate(cfg)
}
