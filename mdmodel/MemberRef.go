package mdmodel

type MemberRefRow struct {
	BaseRow
	Class     Row //TypeRefRow/ModuleRefRow/MethodDefRow/TypeSpecRow
	Name      string
	Signature interface{} // *MethodRefSig/*FieldSig
}

type MemberRefTable struct {
	BaseTable[MemberRefRow]
}

func (this *MemberRefTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	MemberRefParent := this.md.CodedIndexes.MemberRefParent

	for n := range this.Rows {
		row := &this.Rows[n]

		MemberRefParent.ReadAddr(&addr, &row.Class)
		this.md.StringHeap.ReadAddr(&addr, &row.Name)

		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.Signature, _ = ParseMethodRefOrFieldSig(this.md, blob.Data)
	}
	return addr
}
