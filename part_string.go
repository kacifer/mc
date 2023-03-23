package mc

import (
	"github.com/surfinggo/mc"
)

func StringEnsurePrefix(s string, cut string) string {
	return mc.StringEnsurePrefix(s, cut)
}

func StringEnsureSuffix(s string, cut string) string {
	return mc.StringEnsureSuffix(s, cut)
}

func StringToInt(s string) int {
	return mc.StringToInt(s)
}

func StringToInt8(s string) int8 {
	return mc.StringToInt8(s)
}

func StringToInt16(s string) int16 {
	return mc.StringToInt16(s)
}

func StringToInt32(s string) int32 {
	return mc.StringToInt32(s)
}

func StringToInt64(s string) int64 {
	return mc.StringToInt64(s)
}

func StringToUint(s string) uint {
	return mc.StringToUint(s)
}

func StringToUint8(s string) uint8 {
	return mc.StringToUint8(s)
}

func StringToUint16(s string) uint16 {
	return mc.StringToUint16(s)
}

func StringToUint32(s string) uint32 {
	return mc.StringToUint32(s)
}

func StringToUint64(s string) uint64 {
	return mc.StringToUint64(s)
}

const (
	CharsetDefault = mc.CharsetDefault
)

func StringRandWithCharset(length int, charset string) string {
	return mc.StringRandWithCharset(length, charset)
}

func StringRand(length int) string {
	return mc.StringRand(length)
}
