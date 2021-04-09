package main

//#include <string.h>
import "C"
import (
	"unsafe"
)

func MemCopy(dest, src []byte) int {
	n := len(src)
	if len(dest) < len(src) {
		n = len(dest)
	}
	if n == 0 {
		return 0
	}
	C.memcpy(unsafe.Pointer(&dest[0]), unsafe.Pointer(&src[0]), C.size_t(n))
	return n
}

func MemMove(dest, src []byte) int {
	n := len(src)
	if len(dest) < len(src) {
		n = len(dest)
	}
	if n == 0 {
		return 0
	}
	C.memmove(unsafe.Pointer(&dest[0]), unsafe.Pointer(&src[0]), C.size_t(n))
	return n
}
