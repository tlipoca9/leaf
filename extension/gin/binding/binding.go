package binding

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-viper/mapstructure/v2"
	"github.com/tlipoca9/errors"
)

type DefaultBindingBuilder struct {
	data DefaultBinding
}

func NewDefaultBindingBuilder() *DefaultBindingBuilder {
	return &DefaultBindingBuilder{
		data: DefaultBinding{
			TagName: "default",
			DecodeHooks: []mapstructure.DecodeHookFunc{
				mapstructure.StringToIPHookFunc(),
				mapstructure.StringToIPNetHookFunc(),
				mapstructure.StringToNetIPAddrHookFunc(),
				mapstructure.StringToNetIPAddrPortHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
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
				StringToBasicTypeHookFunc(),
			},
		},
	}
}

func (d *DefaultBindingBuilder) TagName(tagName string) {
	d.data.TagName = tagName
}

func (d *DefaultBindingBuilder) DecodeHooks(decodeHooks ...mapstructure.DecodeHookFunc) {
	d.data.DecodeHooks = decodeHooks
}

func (d *DefaultBindingBuilder) AddDecodeHook(decodeHook mapstructure.DecodeHookFunc) {
	d.data.DecodeHooks = append(d.data.DecodeHooks, decodeHook)
}

func (d *DefaultBindingBuilder) Build() *DefaultBinding {
	return &d.data
}

var _ binding.Binding = (*DefaultBinding)(nil)

type DefaultBinding struct {
	TagName     string
	DecodeHooks []mapstructure.DecodeHookFunc
}

// Bind implements binding.Binding.
func (b *DefaultBinding) Bind(_ *http.Request, obj any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:    "json",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(b.DecodeHooks...),
		Result:     obj,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create mapstructure decoder")
	}

	t := reflect.TypeOf(obj)

	var parse func(t reflect.Type) map[string]any
	parse = func(t reflect.Type) map[string]any {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return nil
		}

		defaultValues := make(map[string]any)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			key := f.Name
			if name, ok := f.Tag.Lookup("json"); ok {
				key = name
			}

			if f.Anonymous ||
				f.Type.Kind() == reflect.Struct ||
				(f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct) {
				defaultValues[key] = parse(f.Type)
			}

			if value, ok := f.Tag.Lookup(b.Name()); ok {
				defaultValues[key] = value
			}
		}

		return defaultValues
	}

	defaultValues := parse(t)
	if len(defaultValues) == 0 {
		return nil
	}

	return decoder.Decode(defaultValues)
}

// Name implements binding.Binding.
func (b *DefaultBinding) Name() string {
	return b.TagName
}
