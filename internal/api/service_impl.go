package api

import (
	"cmp"
	"context"
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

func (s *ConfigurationService) AddDevice(ctx context.Context, req *connect.Request[apiv2.AddDeviceRequest]) (*connect.Response[apiv2.AddDeviceResponse], error) {
	newID := req.Msg.GetDevice().GetId()
	err := s.mgr.Modify(func(c *configv2.Configuration) error {
		devs := c.GetDevices()
		for _, dev := range devs {
			if dev.GetId() == newID {
				return fmt.Errorf("device %v already exists", newID)
			}
		}
		devs = append(devs, req.Msg.GetDevice())
		slices.SortFunc(devs, func(a, b *configv2.DeviceConfiguration) int {
			return cmp.Compare(a.GetId(), b.GetId())
		})
		c.SetDevices(devs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&apiv2.AddDeviceResponse{}), nil
}

func (s *ConfigurationService) RemoveDevice(ctx context.Context, req *connect.Request[apiv2.RemoveDeviceRequest]) (*connect.Response[apiv2.RemoveDeviceResponse], error) {
	devID := req.Msg.GetDeviceId()
	err := s.mgr.Modify(func(c *configv2.Configuration) error {
		devs := c.GetDevices()
		idx := slices.IndexFunc(devs, func(i *configv2.DeviceConfiguration) bool { return i.GetId() == devID })
		if idx < 0 {
			return fmt.Errorf("device %v does not exist", devID)
		}
		devs = append(devs[:idx], devs[idx+1:]...)
		c.SetDevices(devs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&apiv2.RemoveDeviceResponse{}), nil
}

func (s *ConfigurationService) UpdateDevice(ctx context.Context, req *connect.Request[apiv2.UpdateDeviceRequest]) (*connect.Response[apiv2.UpdateDeviceResponse], error) {
	devID := req.Msg.GetDevice().GetId()
	err := s.mgr.Modify(func(c *configv2.Configuration) error {
		devs := c.GetDevices()
		idx := slices.IndexFunc(devs, func(i *configv2.DeviceConfiguration) bool { return i.GetId() == devID })
		if idx < 0 {
			return fmt.Errorf("device %v does not exist", devID)
		}

		mask, err := fieldmask.MaskFromProtoFieldMask(req.Msg.GetUpdateMask(), camelCase)
		if err != nil {
			return err
		}
		if err := fieldmask.StructToStruct(mask, req.Msg.GetDevice(), devs[idx]); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&apiv2.UpdateDeviceResponse{}), nil
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
