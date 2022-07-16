package mdmodel

type MethodSemanticsRow struct {
	BaseRow
	Semantics   MethodSemanticsAttributesEnum
	Method      *MethodDefRow
	Association Row //EventRow/PropertyRow
}

type MethodSemanticsTable struct {
	BaseTable[MethodSemanticsRow]
}

func (this *MethodSemanticsTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	MethodDefTable := this.md.Tables.MethodDef
	HasSemantics := this.md.CodedIndexes.HasSemantics

	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.Semantics)
		ReadTableRowByAddr(&addr, MethodDefTable, &row.Method)
		HasSemantics.ReadAddr(&addr, &row.Association)
	}
	return addr
}
