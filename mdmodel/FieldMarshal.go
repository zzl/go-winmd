package mdmodel

type FieldMarshalRow struct {
	BaseRow
	Parent     Row   //FieldRow/ParamRow
	NativeType *Type //[]byte
}

type FieldMarshalTable struct {
	BaseTable[FieldMarshalRow]
}

func (this *FieldMarshalTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	HasFieldMarshall := this.md.CodedIndexes.HasFieldMarshall
	for n := range this.Rows {
		row := &this.Rows[n]

		HasFieldMarshall.ReadAddr(&addr, &row.Parent)
		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.NativeType, _ = ParseType(this.md, blob.Data)
	}
	return addr
}
