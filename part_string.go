package mc

import (
	smc "github.com/surfinggo/mc"
)

func StringEnsurePrefix(s string, cut string) string {
	return smc.StringEnsurePrefix(s, cut)
}

func StringEnsureSuffix(s string, cut string) string {
	return smc.StringEnsureSuffix(s, cut)
}

func StringToInt(s string) int {
	return smc.StringToInt(s)
}

func StringToInt8(s string) int8 {
	return smc.StringToInt8(s)
}

func StringToInt16(s string) int16 {
	return smc.StringToInt16(s)
}

func StringToInt32(s string) int32 {
	return smc.StringToInt32(s)
}

func StringToInt64(s string) int64 {
	return smc.StringToInt64(s)
}

func StringToUint(s string) uint {
	return smc.StringToUint(s)
}

func StringToUint8(s string) uint8 {
	return smc.StringToUint8(s)
}

func StringToUint16(s string) uint16 {
	return smc.StringToUint16(s)
}

func StringToUint32(s string) uint32 {
	return smc.StringToUint32(s)
}

func StringToUint64(s string) uint64 {
	return smc.StringToUint64(s)
}

const (
	CharsetDefault = smc.CharsetDefault
)

func StringRandWithCharset(length int, charset string) string {
	return smc.StringRandWithCharset(length, charset)
}

func StringRand(length int) string {
	return smc.StringRand(length)
}
