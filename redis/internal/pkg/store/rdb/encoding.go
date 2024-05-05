package rdb

import "errors"

func EncodeLength(length uint) []byte {
	if length <= 63 {
		return []byte{byte(length)}
	} else if length <= 16383 {
		byte1 := byte(0b01000000 + (length >> 8))
		byte2 := byte(length & 0b11111111)
		return []byte{byte1, byte2}
	}

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

	return 0, b, errors.New("cannot decode length, it might be a number or string")
}

// TODO: string encoding & decoding
