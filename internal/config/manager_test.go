package config

import (
	"fmt"
	"os"
	"strings"
	"testing"

	configv2 "github.com/syncthing/syncthing/internal/config/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestManagerLoadYAML(t *testing.T) {
	r := strings.NewReader(`
folders:
  - id: foo
    filesystem:
      path: /bar
      minDiskFree:
        value: 1
        unit: SIZE_UNIT_PERCENT
    scanning:
      ownership: false
`)
	m := NewManager(nil)
	if err := m.ReadYAML(r); err != nil {
		t.Fatal(err)
	}

	cur, etag := m.Current()
	if etag == "" {
		t.Error("expected etag")
	}

	if len(cur.GetFolders()) != 1 {
		t.Log(cur.GetFolders())
		t.Fatal("expexted one folder")
	}

	fld := cur.GetFolders()[0]
	if fld.GetId() != "foo" {
		t.Error("bad ID", fld.GetId())
	}
	if fld.GetType() != configv2.FolderType_FOLDER_TYPE_SEND_RECEIVE {
		t.Error("bad type", fld.GetType())
	}
	if fld.GetFilesystem().GetPath() != "/bar" {
		t.Error("bad path", fld.GetFilesystem().GetPath())
	}
	if fld.GetFilesystem().GetFsType() != configv2.FilesystemType_FILESYSTEM_TYPE_BASIC {
		t.Error("bad filesystem type", fld.GetFilesystem(), fld.GetType())
	}

	m.WriteYAML(os.Stdout)
}

func TestDefaultListenAddresses(t *testing.T) {
	m := NewManager(nil)
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
					if !has {
						fmt.Println(strings.Repeat("  ", ind), "#", fld.JSONName(), ":", val.Interface())
					} else {
						fmt.Println(strings.Repeat("  ", ind), fld.JSONName(), ":", val.Interface())
					}
				}
			}
		}
	}

	r := strings.NewReader(`
folders:
  - id: foo
    filesystem:
      path: /bar
      minDiskFree:
        value: 1
        unit: SIZE_UNIT_PERCENT
    scanning:
      ownership: false
`)
	m := NewManager(nil)
	if err := m.ReadYAML(r); err != nil {
		t.Fatal(err)
	}
	cfg := m.current
	print(cfg.ProtoReflect(), 0)
}
