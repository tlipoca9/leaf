package koanfleaf

import (
	"reflect"
	"strconv"

	"github.com/go-viper/mapstructure/v2"
)

// StringToBasicTypeHookFunc returns a DecodeHookFunc that converts
// strings to basic types.
// int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float32, float64, bool, byte, rune, complex64, complex128
func StringToBasicTypeHookFunc() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		StringToInt8HookFunc(),
		StringToUint8HookFunc(),
		StringToInt16HookFunc(),
		StringToUint16HookFunc(),
		StringToInt32HookFunc(),
		StringToUint32HookFunc(),
		StringToInt64HookFunc(),
		StringToUint64HookFunc(),
		StringToIntHookFunc(),
		StringToUintHookFunc(),
		StringToFloat32HookFunc(),
		StringToFloat64HookFunc(),
		StringToBoolHookFunc(),
		// byte and rune are aliases for uint8 and int32 respectively
		// StringToByteHookFunc(),
		// StringToRuneHookFunc(),
		StringToComplex64HookFunc(),
		StringToComplex128HookFunc(),
	)
}

// StringToInt8HookFunc returns a DecodeHookFunc that converts
// strings to int8.
func StringToInt8HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int8 {
			return data, nil
		}

		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 8)
		return int8(i64), err
	}
}

// StringToUint8HookFunc returns a DecodeHookFunc that converts
// strings to uint8.
func StringToUint8HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint8 {
			return data, nil
		}

		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 8)
		return uint8(u64), err
	}
}

// StringToInt16HookFunc returns a DecodeHookFunc that converts
// strings to int16.
func StringToInt16HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int16 {
			return data, nil
		}

		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 16)
		return int16(i64), err
	}
}

// StringToUint16HookFunc returns a DecodeHookFunc that converts
// strings to uint16.
func StringToUint16HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint16 {
			return data, nil
		}

		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 16)
		return uint16(u64), err
	}
}

// StringToInt32HookFunc returns a DecodeHookFunc that converts
// strings to int32.
func StringToInt32HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int32 {
			return data, nil
		}

		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 32)
		return int32(i64), err
	}
}

// StringToUint32HookFunc returns a DecodeHookFunc that converts
// strings to uint32.
func StringToUint32HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint32 {
			return data, nil
		}

		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 32)
		return uint32(u64), err
	}
}

// StringToInt64HookFunc returns a DecodeHookFunc that converts
// strings to int64.
func StringToInt64HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int64 {
			return data, nil
		}

		// Convert it by parsing
		return strconv.ParseInt(data.(string), 0, 64)
	}
}

// StringToUint64HookFunc returns a DecodeHookFunc that converts
// strings to uint64.
func StringToUint64HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint64 {
			return data, nil
		}

		// Convert it by parsing
		return strconv.ParseUint(data.(string), 0, 64)
	}
}

// StringToIntHookFunc returns a DecodeHookFunc that converts
// strings to int.
func StringToIntHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int {
			return data, nil
		}

		// Convert it by parsing
		i64, err := strconv.ParseInt(data.(string), 0, 0)
		return int(i64), err
	}
}

// StringToUintHookFunc returns a DecodeHookFunc that converts
// strings to uint.
func StringToUintHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint {
			return data, nil
		}

		// Convert it by parsing
		u64, err := strconv.ParseUint(data.(string), 0, 0)
		return uint(u64), err
	}
}

// StringToFloat32HookFunc returns a DecodeHookFunc that converts
// strings to float32.
func StringToFloat32HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Float32 {
			return data, nil
		}

		// Convert it by parsing
		f64, err := strconv.ParseFloat(data.(string), 32)
		return float32(f64), err
	}
}

// StringToFloat64HookFunc returns a DecodeHookFunc that converts
// strings to float64.
func StringToFloat64HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Float64 {
			return data, nil
		}

		// Convert it by parsing
		return strconv.ParseFloat(data.(string), 64)
	}
}

// StringToBoolHookFunc returns a DecodeHookFunc that converts
// strings to bool.
func StringToBoolHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Bool {
			return data, nil
		}

		// Convert it by parsing
		return strconv.ParseBool(data.(string))
	}
}

// StringToByteHookFunc returns a DecodeHookFunc that converts
// strings to byte.
func StringToByteHookFunc() mapstructure.DecodeHookFunc {
	return StringToUint8HookFunc()
}

// StringToRuneHookFunc returns a DecodeHookFunc that converts
// strings to rune.
func StringToRuneHookFunc() mapstructure.DecodeHookFunc {
	return StringToInt32HookFunc()
}

// StringToComplex64HookFunc returns a DecodeHookFunc that converts
// strings to complex64.
func StringToComplex64HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Complex64 {
			return data, nil
		}

		// Convert it by parsing
		c128, err := strconv.ParseComplex(data.(string), 64)
		return complex64(c128), err
	}
}

// StringToComplex128HookFunc returns a DecodeHookFunc that converts
// strings to complex128.
func StringToComplex128HookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Complex128 {
			return data, nil
		}

		// Convert it by parsing
		return strconv.ParseComplex(data.(string), 128)
	}
}
