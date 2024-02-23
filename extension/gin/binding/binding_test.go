package binding_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tlipoca9/leaf/extension/gin/binding"
)

type User struct {
	StringField   string        `default:"John"`
	BoolField     bool          `default:"true"`
	IntField      int           `default:"-27"`
	UintField     uint          `default:"27"`
	Int8Field     int8          `default:"27"`
	Uint8Field    uint8         `default:"27"`
	Int16Field    int16         `default:"27"`
	Uint16Field   uint16        `default:"27"`
	Int32Field    int32         `default:"27"`
	Uint32Field   uint32        `default:"27"`
	Int64Field    int64         `default:"27"`
	Uint64Field   uint64        `default:"27"`
	Float32Field  float32       `default:"27.1"`
	Float64Field  float64       `default:"27.1"`
	ByteField1    byte          `default:"27"`
	ByteField2    byte          `default:"a"`
	RuneField     rune          `default:"好"`
	TimeField     time.Time     `default:"2006-01-02T15:04:05Z"`
	DurationField time.Duration `default:"1h"`
	IpField       net.IP        `default:"1.1.1.1"`

	// Pointer fields
	StringPtrField   *string        `default:"John"`
	BoolPtrField     *bool          `default:"true"`
	IntPtrField      *int           `default:"-27"`
	UintPtrField     *uint          `default:"27"`
	Int8PtrField     *int8          `default:"27"`
	Uint8PtrField    *uint8         `default:"27"`
	Int16PtrField    *int16         `default:"27"`
	Uint16PtrField   *uint16        `default:"27"`
	Int32PtrField    *int32         `default:"27"`
	Uint32PtrField   *uint32        `default:"27"`
	Int64PtrField    *int64         `default:"27"`
	Uint64PtrField   *uint64        `default:"27"`
	Float32PtrField  *float32       `default:"27.1"`
	Float64PtrField  *float64       `default:"27.1"`
	BytePtrField1    *byte          `default:"27"`
	BytePtrField2    *byte          `default:"a"`
	RunePtrField     *rune          `default:"好"`
	TimePtrField     *time.Time     `default:"2006-01-02T15:04:05Z"`
	DurationPtrField *time.Duration `default:"1h"`
	IpPtrField       *net.IP        `default:"1.1.1.1"`
}

type House struct {
	User1 User
	User2 *User
	User
}

func TestDefaultBinding(t *testing.T) {
	assert := assert.New(t)
	assertUser := func(user User) {
		assert.Equal("John", user.StringField)
		assert.Equal(true, user.BoolField)
		assert.Equal(-27, user.IntField)
		assert.Equal(uint(27), user.UintField)
		assert.Equal(int8(27), user.Int8Field)
		assert.Equal(uint8(27), user.Uint8Field)
		assert.Equal(int16(27), user.Int16Field)
		assert.Equal(uint16(27), user.Uint16Field)
		assert.Equal(int32(27), user.Int32Field)
		assert.Equal(uint32(27), user.Uint32Field)
		assert.Equal(int64(27), user.Int64Field)
		assert.Equal(uint64(27), user.Uint64Field)
		assert.Equal(float32(27.1), user.Float32Field)
		assert.Equal(float64(27.1), user.Float64Field)
		assert.Equal(byte(27), user.ByteField1)
		assert.Equal(byte('a'), user.ByteField2)
		assert.Equal('好', user.RuneField)
		assert.Equal(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), user.TimeField)
		assert.Equal(time.Hour, user.DurationField)
		assert.Equal(net.IPv4(1, 1, 1, 1), user.IpField)
		// Pointer fields
		assert.Equal("John", *user.StringPtrField)
		assert.Equal(true, *user.BoolPtrField)
		assert.Equal(-27, *user.IntPtrField)
		assert.Equal(uint(27), *user.UintPtrField)
		assert.Equal(int8(27), *user.Int8PtrField)
		assert.Equal(uint8(27), *user.Uint8PtrField)
		assert.Equal(int16(27), *user.Int16PtrField)
		assert.Equal(uint16(27), *user.Uint16PtrField)
		assert.Equal(int32(27), *user.Int32PtrField)
		assert.Equal(uint32(27), *user.Uint32PtrField)
		assert.Equal(int64(27), *user.Int64PtrField)
		assert.Equal(uint64(27), *user.Uint64PtrField)
		assert.Equal(float32(27.1), *user.Float32PtrField)
		assert.Equal(float64(27.1), *user.Float64PtrField)
		assert.Equal(byte(27), *user.BytePtrField1)
		assert.Equal(byte('a'), *user.BytePtrField2)
		assert.Equal('好', *user.RunePtrField)
		assert.Equal(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), *user.TimePtrField)
		assert.Equal(time.Hour, *user.DurationPtrField)
		assert.Equal(net.IPv4(1, 1, 1, 1), *user.IpPtrField)
	}

	b := binding.NewDefaultBindingBuilder().Build()
	user := User{}
	err := b.Bind(nil, &user)
	if err != nil {
		t.Fatal(err)
	}
	assertUser(user)

	house := House{}
	err = b.Bind(nil, &house)
	if err != nil {
		t.Fatal(err)
	}

	assertUser(house.User)
	assertUser(house.User1)
	assertUser(*house.User2)
}
