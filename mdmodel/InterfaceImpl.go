package mdmodel

type InterfaceImplRow struct {
	BaseRow
	Class     *TypeDefRow
	Interface Row //TypeDefRow/TypeRefRow/TypeSpecRow
}

type InterfaceImplTable struct {
	BaseTable[InterfaceImplRow]

	typeDefRowsMap map[*TypeDefRow][]*InterfaceImplRow
}

func (this *InterfaceImplTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	this.typeDefRowsMap = make(map[*TypeDefRow][]*InterfaceImplRow)

	TypeDefTable := this.md.Tables.TypeDef
	TypeDefOrRef := this.md.CodedIndexes.TypeDefOrRef

	for n := range this.Rows {
		row := &this.Rows[n]

		var index uint32
		readRowIndexByAddr(TypeDefTable, &addr, &index)
		row.Class = GetTableRow(TypeDefTable, index).(*TypeDefRow)
		TypeDefOrRef.ReadAddr(&addr, &row.Interface)

		this.typeDefRowsMap[row.Class] = append(
			this.typeDefRowsMap[row.Class], row)
	}
	return addr
}

func (this *InterfaceImplTable) GetByTypeDef(typeDefRow *TypeDefRow) []*InterfaceImplRow {
	return this.typeDefRowsMap[typeDefRow]
}
