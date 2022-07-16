package mdmodel

type TypeDefRow struct {
	BaseRow
	Flags         TypeAttributesEnum
	TypeName      string
	TypeNamespace string
	Extends       Row //TypeDefRow/TypeRefRow/TypeSpecRow
	FieldList     []*FieldRow
	MethodList    []*MethodDefRow
}

func (this *TypeDefRow) GetFullTypeName() string {
	if this.TypeNamespace != "" {
		return this.TypeNamespace + "." + this.TypeName
	} else {
		return this.TypeName
	}
}

func (this *TypeDefRow) String() string {
	return this.TypeNamespace + "." + this.TypeName
}

func (this *TypeDefRow) DerivesFrom(ancestorFullName string) bool {
	if this.TypeNamespace+"."+this.TypeName == ancestorFullName {
		return true
	}
	switch v := this.Extends.(type) {
	case *TypeDefRow:
		return v.DerivesFrom(ancestorFullName)
	case *TypeRefRow:
		return v.GetFullTypeName() == ancestorFullName
	}
	return false
}

type TypeDefTable struct {
	BaseTable[TypeDefRow]

	RowMap map[string]*TypeDefRow //keyed by fullname
}

func (this *TypeDefTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	TypeDefOrRef := this.md.CodedIndexes.TypeDefOrRef
	FieldTable := this.md.Tables.Field
	MethodDefTable := this.md.Tables.MethodDef

	var prevFieldIndex uint32
	var prevMethodDefIndex uint32
	_, _ = prevFieldIndex, prevMethodDefIndex
	var prevRow *TypeDefRow
	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.Flags)
		this.md.StringHeap.ReadAddr(&addr, &row.TypeName)
		this.md.StringHeap.ReadAddr(&addr, &row.TypeNamespace)
		TypeDefOrRef.ReadAddr(&addr, &row.Extends)

		//fieldlist
		var fieldIndex uint32
		readRowIndexByAddr(FieldTable, &addr, &fieldIndex)
		var methodDefIndex uint32
		readRowIndexByAddr(MethodDefTable, &addr, &methodDefIndex)

		if prevRow != nil {
			GetTableRows(FieldTable, prevFieldIndex, fieldIndex, &prevRow.FieldList)
			GetTableRows(MethodDefTable, prevMethodDefIndex,
				methodDefIndex, &prevRow.MethodList)
		}

		prevFieldIndex = fieldIndex
		prevMethodDefIndex = methodDefIndex
		prevRow = row
	}
	if prevRow != nil {
		GetTableRows(FieldTable, prevFieldIndex, 0, &prevRow.FieldList)
		GetTableRows(MethodDefTable, prevMethodDefIndex, 0, &prevRow.MethodList)
	}

	//
	this.RowMap = make(map[string]*TypeDefRow)
	for n := range this.Rows {
		row := &this.Rows[n]
		this.RowMap[row.TypeNamespace+"."+row.TypeName] = row
		for _, mRow := range row.MethodList {
			mRow.OwnerType = row
		}
	}

	return addr
}
