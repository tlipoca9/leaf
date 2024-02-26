package koanfleaf

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"

	"github.com/go-viper/mapstructure/v2"

	"github.com/tlipoca9/errors"
)

type ConfigLoaderBuilder struct {
	data ConfigLoader
}

func NewConfigLoaderBuilder() *ConfigLoaderBuilder {
	return &ConfigLoaderBuilder{
		data: ConfigLoader{
			logger:     slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
			tag:        "config",
			delim:      ".",
			envPrefix:  "",
			dotEnvFile: ".env",
			decodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToBasicTypeHookFunc(),
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.OrComposeDecodeHookFunc(
					mapstructure.StringToTimeHookFunc(time.RFC822),
					mapstructure.StringToTimeHookFunc(time.RFC822Z),
					mapstructure.StringToTimeHookFunc(time.RFC850),
					mapstructure.StringToTimeHookFunc(time.RFC1123),
					mapstructure.StringToTimeHookFunc(time.RFC1123Z),
					mapstructure.StringToTimeHookFunc(time.RFC3339),
					mapstructure.StringToTimeHookFunc(time.RFC3339Nano),
				),
				mapstructure.StringToIPHookFunc(),
				mapstructure.StringToIPNetHookFunc(),
				mapstructure.StringToNetIPAddrHookFunc(),
				mapstructure.StringToNetIPAddrPortHookFunc(),
			),
		},
	}
}

func (b *ConfigLoaderBuilder) Verbose(verbose bool) *ConfigLoaderBuilder {
	b.data.verbose = verbose
	return b
}

func (b *ConfigLoaderBuilder) Logger(logger *slog.Logger) *ConfigLoaderBuilder {
	b.data.logger = logger
	return b
}

func (b *ConfigLoaderBuilder) Tag(tag string) *ConfigLoaderBuilder {
	b.data.tag = tag
	return b
}

func (b *ConfigLoaderBuilder) Delim(delim string) *ConfigLoaderBuilder {
	b.data.delim = delim
	return b
}

func (b *ConfigLoaderBuilder) EnvPrefix(envPrefix string) *ConfigLoaderBuilder {
	b.data.envPrefix = envPrefix
	return b
}

func (b *ConfigLoaderBuilder) DecodeHook(decodeHook mapstructure.DecodeHookFunc) *ConfigLoaderBuilder {
	b.data.decodeHook = decodeHook
	return b
}

func (b *ConfigLoaderBuilder) AppendDecodeHook(decodeHook ...mapstructure.DecodeHookFunc) *ConfigLoaderBuilder {
	b.data.decodeHook = mapstructure.ComposeDecodeHookFunc(
		b.data.decodeHook,
		mapstructure.ComposeDecodeHookFunc(decodeHook...),
	)
	return b
}

func (b *ConfigLoaderBuilder) Build() *ConfigLoader {
	return &b.data
}

type ConfigLoader struct {
	verbose bool
	logger  *slog.Logger

	tag        string
	delim      string
	envPrefix  string
	dotEnvFile string
	decodeHook mapstructure.DecodeHookFunc
}

func (c *ConfigLoader) log(msg string, args ...any) {
	if c.verbose {
		c.logger.Info(msg, args...)
	}
}

func (c *ConfigLoader) Load(o any) error {
	envCb := func(s string) string {
		key := strings.ToLower(strings.TrimPrefix(s, c.envPrefix))
		return strings.ReplaceAll(key, "_", c.delim)
	}
	k := koanf.NewWithConf(koanf.Conf{
		Delim:       c.delim,
		StrictMerge: true,
	})

	// Load default config.
	if err := k.Load(structs.ProviderWithDelim(o, c.tag, c.delim), nil); err != nil {
		return errors.Wrapf(err, "failed to load default config")
	}
	c.log("loaded default config", "config", k.All())

	// Load environment variables.
	if err := k.Load(env.Provider(c.envPrefix, c.delim, envCb), nil); err != nil {
		return errors.Wrapf(err, "failed to load environment variables")
	}
	c.log("loaded environment variables", "config", k.All())

	// Load dotenv file.
	if fileExist(c.dotEnvFile) {
		if err := k.Load(file.Provider(c.dotEnvFile), dotenv.ParserEnv(c.envPrefix, c.delim, envCb)); err != nil {
			return errors.Wrapf(err, "failed to load dotenv file")
		}
		c.log("loaded dotenv file", "filepath", c.dotEnvFile, "config", k.All())
	}

	if err := k.UnmarshalWithConf("", o, koanf.UnmarshalConf{
		Tag: c.tag,
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: c.decodeHook,
			Result:     o,
		},
	}); err != nil {
		return errors.Wrapf(err, "failed to unmarshal config")
	}
	c.log("unmarshaled config", "config", o)

	return nil
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
