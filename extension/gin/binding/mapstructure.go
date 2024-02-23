package binding

import (
	"reflect"
	"strconv"

	"github.com/go-viper/mapstructure/v2"
	"github.com/tlipoca9/errors"
)

func StringToBasicTypeHookFunc() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		StringToBoolHookFunc(),
		StringToIntHookFunc(),
		StringToUintHookFunc(),
		StringToInt8HookFunc(),
		StringToUint8HookFunc(),
		StringToInt16HookFunc(),
		StringToUint16HookFunc(),
		StringToInt32HookFunc(),
		StringToUint32HookFunc(),
		StringToInt64HookFunc(),
		StringToUint64HookFunc(),
		StringToFloat32HookFunc(),
		StringToFloat64HookFunc(),
		StringToByteHookFunc(),
		StringToRuneHookFunc(),
	)
}

func StringToBoolHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Bool {
			return data, nil
		}

		s := data.(string)
		return strconv.ParseBool(s)
	}
}

func StringToIntHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Int {
			return data, nil
		}

		s := data.(string)
		return strconv.Atoi(s)
	}
}

func StringToUintHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Uint {
			return data, nil
		}

		s := data.(string)
		ui64, err := strconv.ParseUint(s, 0, 0)
		return uint(ui64), err
	}
}

func StringToInt8HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Int8 {
			return data, nil
		}

		s := data.(string)
		i64, err := strconv.ParseInt(s, 0, 8)
		return int8(i64), err
	}
}

func StringToUint8HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Uint8 {
			return data, nil
		}
		var err error
		defer func() {
			if err != nil {
				err = errors.Wrapf(err, "unable to convert string %s to uint8", data)
			}
		}()

		s := data.(string)
		ui64, err := strconv.ParseUint(s, 0, 8)
		if err == nil {
			return uint8(ui64), nil
		}

		if len(s) == 1 {
			return s[0], nil
		}

		return nil, err
	}
}

func StringToInt16HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Int16 {
			return data, nil
		}

		s := data.(string)
		i64, err := strconv.ParseInt(s, 0, 16)
		return int16(i64), err
	}
}

func StringToUint16HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Uint16 {
			return data, nil
		}

		s := data.(string)
		ui64, err := strconv.ParseUint(s, 0, 16)
		return uint16(ui64), err
	}
}

func StringToInt32HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Int32 {
			return data, nil
		}
		var err error
		defer func() {
			if err != nil {
				err = errors.Wrapf(err, "unable to convert string %s to int32", data)
			}
		}()

		s := data.(string)
		i64, err := strconv.ParseInt(s, 0, 32)
		if err == nil {
			return int32(i64), nil
		}

		rs := []rune(s)
		if len(rs) == 1 {
			return rs[0], nil
		}

		return nil, err
	}
}

func StringToUint32HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Uint32 {
			return data, nil
		}

		s := data.(string)
		ui64, err := strconv.ParseUint(s, 0, 32)
		return uint32(ui64), err
	}
}

func StringToInt64HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Int64 {
			return data, nil
		}

		s := data.(string)
		i64, err := strconv.ParseInt(s, 0, 64)
		return i64, err
	}
}

func StringToUint64HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Uint64 {
			return data, nil
		}

		s := data.(string)
		ui64, err := strconv.ParseUint(s, 0, 64)
		return ui64, err
	}
}

func StringToFloat32HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Float32 {
			return data, nil
		}

		s := data.(string)
		f64, err := strconv.ParseFloat(s, 32)
		return float32(f64), err
	}
}

func StringToFloat64HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Float64 {
			return data, nil
		}

		s := data.(string)
		f64, err := strconv.ParseFloat(s, 64)
		return f64, err
	}
}

func StringToByteHookFunc() mapstructure.DecodeHookFunc {
	return StringToUint8HookFunc()
}

func StringToRuneHookFunc() mapstructure.DecodeHookFunc {
	return StringToInt32HookFunc()
}

// StringToComplex64HookFunc converts a string to a complex64.
// mapstructure does not support complex64 (Deprecated)
func StringToComplex64HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Complex64 {
			return data, nil
		}

		s := data.(string)
		c128, err := strconv.ParseComplex(s, 64)
		return complex64(c128), err
	}
}

// StringToComplex128HookFunc converts a string to a complex128.
// mapstructure does not support complex128 (Deprecated)
func StringToComplex128HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.Complex128 {
			return data, nil
		}

		s := data.(string)
		c128, err := strconv.ParseComplex(s, 128)
		return c128, err
	}
}
