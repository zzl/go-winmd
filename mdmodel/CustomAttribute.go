package mdmodel

type CustomAttributeRow struct {
	BaseRow
	Parent Row    //
	Type   Row    //MethodDefRow/MemberRefRow
	Value  []byte //CustomAttrib sig..

	ValueSig *CustomAttrib
}

type CustomAttributeTable struct {
	BaseTable[CustomAttributeRow]

	fieldAttrRowMap     map[*FieldRow][]*CustomAttributeRow
	typeDefAttrRowMap   map[*TypeDefRow][]*CustomAttributeRow
	methodDefAttrRowMap map[*MethodDefRow][]*CustomAttributeRow
	interfaceImplRowMap map[*InterfaceImplRow][]*CustomAttributeRow
}

func (this *CustomAttributeTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	this.fieldAttrRowMap = make(map[*FieldRow][]*CustomAttributeRow)
	this.typeDefAttrRowMap = make(map[*TypeDefRow][]*CustomAttributeRow)
	this.methodDefAttrRowMap = make(map[*MethodDefRow][]*CustomAttributeRow)
	this.interfaceImplRowMap = make(map[*InterfaceImplRow][]*CustomAttributeRow)

	HasCustomAttribute := this.md.CodedIndexes.HasCustomAttribute
	CustomAttributeType := this.md.CodedIndexes.CustomAttributeType
	for n := range this.Rows {
		row := &this.Rows[n]
		HasCustomAttribute.ReadAddr(&addr, &row.Parent)
		switch v := row.Parent.(type) {
		case *FieldRow:
			this.fieldAttrRowMap[v] = append(this.fieldAttrRowMap[v], row)
		case *TypeDefRow:
			this.typeDefAttrRowMap[v] = append(this.typeDefAttrRowMap[v], row)
		case *MethodDefRow:
			this.methodDefAttrRowMap[v] = append(this.methodDefAttrRowMap[v], row)
		case *InterfaceImplRow:
			this.interfaceImplRowMap[v] = append(this.interfaceImplRowMap[v], row)
		default:
			break //println("??")
		}
		CustomAttributeType.ReadAddr(&addr, &row.Type)

		var params []*Param
		switch v := row.Type.(type) {
		case *MethodDefRow:
			params = v.Signature.Params
		case *MemberRefRow:
			params = v.Signature.(*MethodRefSig).Params
		}

		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.Value = blob.Data
		row.ValueSig = ParseCustomAttrib(this.md, params, row.Value) //
	}
	return addr
}

func (this *CustomAttributeTable) ListByField(fieldRow *FieldRow) []*CustomAttributeRow {
	return this.fieldAttrRowMap[fieldRow]
}

func (this *CustomAttributeTable) ListByTypeDef(typeDefRow *TypeDefRow) []*CustomAttributeRow {
	return this.typeDefAttrRowMap[typeDefRow]
}

func (this *CustomAttributeTable) ListByMethodDef(methodDefRow *MethodDefRow) []*CustomAttributeRow {
	return this.methodDefAttrRowMap[methodDefRow]
}

func (this *CustomAttributeTable) ListByInterfaceImpl(
	interfaceImplRow *InterfaceImplRow) []*CustomAttributeRow {
	return this.interfaceImplRowMap[interfaceImplRow]
}
