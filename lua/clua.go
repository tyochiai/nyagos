package lua

import "C"
import "unsafe"

func CGoBytes(p, length uintptr) []byte {
	return C.GoBytes(unsafe.Pointer(p), C.int(length))
}

func CGoStringN(p, length uintptr) string {
	return C.GoStringN((*C.char)(unsafe.Pointer(p)), C.int(length))
}