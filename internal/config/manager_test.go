package config

import (
	"os"
	"strings"
	"testing"

	configv2 "github.com/syncthing/syncthing/internal/config/v2"
)

func TestManagerLoadYAML(t *testing.T) {
	r := strings.NewReader(`
folders:
  - id: foo
    filesystem:
      path: /bar
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
