package uuid

import (
	crand "crypto/rand"
	"encoding/binary"
	"unsafe"
)

func toBytes(value int64) (bs []byte) {
	return (*[8]byte)(unsafe.Pointer(&value))[:]
}

func toValue(bs []byte) (value int64) {
	return *(*int64)(unsafe.Pointer(&bs[0]))
}

func prependByte(bs []byte, b byte) {
	// Iterate through byteslice in reverse, starting at the second to last value
	for i := len(bs) - 2; i > -1; i-- {
		// Set next index as the current index
		bs[i+1] = bs[i]
	}

	// Set the first index as the provided byte
	bs[0] = b
}

func removeByte(bs []byte, i int) {
	if i < 0 {
		// This should NEVER happen, if it happens - it's worth a panic
		panic("invalid index, cannot be negative")
	}

	// Iterate through all values starting at i, ending at the second to last value
	for ; i < len(bs)-1; i++ {
		// Set current index as the value from the following index
		bs[i] = bs[i+1]
	}

	// Set last byte in byteslice to zero-value
	bs[len(bs)-1] = 0
}

func generateSeed() (seed int64) {
	var buf [8]byte
	// Read 8 random bytes to our buffer
	crand.Read(buf[:])
	// Parse a base64 uint64 from the created bytes
	return int64(binary.LittleEndian.Uint64(buf[:]))
}

func removeSeparators(bs []byte) {
	// Remove the four separators from the provided byteslice
	// Note: This will leave zero'd bytes for the removed indexes at the end of the byteslice.
	// It should be noted that we remove the separators in reverse-order so we ensure
	// we don't lose any crucial bytes.
	removeByte(bs, separator4)
	removeByte(bs, separator3)
	removeByte(bs, separator2)
	removeByte(bs, separator1)
}
