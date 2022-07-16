package mdmodel

import (
	"reflect"
)

type Tables struct {
	MetaData *Model

	//
	Module                 *ModuleTable                 //0x00 *
	TypeRef                *TypeRefTable                //0x01 *
	TypeDef                *TypeDefTable                //0x02 *
	Field                  *FieldTable                  //0x04 *
	MethodDef              *MethodDefTable              //0x06 *
	Param                  *ParamTable                  //0x08 *
	InterfaceImpl          *InterfaceImplTable          //0x09 *
	MemberRef              *MemberRefTable              //0x0A *
	Constant               *ConstantTable               //0x0B *
	CustomAttribute        *CustomAttributeTable        //0x0C *
	FieldMarshal           *FieldMarshalTable           //0x0D
	DeclSecurity           *DeclSecurityTable           //0x0E
	ClassLayout            *ClassLayoutTable            //0x0F
	FieldLayout            *FieldLayoutTable            //0x10
	StandAloneSig          *StandAloneSigTable          //0x11
	EventMap               *EventMapTable               //0x12 *
	Event                  *EventTable                  //0x14 *
	PropertyMap            *PropertyMapTable            //0x15 *
	Property               *PropertyTable               //0x17 *
	MethodSemantics        *MethodSemanticsTable        //0x18 *
	MethodImpl             *MethodImplTable             //0x19 *
	ModuleRef              *ModuleRefTable              //0x1A **
	TypeSpec               *TypeSpecTable               //0x1B *
	ImplMap                *ImplMapTable                //0x1C **
	FieldRva               *FieldRvaTable               //0x1D
	Assembly               *AssemblyTable               //0x20 *
	AssemblyProcessor      *AssemblyProcessorTable      //0x21
	AssemblyOS             *AssemblyOSTable             //0x22 *
	AssemblyRef            *AssemblyRefTable            //0x23
	AssemblyRefProcessor   *AssemblyRefProcessorTable   //0x24
	AssemblyRefOS          *AssemblyRefOSTable          //0x25
	File                   *FileTable                   //0x26
	ExportedType           *ExportedTypeTable           //0x27
	ManifestResource       *ManifestResourceTable       //0x28
	NestedClass            *NestedClassTable            //0x29 **
	GenericParam           *GenericParamTable           //0x2A
	MethodSpec             *MethodSpecTable             //0x2B
	GenericParamConstraint *GenericParamConstraintTable //0x2C

	//
	tables []Table
}

func NewTables(metaData *Model) *Tables {
	tables := &Tables{MetaData: metaData}
	tables.init()
	return tables
}

func (this *Tables) init() {
	tablesValue := reflect.ValueOf(this).Elem()
	tableType := reflect.TypeOf((*Table)(nil)).Elem()
	for n := 0; n < tablesValue.NumField(); n++ {
		fieldValue := tablesValue.Field(n)
		if fieldValue.Type().Implements(tableType) {
			pTableValue := reflect.New(fieldValue.Type().Elem())
			fieldValue.Set(pTableValue)
			this.tables = append(this.tables, pTableValue.Interface().(Table))
		}
	}
	tableCodes := []byte{
		0x00, 0x01, 0x02, 0x04, 0x06, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
		0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x14, 0x15, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x20, 0x21, 0x22, 0x23, 0x24,
		0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C,
	}
	for n, table := range this.tables {
		table.SetMetaData(this.MetaData)
		table.SetCode(tableCodes[n])
	}
}

func (this *Tables) SetRowCounts(codes []byte, rowCounts []uint32) {
	codeTableMap := make(map[byte]Table)
	for _, table := range this.tables {
		codeTableMap[table.GetCode()] = table
	}
	for n, code := range codes {
		rowCount := int(rowCounts[n])
		table := codeTableMap[code]
		table.SetRowCount(rowCount)
	}
}

func (this *Tables) ParseRows(baseAddr uintptr) {
	addr := baseAddr
	for _, table := range this.tables {
		addr = table.ParseRows(addr)
	}
}

func GetMaxRowCount(tables ...Table) int {
	maxCount := 0
	for _, table := range tables {
		if table == nil {
			continue
		}
		tableValue := reflect.ValueOf(table).Elem()
		rowsFieldValue := tableValue.FieldByName("Rows")
		count := rowsFieldValue.Len()
		if count > maxCount {
			maxCount = count
		}
	}
	return maxCount
}
