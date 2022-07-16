package mdmodel

type GenericParamConstraintRow struct {
	BaseRow
	Owner      *GenericParamRow
	Constraint Row //TypeDefRow,TypeRefRow,TypeSpecRow
}

type GenericParamConstraintTable struct {
	BaseTable[GenericParamConstraintRow]
}

func (this *GenericParamConstraintTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		ReadTableRowByAddr(&addr, this.md.Tables.GenericParam, &row.Owner)
		this.md.CodedIndexes.TypeDefOrRef.ReadAddr(&addr, &row.Constraint)
	}
	return addr
}
