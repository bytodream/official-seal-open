package collx

import "encoding/binary"

// Bytes2Int8 按小端序把字节解析为 int8。
func Bytes2Int8(bytes []byte) int8 {
	return int8(Byte2Uint16(bytes))
}

// Bytes2Int 按小端序把 8 字节解析为 int。
func Bytes2Int(bytes []byte) int {
	return int(Byte2Uint64(bytes))
}

// Bytes2Int64 按小端序把 8 字节解析为 int64。
func Bytes2Int64(bytes []byte) int64 {
	return int64(Byte2Uint64(bytes))
}

// Byte2Uint64 按小端序把 8 字节解析为 uint64。
func Byte2Uint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

// Byte2Uint32 按小端序把 4 字节解析为 uint32。
func Byte2Uint32(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

// Byte2Uint16 按小端序把 2 字节解析为 uint16。
func Byte2Uint16(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}
