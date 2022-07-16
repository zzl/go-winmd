package mdmodel

type ConstantRow struct {
	BaseRow
	Type   ElementTypeEnum
	Parent Row //ParamRow/FieldRow/PropertyRow
	Value  []byte
}

type ConstantTable struct {
	BaseTable[ConstantRow]

	fieldRowConstantMap map[uint32]*ConstantRow
}

func (this *ConstantTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	HasConstant := this.md.CodedIndexes.HasConstant
	this.fieldRowConstantMap = make(map[uint32]*ConstantRow)

	for n := range this.Rows {
		row := &this.Rows[n]

		readAddr(&addr, &row.Type)
		HasConstant.ReadAddr(&addr, &row.Parent)
		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.Value = blob.Data

		switch v := row.Parent.(type) {
		case *FieldRow:
			this.fieldRowConstantMap[v.RowIndex()] = row
		case *ParamRow:
			break
		default:
			//println("??")
		}
	}
	return addr
}

func (this *ConstantTable) GetByField(fieldRow *FieldRow) *ConstantRow {
	return this.fieldRowConstantMap[fieldRow.RowIndex()]
}
