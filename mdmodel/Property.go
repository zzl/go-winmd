package mdmodel

type PropertyRow struct {
	BaseRow
	Flags PropertyAttributesEnum
	Name  string
	Type  *PropertySig
}

type PropertyTable struct {
	BaseTable[PropertyRow]
}

func (this *PropertyTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.Flags)
		this.md.StringHeap.ReadAddr(&addr, &row.Name)

		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.Type, _ = ParsePropertySig(this.md, blob.Data)
	}
	return addr
}
