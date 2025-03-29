package apiv2

import (
	"testing"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/proto"
)

func TestRenderConfig(t *testing.T) {
	cfgb := &Configuration_builder{}
	fcfg := &FolderConfiguration{}
	fcfg.SetId("foo")
	fcfg.SetLabel("this is the label")
	// fcfg.SetType(FolderType_FOLDER_TYPE_SEND_RECEIVE)
	fcfg.SetSyncing(FolderConfiguration_Syncing_builder{
		Copiers: proto.Int32(8),
	}.Build())
	fcfg.SetScanning(FolderConfiguration_Scanning_builder{
		RescanIntervalS: proto.Int32(3600),
		Hashers:         proto.Int32(8),
	}.Build())
	fcfg.SetFilesystem(FolderConfiguration_Filesystem_builder{
		Path: proto.String("/tmp/foo"),
	}.Build())
	cfgb.Folders = append(cfgb.Folders, fcfg)
	cfg := cfgb.Build()

	t.Log(cfg.GetFolders()[0].GetType())
	t.Log(cfg.GetFolders()[0].GetPaused())
	t.Log(cfg.GetFolders()[0].GetScanning().GetWatcher().GetEnabled())

	bs, err := protoyaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s\n", bs)
}
