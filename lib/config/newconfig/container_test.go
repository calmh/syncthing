package newconfig

import (
	"encoding/xml"
	"testing"

	"github.com/syncthing/syncthing/lib/structutil"
)

func TestContainer(t *testing.T) {
	t.Parallel()

	type config struct {
		Int    Setting[int]     `xml:"int" default:"4"`
		String Setting[string]  `xml:"string" default:"something"`
		Float  Setting[float64] `xml:"float" default:"1.5"`
		Bool   Setting[bool]    `xml:"bool" default:"false"`
	}

	var cfg config
	structutil.SetDefaults(&cfg)

	cfg.Int.valid = true
	cfg.Int.value = 42
	cfg.String.valid = true
	cfg.String.value = "four"
	cfg.Float.valid = true
	cfg.Float.value = 0.8
	cfg.Bool.valid = true
	cfg.Bool.value = true

	bs, err := xml.MarshalIndent(cfg, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bs))

	var dec config
	structutil.SetDefaults(&cfg)
	if err := xml.Unmarshal(bs, &dec); err != nil {
		t.Fatal(err)
	}

	if dec.Int.Get() != 42 {
		t.Error("wrong int value")
	}
	if dec.String.Get() != "four" {
		t.Error("wrong string value")
	}
	if dec.Float.Get() != 0.8 {
		t.Error("wrong float value")
	}
	if dec.Bool.Get() != true {
		t.Error("wrong bool value")
	}
}

func TestContainerZeroValues(t *testing.T) {
	t.Parallel()

	type config struct {
		Int    Setting[int]     `xml:"int" default:"4"`
		String Setting[string]  `xml:"string" default:"something"`
		Float  Setting[float64] `xml:"float" default:"1.5"`
		Bool   Setting[bool]    `xml:"bool" default:"true"`
	}

	var cfg config
	structutil.SetDefaults(&cfg)

	cfg.Int.Set(0)
	cfg.String.Set("")
	cfg.Float.Set(0)
	cfg.Bool.Set(false)

	bs, err := xml.MarshalIndent(cfg, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bs))

	var dec config
	structutil.SetDefaults(&dec)
	if err := xml.Unmarshal(bs, &dec); err != nil {
		t.Fatal(err)
	}

	if dec.Int.Get() != 0 {
		t.Error("wrong int value")
	}
	if dec.String.Get() != "" {
		t.Error("wrong string value")
	}
	if dec.Float.Get() != 0 {
		t.Error("wrong float value")
	}
	if dec.Bool.Get() != false {
		t.Error("wrong bool value")
	}
}

func TestContainerDefaults(t *testing.T) {
	t.Parallel()

	type config struct {
		Int    Setting[int]     `xml:"int" default:"4"`
		String Setting[string]  `xml:"string" default:"something"`
		Float  Setting[float64] `xml:"float" default:"1.5"`
		Bool   Setting[bool]    `xml:"bool" default:"true"`
	}

	var cfg config
	structutil.SetDefaults(&cfg)

	bs, err := xml.MarshalIndent(cfg, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bs))

	var dec config
	structutil.SetDefaults(&dec)
	if err := xml.Unmarshal(bs, &dec); err != nil {
		t.Fatal(err)
	}

	if dec.Int.Get() != 4 {
		t.Error("wrong int value")
	}
	if dec.String.Get() != "something" {
		t.Error("wrong string value")
	}
	if dec.Float.Get() != 1.5 {
		t.Error("wrong float value")
	}
	if dec.Bool.Get() != true {
		t.Error("wrong bool value")
	}
}

func TestOptionsStruct(t *testing.T) {
	t.Parallel()

	var cfg OptionsConfiguration
	structutil.SetDefaults(&cfg)

	bs, err := xml.MarshalIndent(cfg, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bs))
}
