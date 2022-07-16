package mdmodel

import (
	"bytes"
	"log"
	"syscall"
	"unsafe"
)

type StringHeap struct {
	IndexSize      uintptr
	indexStringMap map[uint32]string
	Data           []byte
}

func (this *StringHeap) LoadData(addr uintptr, cb uint32) {
	this.indexStringMap = make(map[uint32]string)
	data := unsafe.Slice((*byte)(unsafe.Pointer(addr)), cb)
	this.Data = data
	var start = 0
	for n, b := range data {
		if b == 0 {
			s := string(data[start:n])
			this.indexStringMap[uint32(start)] = s
			start = n + 1
		}
	}
}

func (this *StringHeap) GetByIndex(index int) string {
	s, ok := this.indexStringMap[uint32(index)]
	if !ok { //?
		pos0 := bytes.IndexByte(this.Data[index:], 0)
		s = string(this.Data[index : index+pos0])
	}
	return s
}

func (this *StringHeap) GetByIndexAddr(addr uintptr) string {
	if this.IndexSize == 2 {
		index := *(*uint16)(unsafe.Pointer(addr))
		return this.GetByIndex(int(index))
	} else if this.IndexSize == 4 {
		index := *(*uint32)(unsafe.Pointer(addr))
		return this.GetByIndex(int(index))
	} else {
		log.Panic("?")
		return ""
	}
}

func (this *StringHeap) ReadAddr(paddr *uintptr, ps *string) {
	*ps = this.GetByIndexAddr(*paddr)
	*paddr += this.IndexSize
}

//
type GuidHeap struct {
	IndexSize uintptr
	guids     []syscall.GUID
}

func (this *GuidHeap) LoadData(addr uintptr, cb uint32) {
	this.guids = unsafe.Slice((*syscall.GUID)(unsafe.Pointer(addr)), cb/16)
}

func (this *GuidHeap) GetByIndex(index int) syscall.GUID {
	return this.guids[index-1]
}

func (this *GuidHeap) GetByIndexAddr(addr uintptr) syscall.GUID {
	if this.IndexSize == 2 {
		index := *(*uint16)(unsafe.Pointer(addr))
		return this.GetByIndex(int(index))
	} else if this.IndexSize == 4 {
		index := *(*uint32)(unsafe.Pointer(addr))
		return this.GetByIndex(int(index))
	} else {
		log.Panic("?")
		return syscall.GUID{}
	}
}

func (this *GuidHeap) ReadAddr(paddr *uintptr, ps *syscall.GUID) {
	*ps = this.GetByIndexAddr(*paddr)
	*paddr += this.IndexSize
}

//
type Blob struct {
	Data []byte
}
type BlobHeap struct {
	IndexSize    uintptr
	indexBlobMap map[uint32]Blob
}

func (this *BlobHeap) LoadData(addr uintptr, cbHeap uint32) {
	this.indexBlobMap = make(map[uint32]Blob)
	data := unsafe.Slice((*byte)(unsafe.Pointer(addr)), cbHeap)
	for n := uint32(0); n < cbHeap; n++ {
		b := data[n]
		var start, cb uint32
		if b&0b10000000 == 0 {
			cb = uint32(b & 0b01111111)
			start = n + 1
		} else if b&0b11000000 == 0b10000000 {
			cb = uint32(uint16(b&0b00111111)<<8 | uint16(data[n+1]))
			start = n + 2
		} else if b&0b11100000 == 0b11000000 {
			cb = uint32(b&0b00111111)<<24 | uint32(data[n+1])<<16 |
				uint32(data[n+2])<<8 | uint32(data[n+3])
			start = n + 4
		} else {
			log.Panic("??")
		}
		blob := Blob{Data: data[start : start+cb]}
		this.indexBlobMap[n] = blob
		n = start + cb - 1
	}
}

func (this *BlobHeap) ReadAddr(paddr *uintptr) Blob {
	var index uint32
	if this.IndexSize == 2 {
		index = uint32(*(*uint16)(unsafe.Pointer(*paddr)))
	} else if this.IndexSize == 4 {
		index = *(*uint32)(unsafe.Pointer(*paddr))
	}
	*paddr += this.IndexSize
	blob := this.indexBlobMap[index]
	return blob
}
