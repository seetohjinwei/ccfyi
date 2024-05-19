package encoding

import (
	"bytes"
	"hash/crc64"

	"github.com/rs/zerolog/log"
)

// normal CRC-64-ECMA polynomial representation
var table = crc64.MakeTable(0x42F0E1EBA9EA3693)

func GenerateChecksum(b []byte) []byte {
	return EncodeInteger(int64(crc64.Checksum(b, table)))
}

func VerifyChecksum(b []byte, checksum []byte) bool {
	expected := EncodeInteger(int64(crc64.Checksum(b, table)))
	ret := bytes.Equal(expected, checksum)
	if !ret {
		log.Warn().Bytes("expected", expected).Bytes("actual", checksum).Msg("VerifyChecksum")
	}
	return ret
}
