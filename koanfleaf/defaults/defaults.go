package defaults

import (
	"reflect"

	"github.com/tlipoca9/errors"
)

// Defaults implements a structs provider.
type Defaults struct {
	s   any
	tag string
}

// Provider returns a provider that takes a takes a struct and a struct tag
// and uses structs to parse and provide it to koanf.
func Provider(s any, tag string) *Defaults {
	return &Defaults{s: s, tag: tag}
}

// ReadBytes is not supported by the structs provider.
func (d *Defaults) ReadBytes() ([]byte, error) {
	return nil, errors.New("defaults provider does not support this method")
}

// Read reads the struct and returns a nested config map.
func (d *Defaults) Read() (map[string]any, error) {
	t := reflect.TypeOf(d.s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("defaults provider only supports structs")
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
			k, ok := field.Tag.Lookup(d.tag)
			if !ok {
				continue
			}
			v, found := field.Tag.Lookup("default")
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

	return read(t), nil
}
