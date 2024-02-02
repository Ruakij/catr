package textDetect

import (
	"bytes"
	"unicode/utf8"
)

type Encoding int

const (
	ASCII Encoding = iota
	UTF8
	UTF16LE
	UTF16BE
	Unknown
)

func DetectEncoding(buffer []byte) Encoding {
	if IsASCII(buffer) {
		return ASCII
	}
	if IsUTF8(buffer) {
		return UTF8
	}
	if IsUTF16LE(buffer) {
		return UTF16LE
	}
	if IsUTF16BE(buffer) {
		return UTF16BE
	}
	return Unknown
}

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}
func IsUTF8(buffer []byte) bool {
	return bytes.HasPrefix(buffer, utf8BOM) || utf8.Valid(buffer)
}

var utf16LEBOM = []byte{0xFF, 0xFE}
func IsUTF16LE(buffer []byte) bool {
	if bytes.HasPrefix(buffer, utf16LEBOM) {
		return true
	}

	if len(buffer)%2 == 0 {
		isUtf16 := true
		for i := 0; i < len(buffer); i += 2 {
			if buffer[i+1] < 0xD8 || buffer[i+1] > 0xDF {
				isUtf16 = false
				break
			}
		}
		return isUtf16
	}
	return false
}

var utf16BEBOM = []byte{0xFE, 0xFF}
func IsUTF16BE(buffer []byte) bool {
	if bytes.HasPrefix(buffer, utf16BEBOM) {
		return true
	}

	// Check for UTF-16 Big Endian
	if len(buffer)%2 == 0 {
		isUtf16 := true
		for i := 0; i < len(buffer); i += 2 {
			if buffer[i] < 0xD8 || buffer[i] > 0xDF {
				isUtf16 = false
				break
			}
		}
		return isUtf16
	}
	return false
}

func IsUTF16(buffer []byte) bool {
	return IsUTF16BE(buffer) || IsUTF16LE(buffer)
}

func IsASCII(buffer []byte) bool {
	// Check for ASCII
	for _, b := range buffer {
		if (b < 32 || b > 126) && b != '\n' && b != '\r' && b != '\t' {
			return false
		}
	}
	return true
}
