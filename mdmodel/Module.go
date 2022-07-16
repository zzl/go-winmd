package mdmodel

import (
	"syscall"
	"unsafe"
)

type ModuleRow struct {
	BaseRow
	Generation uint16
	Name       string
	Mvid       syscall.GUID
	EncId      syscall.GUID
	EncBaseId  syscall.GUID
}

type ModuleTable struct {
	BaseTable[ModuleRow]
}

func (this *ModuleTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		addr += unsafe.Sizeof(row.Generation)

		this.md.StringHeap.ReadAddr(&addr, &row.Name)

		this.md.GuidHeap.ReadAddr(&addr, &row.Mvid)
		//println(win32.GuidToStr(&row.Mvid))

		addr += this.md.GuidHeap.IndexSize
		addr += this.md.GuidHeap.IndexSize //?
	}
	return addr
}
