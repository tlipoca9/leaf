package koanfleaf_test

import (
	"os"
	"testing"

	"github.com/tlipoca9/errors"
	"github.com/tlipoca9/leaf/koanfleaf"
)

func TestMain(m *testing.M) {
	errors.C.Style = errors.StyleStack
	code := m.Run()
	os.Exit(code)
}

func TestUnmarshalWithConf(t *testing.T) {
	type Config1 struct {
		A string            `conf:"a" default:"a_value"`
		B bool              `conf:"b"`
		C int               `conf:"c"`
		D float64           `conf:"d" default:"123.123"`
		E []string          `conf:"e" default:"1,2,3"`
		F map[string]string `conf:"f"`
		G struct {
			H string  `conf:"h"`
			I bool    `conf:"i"`
			J int     `conf:"j"`
			K float64 `conf:"k" default:"42.42"`
			L struct {
				M string `conf:"m"`
				N bool   `conf:"n"`
			} `conf:"l"`
		} `conf:"g"`
	}

	var conf1 Config1
	err := koanfleaf.UnmarshalWithConf(&conf1, koanfleaf.Conf{
		Verbose:   true,
		Tag:       "conf",
		Delim:     ".",
		EnvPrefix: "LEAF_UNITTEST_",
	})
	if err != nil {
		t.Fatalf("UnmarshalWithConf failed: %s", err.Error())
	}
}
