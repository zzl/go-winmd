package mdmodel

type GenericParamRow struct {
	BaseRow
	Number uint16
	Flags  GenericParamAttributesEnum
	Owner  Row //TypeDefRow/MethodDefRow
	Name   string
}

type GenericParamTable struct {
	BaseTable[GenericParamRow]

	typeDefRowsMap   map[*TypeDefRow][]*GenericParamRow
	methodDefRowsMap map[*MethodDefRow][]*GenericParamRow
}

func (this *GenericParamTable) ParseRows(baseAddr uintptr) uintptr {
	this.typeDefRowsMap = make(map[*TypeDefRow][]*GenericParamRow)
	this.methodDefRowsMap = make(map[*MethodDefRow][]*GenericParamRow)
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.Number)
		readAddr(&addr, &row.Flags)
		this.md.CodedIndexes.TypeOrMethodDef.ReadAddr(&addr, &row.Owner)
		this.md.StringHeap.ReadAddr(&addr, &row.Name)

		switch v := row.Owner.(type) {
		case *TypeDefRow:
			this.typeDefRowsMap[v] = append(this.typeDefRowsMap[v], row)
		case *MethodDefRow:
			this.methodDefRowsMap[v] = append(this.methodDefRowsMap[v], row)
		}
	}
	return addr
}

func (this *GenericParamTable) ListByTypeDef(typeDefRow *TypeDefRow) []*GenericParamRow {
	return this.typeDefRowsMap[typeDefRow]
}

func (this *GenericParamTable) ListByMethodDef(methodDefRow *MethodDefRow) []*GenericParamRow {
	return this.methodDefRowsMap[methodDefRow]
}
