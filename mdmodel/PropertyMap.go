package mdmodel

type PropertyMapRow struct {
	BaseRow
	Parent       *TypeDefRow
	PropertyList []*PropertyRow
}

type PropertyMapTable struct {
	BaseTable[PropertyMapRow]
}

func (this *PropertyMapTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	TypeDefTable := this.md.Tables.TypeDef
	PropertyTable := this.md.Tables.Property

	var prevPropertyIndex uint32
	var prevRow *PropertyMapRow

	for n := range this.Rows {
		row := &this.Rows[n]
		ReadTableRowByAddr(&addr, TypeDefTable, &row.Parent)

		var propertyIndex uint32
		readRowIndexByAddr(PropertyTable, &addr, &propertyIndex)
		if prevRow != nil {
			GetTableRows(PropertyTable, prevPropertyIndex,
				propertyIndex, &prevRow.PropertyList)
		}
		prevPropertyIndex = propertyIndex
		prevRow = row
	}
	if prevRow != nil {
		GetTableRows(PropertyTable, prevPropertyIndex, 0, &prevRow.PropertyList)
	}
	return addr
}
