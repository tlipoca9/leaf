package koanfleaf

import (
	"log/slog"
	"os"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/golang-cz/devslog"
	"github.com/tlipoca9/errors"
	"github.com/tlipoca9/leaf/koanfleaf/defaults"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Conf struct {
	Verbose bool
	Logger  *slog.Logger

	Tag        string
	Delim      string
	EnvPrefix  string
	DecodeHook mapstructure.DecodeHookFunc
}

func UnmarshalWithConf(o any, conf Conf) error {
	if conf.DecodeHook == nil {
		conf.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeDurationHookFunc(),
			StringToBasicTypeHookFunc(),
		)
	}
	if conf.Verbose && conf.Logger == nil {
		conf.Logger = slog.New(devslog.NewHandler(os.Stdout, nil))
	}

	k := koanf.NewWithConf(koanf.Conf{
		Delim:       conf.Delim,
		StrictMerge: true,
	})

	// Load default config.
	if err := k.Load(defaults.Provider(o, conf.Tag), nil); err != nil {
		return errors.Wrapf(err, "failed to load default config")
	}

	// Load environment variables.
	if err := k.Load(
		env.Provider(
			conf.EnvPrefix,
			conf.Delim,
			func(s string) string {
				key := strings.ToLower(strings.TrimPrefix(s, conf.EnvPrefix))
				return strings.ReplaceAll(key, "_", conf.Delim)
			},
		),
		nil,
	); err != nil {
		return errors.Wrapf(err, "failed to load environment variables")
	}

	// Load dotenv file.
	if fileExist(".env") {
		if err := k.Load(
			file.Provider(".env"),
			dotenv.ParserEnv(
				"",
				conf.Delim,
				func(s string) string {
					key := strings.ToLower(s)
					return strings.ReplaceAll(key, "_", conf.Delim)
				},
			),
		); err != nil {
			return errors.Wrapf(err, "failed to load dotenv file")
		}
	}

	if conf.Verbose {
		allConfig := k.All()
		allConfigStr := make(map[string]string)
		for k, v := range allConfig {
			allConfigStr[k] = v.(string)
		}
		conf.Logger.Info("loaded config", "loaded_conf", allConfigStr)
	}

	if err := k.UnmarshalWithConf("", o, koanf.UnmarshalConf{
		Tag: conf.Tag,
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: conf.DecodeHook,
			Result:     o,
		},
	}); err != nil {
		return errors.Wrapf(err, "failed to unmarshal config")
	}

	if conf.Verbose {
		conf.Logger.Info("unmarshaled config", "unmarshaled_conf", o)
	}

	return nil
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
