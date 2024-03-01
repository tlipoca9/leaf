package ginleaf

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
				mapstructure.StringToSliceHookFunc(","),
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

	var read func(t reflect.Type) map[string]any
	read = func(t reflect.Type) map[string]any {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return nil
		}

		out := make(map[string]any)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			k := field.Name
			if kk, ok := field.Tag.Lookup("json"); ok {
				k = kk
			}
			v, found := field.Tag.Lookup(b.Name())
			if !found {
				if v := read(field.Type); v != nil {
					out[k] = v
				}
				continue
			}
			out[k] = v
		}

		if len(out) == 0 {
			return nil
		}
		return out
	}

	defaultValues := read(reflect.TypeOf(obj))
	if len(defaultValues) == 0 {
		return nil
	}

	return decoder.Decode(defaultValues)
}

// Name implements binding.Binding.
func (b *DefaultBinding) Name() string {
	return b.TagName
}
