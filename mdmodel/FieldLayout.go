package mdmodel

type FieldLayoutRow struct {
	BaseRow
	Offset uint32
	Field  *FieldRow
}

type FieldLayoutTable struct {
	BaseTable[FieldLayoutRow]
	fieldLayoutMap map[*FieldRow]*FieldLayoutRow
}

func (this *FieldLayoutTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	this.fieldLayoutMap = make(map[*FieldRow]*FieldLayoutRow)
	FieldTable := this.md.Tables.Field
	for n := range this.Rows {
		row := &this.Rows[n]

		readAddr(&addr, &row.Offset)
		ReadTableRowByAddr(&addr, FieldTable, &row.Field)

		this.fieldLayoutMap[row.Field] = row
	}
	return addr

}

func (this *FieldLayoutTable) GetByField(row *FieldRow) *FieldLayoutRow {
	return this.fieldLayoutMap[row]
}
