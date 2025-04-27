package config

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	configv1 "github.com/syncthing/syncthing/internal/config/v1"
	configv2 "github.com/syncthing/syncthing/internal/config/v2"
	"github.com/syncthing/syncthing/lib/protocol"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestManagerLoadYAML(t *testing.T) {
	r := strings.NewReader(`
folders:
  - folderId: foo
    filesystem:
      path: /bar
    filesystemOptions:
      minDiskFree:
        percent: 1.0
    scanning:
      scanOwnership: false
`)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	m := NewManager(nil)
	go m.Serve(ctx)
	if err := m.ReadYAML(r); err != nil {
		t.Fatal(err)
	}

	cur := m.Current()
	if cur.ETag() == "" {
		t.Error("expected etag")
	}

	if len(cur.GetFolders()) != 1 {
		t.Log(cur.GetFolders())
		t.Fatal("expexted one folder")
	}

	fld := cur.GetFolders()[0]
	if fld.GetFolderId() != "foo" {
		t.Error("bad ID", fld.GetFolderId())
	}
	if fld.GetType() != configv2.FolderType_FOLDER_TYPE_SEND_RECEIVE {
		t.Error("bad type", fld.GetType())
	}
	if fld.GetFilesystem().GetPath() != "/bar" {
		t.Error("bad path", fld.GetFilesystem().GetPath())
	}
	if fld.GetFilesystem().GetType() != configv2.FilesystemType_FILESYSTEM_TYPE_BASIC {
		t.Error("bad filesystem type", fld.GetFilesystem(), fld.GetType())
	}

	m.WriteYAML(os.Stdout)
}

func TestDefaultListenAddresses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	m := NewManager(nil)
	go m.Serve(ctx)
	if m.current.GetOptions().GetListen() != nil {
		t.Error("expected nil listen")
	}

	err := m.ReadYAML(strings.NewReader(`
options:
  listen:
    urls: []
`))
	if err != nil {
		t.Fatal(err)
	}
	if m.current.GetOptions().GetListen() == nil {
		t.Error("expected non-nil listen")
	}
	if m.current.GetOptions().GetListen().GetUrls() != nil {
		t.Error("expected nil listen urls")
	}

	m = NewManager(nil)
	go m.Serve(ctx)
	err = m.ReadYAML(strings.NewReader(`
options:
  listen: {}
`))
	if err != nil {
		t.Fatal(err)
	}
	if m.current.GetOptions().GetListen() == nil {
		t.Error("expected non-nil listen")
	}
	if m.current.GetOptions().GetListen().GetUrls() != nil {
		t.Error("expected nil listen urls")
	}
}

func TestReflect(t *testing.T) {
	var print func(m protoreflect.Message, ind int)
	print = func(m protoreflect.Message, ind int) {
		desc := m.Descriptor()
		for i := 0; i < desc.Fields().Len(); i++ {
			fld := desc.Fields().Get(i)
			val := m.Get(fld)
			has := m.Has(fld)
			switch fld.Kind() {
			case protoreflect.MessageKind:
				if fld.IsList() {
					fmt.Println(strings.Repeat("  ", ind), fld.JSONName(), ":")
					lst := val.List()
					for i := 0; i < lst.Len(); i++ {
						print(lst.Get(i).Message(), ind+1)
					}
				} else {
					if !has {
						fmt.Println(strings.Repeat("  ", ind), "#", fld.JSONName(), ":")
					} else {
						fmt.Println(strings.Repeat("  ", ind), fld.JSONName(), ":")
					}
					print(val.Message(), ind+1)
				}
			default:
				if fld.IsList() {
					fmt.Println(strings.Repeat("  ", ind), fld.JSONName(), ":")
					lst := val.List()
					for i := 0; i < lst.Len(); i++ {
						fmt.Println(strings.Repeat("  ", ind+1), "-", lst.Get(i).Interface())
					}
				} else {
					ev := val.Interface()
					if fld.Kind() == protoreflect.EnumKind {
						ev = fld.Enum().Values().ByNumber(val.Enum()).Name()
					} else if fld.Kind() == protoreflect.StringKind {
						ev = strconv.Quote(val.String())
					}
					if !has {
						fmt.Println(strings.Repeat("  ", ind), "#", fld.JSONName(), ":", ev)
					} else {
						fmt.Println(strings.Repeat("  ", ind), fld.JSONName(), ":", ev)
					}
				}
			}
		}
	}

	r := strings.NewReader(`
folders:
  - folderId: foo
    filesystem:
      path: /bar
    filesystemOptions:
      minDiskFree:
        percent: 0.5
    scanning:
      scanOwnership: false
options:
  usageReporting:
    uniqueId: abcd1234
`)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	m := NewManager(nil)
	go m.Serve(ctx)
	if err := m.ReadYAML(r); err != nil {
		t.Fatal(err)
	}
	cfg := m.current
	print(cfg.ProtoReflect(), 0)
}

func TestConvertV2(t *testing.T) {
	bs, err := os.ReadFile("_testdata/config.xml")
	if err != nil {
		t.Fatal(err)
	}
	v1cfg, _, err := configv1.ReadXML(bytes.NewReader(bs), protocol.EmptyDeviceID)
	if err != nil {
		t.Fatal(err)
	}
	v2cfg := FromV1(&v1cfg)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	m := NewManager(v2cfg)
	go m.Serve(ctx)
	m.WriteYAML(os.Stdout)
}
