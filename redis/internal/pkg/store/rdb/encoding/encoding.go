package encoding

import (
	"errors"

	"github.com/rs/zerolog/log"
)

func EncodeLength(length uint) []byte {
	if length <= 63 {
		// MSB == 00
		return []byte{byte(length)}
	} else if length <= 16383 {
		// MSB == 01
		byte1 := byte(0b01000000 + (length >> 8))
		byte2 := byte(length & 0b11111111)
		return []byte{byte1, byte2}
	}

	// MSB == 10
	return []byte{0b10000000, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)}
}

func DecodeLength(b []byte) (uint, []byte, error) {
	if len(b) < 1 {
		return 0, b, errors.New("cannot decode length from empty bytes")
	}

	// first two MSB
	switch b[0] >> 6 {
	case 0b0:
		val := uint(b[0])
		return val, b[1:], nil
	case 0b01:
		if len(b) < 2 {
			return 0, b, errors.New("length is in a bad format")
		}
		val := uint(b[0]&0b00111111)<<8 + uint(b[1])
		return val, b[2:], nil
	case 0b10:
		if len(b) < 5 {
			return 0, b, errors.New("length is in a bad format")
		}
		val := uint(b[1])<<24 + uint(b[2])<<16 + uint(b[3])<<8 + uint(b[4])
		return val, b[5:], nil
	}

	log.Error().Bytes("bytes", b).Msg("cannot decode length")

	return 0, b, errors.New("cannot decode length, it might be a number or string")
}

func EncodeString(str string) []byte {
	// encodes as a length-prefixed string (not supporting LZF compressed strings)

	length := len(str)
	b := EncodeLength(uint(length))

	for _, r := range str {
		b = append(b, byte(r))
	}

	// let the caller merge the strings, if necessary
	return b
}

func DecodeString(b []byte) (string, []byte, error) {
	length, remaining, err := DecodeLength(b)
	if err != nil {
		return "", b, err
	}
	if len(remaining) < int(length) {
		log.Error().Bytes("bytes", b).Bytes("remaining", remaining).Uint("length", length).Msg("cannot decode string (length is too short)")

		return "", b, errors.New("cannot decode string, the length is too short")
	}

	return string(remaining[:length]), remaining[length:], nil
}

// first 2 bits must be 11
// 3rd bit is 1 <=> integer is negative (ignore the complement bit - set it to 0)
func EncodeInteger(integer int64) []byte {
	prefix := byte(0b11000000)
	if integer < 0 {
		prefix += byte(0b00100000)
		integer *= -1
	}

	if integer < (1 << 8) {
		// 8 bit integer follows
		return []byte{prefix + 0, byte(integer)}
	} else if integer < (1 << 16) {
		// 16 bit integer follows
		return []byte{prefix + 1, byte(integer >> 8), byte(integer)}
	} else if integer < (1 << 32) {
		// 32 bit integer follows
		return []byte{prefix + 2, byte(integer >> 24), byte(integer >> 16), byte(integer >> 8), byte(integer)}
	}

	// 64 bit integer follows
	return []byte{prefix + 3, byte(integer >> 56), byte(integer >> 48), byte(integer >> 40), byte(integer >> 32), byte(integer >> 24), byte(integer >> 16), byte(integer >> 8), byte(integer)}
}

func DecodeInteger(b []byte) (int64, []byte, error) {
	if len(b) < 1 {
		return 0, b, errors.New("cannot decode length from empty bytes")
	}

	isNegative := ((b[0] >> 5) & 1) == 1

	switch b[0] & 0b00011111 {
	case 0:
		if len(b) < 2 {
			return 0, b, errors.New("length is in a bad format")
		}
		val := int64(b[1])
		if isNegative {
			val *= -1
		}
		return val, b[2:], nil
	case 1:
		if len(b) < 3 {
			return 0, b, errors.New("length is in a bad format")
		}
		val := int64(b[1])<<8 + int64(b[2])
		if isNegative {
			val *= -1
		}
		return val, b[3:], nil
	case 2:
		if len(b) < 5 {
			return 0, b, errors.New("length is in a bad format")
		}
		val := int64(b[1])<<24 + int64(b[2])<<16 + int64(b[3])<<8 + int64(b[4])
		if isNegative {
			val *= -1
		}
		return val, b[5:], nil
	case 3:
		if len(b) < 9 {
			return 0, b, errors.New("length is in a bad format")
		}
		val := int64(b[1])<<56 + int64(b[2])<<48 + int64(b[3])<<40 + int64(b[4])<<32 + int64(b[5])<<24 + int64(b[6])<<16 + int64(b[7])<<8 + int64(b[8])
		if isNegative {
			val *= -1
		}
		return val, b[9:], nil
	}

	log.Error().Hex("bytes", b).Msg("cannot decode integer")

	return 0, b, errors.New("cannot decode integer")
}
