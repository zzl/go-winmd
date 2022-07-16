package mdmodel

type FieldRow struct {
	BaseRow
	Flags     FieldAttributesEnum
	Name      string
	Signature *FieldSig
}

type FieldTable struct {
	BaseTable[FieldRow]
}

func (this *FieldTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.Flags)
		this.md.StringHeap.ReadAddr(&addr, &row.Name)
		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.Signature, _ = ParseFieldSig(this.md, blob.Data)
	}
	return addr
}
