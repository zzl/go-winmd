package mdmodel

type NestedClassRow struct {
	BaseRow
	NestedClass    *TypeDefRow
	EnclosingClass *TypeDefRow
}

type NestedClassTable struct {
	BaseTable[NestedClassRow]
	nestedEnclosingMap map[*TypeDefRow]*TypeDefRow
}

func (this *NestedClassTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	this.nestedEnclosingMap = make(map[*TypeDefRow]*TypeDefRow)
	TypeDefTable := this.md.Tables.TypeDef

	for n := range this.Rows {
		row := &this.Rows[n]
		ReadTableRowByAddr(&addr, TypeDefTable, &row.NestedClass)
		ReadTableRowByAddr(&addr, TypeDefTable, &row.EnclosingClass)
		this.nestedEnclosingMap[row.NestedClass] = row.EnclosingClass
	}
	return addr
}

func (this *NestedClassTable) GetEnclosingType(typDefRow *TypeDefRow) *TypeDefRow {
	return this.nestedEnclosingMap[typDefRow]
}
