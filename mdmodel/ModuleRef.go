package mdmodel

type ModuleRefRow struct {
	BaseRow
	Name string
}

type ModuleRefTable struct {
	BaseTable[ModuleRefRow]
}

func (this *ModuleRefTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		this.md.StringHeap.ReadAddr(&addr, &row.Name)
	}
	return addr
}
