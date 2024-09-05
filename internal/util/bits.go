package util

import "encoding/binary"

func EncodeSelection(index uint64, set bool) []byte {
	bytes := make([]byte, 9)
	if set {
		bytes[0] = 1
	} else {
		bytes[0] = 0
	}
	binary.LittleEndian.PutUint64(bytes[1:], index)
	return bytes
}

func DecodeSelection(bytes []byte) (uint64, bool) {
	if len(bytes) != 9 {
		return 0, false
	}
	return binary.LittleEndian.Uint64(bytes[1:]), bytes[0] == 1
}
