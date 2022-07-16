package mdmodel

import (
	"errors"
	"github.com/zzl/go-win32api/win32"
	"unsafe"
)

type Model struct {
	Tables       *Tables
	CodedIndexes *CodedIndexes

	StringHeap StringHeap
	GuidHeap   GuidHeap
	BlobHeap   BlobHeap

	loadedImage *win32.LOADED_IMAGE
}

type ModelParser struct {
	//
}

func NewModelParser() *ModelParser {
	return &ModelParser{}
}

func (this *ModelParser) Parse(mdFilePath string) (*Model, error) {
	md := &Model{}
	var loadedImage win32.LOADED_IMAGE
	var pDosHeader *win32.IMAGE_DOS_HEADER

	ok, errno := win32.MapAndLoad(win32.StrToPstr(mdFilePath),
		nil, &loadedImage, win32.FALSE, win32.TRUE)

	if ok != win32.TRUE {
		return nil, errno
	}
	imageBaseAddr := uintptr(unsafe.Pointer(loadedImage.MappedAddress))

	pDosHeader = (*win32.IMAGE_DOS_HEADER)(unsafe.Pointer(loadedImage.MappedAddress))
	if pDosHeader.E_magic != win32.IMAGE_DOS_SIGNATURE {
		return nil, errors.New("bad dos signature")
	}
	imageHeader := loadedImage.FileHeader
	imageHeader = (*win32.IMAGE_NT_HEADERS64)(
		unsafe.Pointer(imageBaseAddr + uintptr(pDosHeader.E_lfanew)))
	if imageHeader.Signature != uint32(win32.IMAGE_NT_SIGNATURE) {
		return nil, errors.New("bad nt signature")
	}
	var dataDirectory []win32.IMAGE_DATA_DIRECTORY
	if uintptr(imageHeader.FileHeader.SizeOfOptionalHeader) ==
		unsafe.Sizeof(win32.IMAGE_OPTIONAL_HEADER64{}) {
		dataDirectory = imageHeader.OptionalHeader.DataDirectory[0:]
	} else {
		dataDirectory = (*win32.IMAGE_OPTIONAL_HEADER32)(
			unsafe.Pointer(&imageHeader.OptionalHeader)).DataDirectory[0:]
	}
	clrDd := dataDirectory[int(win32.IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR)]
	if clrDd.VirtualAddress == 0 {
		return nil, errors.New("??")
	}

	//
	var clrHeader *win32.IMAGE_COR20_HEADER

	clrHeader = (*win32.IMAGE_COR20_HEADER)(mappedAddrFromRva(&loadedImage, clrDd.VirtualAddress))
	//println(clrHeader.MajorRuntimeVersion)

	mdrPtr := mappedAddrFromRva(&loadedImage, clrHeader.MetaData.VirtualAddress)
	var mdr *MetadataRoot
	mdr = (*MetadataRoot)(mdrPtr)

	mdr2 := (*MetadataRoot2)(unsafe.Pointer(
		uintptr(unsafe.Pointer(&mdr.Version)) + uintptr(mdr.Length)))

	streamHeadersAddr := uintptr(unsafe.Pointer(&mdr2.StreamHeaders))
	var mti *MetadataTablesInfo
	for n := 0; n < int(mdr2.Streams); n++ {
		streamHeader := (*StreamHeader)(unsafe.Pointer(streamHeadersAddr))
		name, cb := pszToStr(&streamHeader.Name[0])
		//println(name)
		if name == "#Strings" {
			addr := uintptr(unsafe.Pointer(mdr)) + uintptr(streamHeader.Offset)
			md.StringHeap.LoadData(addr, streamHeader.Size)
		}
		if name == "#US" {
			//println("?")
		}
		if name == "#Blob" {
			addr := uintptr(unsafe.Pointer(mdr)) + uintptr(streamHeader.Offset)
			md.BlobHeap.LoadData(addr, streamHeader.Size)
		}
		if name == "#GUID" {
			addr := uintptr(unsafe.Pointer(mdr)) + uintptr(streamHeader.Offset)
			md.GuidHeap.LoadData(addr, streamHeader.Size)
		}
		if name == "#~" {
			mti = (*MetadataTablesInfo)(unsafe.Pointer(
				uintptr(unsafe.Pointer(mdr)) + uintptr(streamHeader.Offset)))
		}
		streamHeadersAddr += 8 + sizeAlign4(cb)
	}

	//
	md.StringHeap.IndexSize = 2
	if hasBit8(mti.HeapSizes, 0) {
		md.StringHeap.IndexSize = 4
	}
	md.GuidHeap.IndexSize = 2
	if hasBit8(mti.HeapSizes, 1) {
		md.GuidHeap.IndexSize = 4
	}
	md.BlobHeap.IndexSize = 2
	if hasBit8(mti.HeapSizes, 2) {
		md.BlobHeap.IndexSize = 4
	}

	var validTables []byte
	for n := byte(0); n < 64; n++ {
		if hasBit64(mti.Valid, n) {
			validTables = append(validTables, n)
		}
	}
	tableCount := len(validTables)
	rows := unsafe.Slice(&mti.Rows, tableCount)

	//
	md.Tables = NewTables(md)
	md.Tables.SetRowCounts(validTables, rows) //
	md.CodedIndexes = NewCodedIndexes(md)     //

	//
	tablesAddr := uintptr(unsafe.Pointer(&mti.Rows)) + uintptr(4*tableCount)
	md.Tables.ParseRows(tablesAddr)

	//
	md.loadedImage = &loadedImage
	return md, nil
}

func (this *Model) Close() {
	win32.UnMapAndLoad(this.loadedImage)
	this.loadedImage = nil
}

func mappedAddrFromRva(image *win32.LOADED_IMAGE, rva uint32) unsafe.Pointer {
	secCount := int(image.NumberOfSections)
	sections := unsafe.Slice(image.Sections, secCount)
	var n int
	for n = 0; n < secCount; n++ {
		sh := &sections[n]
		if sh.VirtualAddress <= rva &&
			sh.SizeOfRawData+sh.VirtualAddress > rva {
			rva -= sh.VirtualAddress
			rva += sh.PointerToRawData
			break
		}
	}
	if n >= secCount {
		println("??")
		return nil
	}
	return unsafe.Pointer(uintptr(unsafe.Pointer(image.MappedAddress)) + uintptr(rva))
}

type StreamHeader struct {
	Offset uint32
	Size   uint32
	Name   [32]byte
}

type MetadataRoot struct {
	Signature    uint32
	MajorVersion uint16
	MinorVersion uint16
	Reserved     uint32
	Length       uint32
	Version      byte //
}

type MetadataRoot2 struct {
	Flags         uint16
	Streams       uint16
	StreamHeaders StreamHeader
}

type MetadataTablesInfo struct {
	Reserved     uint32
	MajorVersion byte
	MinorVersion byte
	HeapSizes    byte
	Reserved2    byte
	Valid        uint64
	Sorted       uint64
	Rows         uint32
}

func pszToStr(psz *byte) (string, int) {
	bytes := unsafe.Slice(psz, 256)
	for n, b := range bytes {
		if b == 0 {
			return string(bytes[:n]), n + 1
		}
	}
	return "??", 0
}

func sizeAlign4(n int) uintptr {
	n2 := n / 4 * 4
	if n != n2 {
		n2 += 4
	}
	return uintptr(n2)
}

func hasBit64(bits uint64, pos byte) bool {
	val := bits & (1 << pos)
	return val > 0
}

func hasBit8(bits byte, pos byte) bool {
	val := bits & (1 << pos)
	return val > 0
}
