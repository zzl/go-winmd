package mdmodel

type MethodImplRow struct {
	BaseRow
	Class             *TypeDefRow
	MethodBody        Row //MethodDefRow/MemberRefRow
	MethodDeclaration Row //MethodDefRow/MemberRefRow
}

type MethodImplTable struct {
	BaseTable[MethodImplRow]
}

func (this *MethodImplTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	TypeDefTable := this.md.Tables.TypeDef
	MethodDefOrRef := this.md.CodedIndexes.MethodDefOrRef

	for n := range this.Rows {
		row := &this.Rows[n]

		ReadTableRowByAddr(&addr, TypeDefTable, &row.Class)
		MethodDefOrRef.ReadAddr(&addr, &row.MethodBody)
		MethodDefOrRef.ReadAddr(&addr, &row.MethodDeclaration)
	}
	return addr
}
