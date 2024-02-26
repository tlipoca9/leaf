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
	loader := koanfleaf.NewConfigLoaderBuilder().
		Verbose(true).
		Tag("conf").
		EnvPrefix("KOANFLEAF_UNITTEST_").
		Build()

	type Config1 struct {
		A string            `conf:"a"`
		B bool              `conf:"b"`
		C int               `conf:"c"`
		D float64           `conf:"d"`
		E []string          `conf:"e"`
		F map[string]string `conf:"f"`
		G struct {
			H string  `conf:"h"`
			I bool    `conf:"i"`
			J int     `conf:"j"`
			K float64 `conf:"k"`
			L struct {
				M string `conf:"m"`
				N bool   `conf:"n"`
			} `conf:"l"`
		} `conf:"g"`
	}

	conf1 := Config1{
		A: "a_value",
		D: 123.123,
		E: []string{"1", "2", "3"},
	}
	conf1.G.K = 42.42
	err := loader.Load(&conf1)
	if err != nil {
		t.Fatalf("Load failed: %s", err.Error())
	}
}
