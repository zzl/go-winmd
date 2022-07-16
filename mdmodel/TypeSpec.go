package mdmodel

type TypeSpecRow struct {
	BaseRow
	Signature *TypeSpec
}

func (this *TypeSpecRow) GetFullTypeName() string {
	return this.Signature.String()
}

func (this *TypeSpecRow) String() string {
	return this.Signature.String()
}

type TypeSpecTable struct {
	BaseTable[TypeSpecRow]
}

func (this *TypeSpecTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.Signature, _ = ParseTypeSpec(this.md, blob.Data)
	}
	return addr
}
