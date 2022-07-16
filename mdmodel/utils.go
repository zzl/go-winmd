package mdmodel

import "unsafe"

func readAddr[T any](paddr *uintptr, v *T) {
	*v = *(*T)(unsafe.Pointer(*paddr))
	*paddr += unsafe.Sizeof(*v)
}
