package api

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	"connectrpc.com/connect"
	fieldmask "github.com/mennanov/fieldmask-utils"
	apiv2 "github.com/syncthing/syncthing/internal/api/v2"
	"github.com/syncthing/syncthing/internal/api/v2/apiv2connect"
	"github.com/syncthing/syncthing/internal/config"
	configv2 "github.com/syncthing/syncthing/internal/config/v2"
)

type ConfigurationService struct {
	mgr *config.Manager
	apiv2connect.UnimplementedConfigurationServiceHandler
}

func (s *ConfigurationService) GetConfiguration(context.Context, *connect.Request[apiv2.GetConfigurationRequest]) (*connect.Response[apiv2.GetConfigurationResponse], error) {
	cur := s.mgr.Current()
	resp := connect.NewResponse(apiv2.GetConfigurationResponse_builder{Config: cur}.Build())
	resp.Header().Set("ETag", cur.ETag())
	return resp, nil
}

func (s *ConfigurationService) AddDevice(ctx context.Context, req *connect.Request[apiv2.AddDeviceRequest]) (*connect.Response[apiv2.AddDeviceResponse], error) {
	newID := req.Msg.GetDevice().GetDeviceId()
	err := s.mgr.Modify(func(cfg *configv2.Configuration) error {
		if err := checkEtag(req, cfg.ETag()); err != nil {
			return err
		}

		devs := cfg.GetDevices()
		for _, dev := range devs {
			if dev.GetDeviceId() == newID {
				return connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("device %v already exists", newID))
			}
		}
		devs = append(devs, req.Msg.GetDevice())
		slices.SortFunc(devs, func(a, b *configv2.DeviceConfiguration) int {
			return cmp.Compare(a.GetDeviceId(), b.GetDeviceId())
		})
		cfg.SetDevices(devs)

		return nil
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&apiv2.AddDeviceResponse{}), nil
}

func (s *ConfigurationService) RemoveDevice(ctx context.Context, req *connect.Request[apiv2.RemoveDeviceRequest]) (*connect.Response[apiv2.RemoveDeviceResponse], error) {
	devID := req.Msg.GetDeviceId()
	err := s.mgr.Modify(func(cfg *configv2.Configuration) error {
		if err := checkEtag(req, cfg.ETag()); err != nil {
			return err
		}

		devs := cfg.GetDevices()
		idx := slices.IndexFunc(devs, func(i *configv2.DeviceConfiguration) bool { return i.GetDeviceId() == devID })
		if idx < 0 {
			return connect.NewError(connect.CodeNotFound, fmt.Errorf("device %v does not exist", devID))
		}
		devs = append(devs[:idx], devs[idx+1:]...)
		cfg.SetDevices(devs)

		return nil
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&apiv2.RemoveDeviceResponse{}), nil
}

func (s *ConfigurationService) UpdateDevice(ctx context.Context, req *connect.Request[apiv2.UpdateDeviceRequest]) (*connect.Response[apiv2.UpdateDeviceResponse], error) {
	devID := req.Msg.GetDevice().GetDeviceId()
	err := s.mgr.Modify(func(cfg *configv2.Configuration) error {
		if err := checkEtag(req, cfg.ETag()); err != nil {
			return err
		}

		devs := cfg.GetDevices()
		idx := slices.IndexFunc(devs, func(i *configv2.DeviceConfiguration) bool { return i.GetDeviceId() == devID })
		if idx < 0 {
			return connect.NewError(connect.CodeNotFound, fmt.Errorf("device %v does not exist", devID))
		}

		if mask := req.Msg.GetUpdateMask(); mask != nil {
			// Update the device using the given object and field mask
			mask, err := fieldmask.MaskFromProtoFieldMask(mask, camelCase)
			if err != nil {
				return connect.NewError(connect.CodeInvalidArgument, err)
			}
			if err := fieldmask.StructToStruct(mask, req.Msg.GetDevice(), devs[idx]); err != nil {
				return connect.NewError(connect.CodeInvalidArgument, err)
			}
		} else {
			// There is no field mask so the update becomes a full replace
			devs[idx] = req.Msg.GetDevice()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&apiv2.UpdateDeviceResponse{}), nil
}

func checkEtag[T any](req *connect.Request[T], etag string) error {
	if m := req.Header().Get("If-Match"); m != "" && m != etag {
		return connect.NewError(connect.CodeFailedPrecondition, errors.New("configuration modified, please retry"))
	}
	return nil
}

// https://github.com/golang/protobuf/blob/75de7c059e36b64f01d0dd234ff2fff404ec3374/protoc-gen-go/generator/generator.go#L2650
func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
