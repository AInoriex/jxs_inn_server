package uuid

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

const (
	byteLen   = 14
	stringLen = 32
	separator = '-'

	separator1 = 8
	separator2 = 13
	separator3 = 18
	separator4 = 23
)

var (
	hexLen = hex.EncodedLen(byteLen)
	global = NewGenerator()
)

//获取uuid
func GetUuid() string {
	return strings.Replace(New().String(), "-", "", -1)
}

//获取uuid
func GetUuidLen(len int64) string {
	uuid := makeUUID(len)
	return (&uuid).String()
}

// New will return a new UUID utilizing the global generator
func New() *UUID {
	return global.New()
}

// Make will make a UUID utilizing the global generator
func Make() UUID {
	return global.Make()
}

// makeUUID will return an initialized UUID from a provided value
// Note: The provided value will be truncated from 8 bytes to 6 bytes
func makeUUID(value int64) (u UUID) {
	ts := toBytes(time.Now().UnixNano())
	copy(u[:8], ts)
	copy(u[8:], toBytes(value)[:6])
	return
}

// UUID represents a unique user identifier
type UUID [byteLen]byte

// String will return the string representation of a UUID
func (u *UUID) String() (out string) {
	// Make a byteslice with a capacity of stringLen
	// Note: The string length is greater than the hex length, this provides us space
	// to add our UUID separators after encoding to hexidecimal
	bs := make([]byte, stringLen)
	// Encode our UUID bytes as hexidecimal to our byteslice
	hex.Encode(bs, u[:])
	for i := 0; i < len(bs); i++ {
		// Switch on index
		switch i {
		// Separator index case
		case separator1, separator2, separator3, separator4:
			// Our index matches a separator index, prepend separator byte to this index
			prependByte(bs[i:], separator)
		}
	}

	// Convert byteslice to outbound string
	out = string(bs)
	return
}

// ExtendedString will return the microsoft format string representation of a UUID
func (u *UUID) ExtendedString() (out string) {
	out = u.String() + "0000"
	return
}

// Time will return an associated time
func (u *UUID) Time() (out time.Time) {
	// Make a byteslice with a size of 8 bytes (same size as int64)
	bs := make([]byte, 8)
	// Copy first 8 bytes of UUID to our byteslice
	copy(bs, u[:8])
	// Convert the byteslice to an int64 and parse that value as a Unix timestamp
	return time.Unix(0, toValue(bs))
}

// IsZero will return if the UUID is unset
func (u *UUID) IsZero() (isZero bool) {
	for i := 0; i < byteLen; i++ {
		if u[i] != 0 {
			return
		}
	}

	return true
}

// MarshalJSON is a JSON encoding helper func
func (u UUID) MarshalJSON() (bs []byte, err error) {
	// Make the byteslice big enough for our string plus our quotation marks
	bs = make([]byte, 0, stringLen+2)
	// Append first quotation
	bs = append(bs, '"')
	// Append string representation of our UUID
	bs = append(bs, u.String()...)
	// Append second quotation
	bs = append(bs, '"')
	return
}

// UnmarshalJSON is a JSON decoding helper func
func (u *UUID) UnmarshalJSON(bs []byte) (err error) {
	var str string
	// Unmarshal inbound bytes as a string
	if err = json.Unmarshal(bs, &str); err != nil {
		return
	}

	// Convert parsed string to a byteslice
	strbs := []byte(str)
	// Remove separators from the string byteslice
	removeSeparators(strbs)
	// Decode cleaned bytes as Hexidecimal
	// Note: We can ignore any bytes which extend past our hexLen
	hex.Decode(u[:], strbs[:hexLen])
	return
}
