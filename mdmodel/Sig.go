package mdmodel

import (
	"log"
	"math"
	"unsafe"
)

//II.23.2

type SigKind byte

const (
	SigKindDefault  SigKind = 0x0
	SigKindC        SigKind = 0x1
	SigKindStdCall  SigKind = 0x2
	SigKindThiCall  SigKind = 0x3
	SigKindFastCall SigKind = 0x4
	SigKindVarArg   SigKind = 0x5
	SigKindField    SigKind = 0x6
	SigKindProperty SigKind = 0x8 //?
)

type Sig struct {
	Data []byte
	Kind SigKind
}

func ParseMethodRefOrFieldSig(md *Model, data []byte) (interface{}, int) {
	if len(data) == 0 {
		return nil, 0
	}
	kind := SigKind(data[0] & 0b00001111)
	if kind == SigKindField {
		return ParseFieldSig(md, data)
	} else if kind < SigKindField {
		return ParseMethodRefSig(md, data)
	} else {
		log.Panic("?")
		return nil, 0
	}
}

type MethodRefSig struct {
	Sig
	HasThis         bool
	ExplicitThis    bool
	VarArg          bool
	RetType         *RetType
	FixedParamCount int
	Params          []*Param
}

func (this *MethodRefSig) String() string {
	s := "func ("
	for n, param := range this.Params {
		if n > 0 {
			s += ", "
		}
		s += param.String()
	}
	s += ") " + this.RetType.String()
	return s
}

func ParseMethodRefSig(md *Model, data []byte) (*MethodRefSig, int) {
	if len(data) == 0 {
		return nil, 0
	}
	index := 0
	var cb int
	sig := &MethodRefSig{}
	sig.Data = data

	b := data[index]
	index += 1

	sig.HasThis = b&0x20 != 0
	sig.ExplicitThis = b&0x40 != 0
	if b&0x5 != 0 {
		sig.VarArg = true
	}

	paramCount, cb := DecompressUInt(data[index:])
	index += cb

	sig.RetType, cb = ParseRetType(md, data[index:])
	index += cb

	for n := uint32(0); n < paramCount; n++ {
		param, cb := ParseParam(md, data[index:])
		index += cb
		sig.Params = append(sig.Params, param)
		if n < paramCount-1 {
			if data[index] == 0x41 { //sentinel
				sig.FixedParamCount = len(sig.Params)
				index += 1
			}
		}
	}

	return sig, index
}

//
type MethodDefSig struct {
	Sig
	HasThis            bool
	ExplicitThis       bool
	Default            bool
	VarArg             bool
	Generic            bool
	Generic_ParamCount uint32
	//ParamCount         int
	RetType *RetType
	Params  []*Param
}

func (this *MethodDefSig) String() string {
	s := "func ("
	for n, param := range this.Params {
		if n > 0 {
			s += ", "
		}
		s += param.String()
	}
	s += ") " + this.RetType.String()
	return s
}

func ParseMethodDefSig(md *Model, data []byte) (*MethodDefSig, int) {
	if len(data) == 0 {
		return nil, 0
	}
	index := 0
	var cb int
	sig := &MethodDefSig{}
	sig.Data = data

	b := data[index]
	index += 1

	sig.HasThis = b&0x20 != 0
	sig.ExplicitThis = b&0x40 != 0
	if b&0x10 != 0 {
		sig.Generic = true
	} else if b&0x5 != 0 {
		sig.VarArg = true
	} else {
		sig.Default = true
	}

	if sig.Generic {
		sig.Generic_ParamCount, cb = DecompressUInt(data[index:])
		index += cb
	}
	paramCount, cb := DecompressUInt(data[index:])
	index += cb

	sig.RetType, cb = ParseRetType(md, data[index:])
	index += cb

	for n := uint32(0); n < paramCount; n++ {
		param, cb := ParseParam(md, data[index:])
		index += cb
		sig.Params = append(sig.Params, param)
	}

	return sig, index
}

//
type FieldSig struct {
	Sig
	CustomMods []*CustomMod
	Type       *Type
}

func ParseFieldSig(md *Model, data []byte) (*FieldSig, int) {
	if len(data) == 0 {
		return nil, 0
	}
	index := 0
	var cb int
	sigKind := SigKind(data[0] & 0b00001111)
	index += 1
	if sigKind != SigKindField {
		return nil, 0
	}
	sig := &FieldSig{}
	sig.Data = data

	sig.CustomMods, cb = ParseCustomMods(md, data[index:])
	index += cb

	sig.Type, cb = ParseType(md, data[index:])
	index += cb
	return sig, index
}

//
type PropertySig struct {
	Sig
	HasThis    bool
	CustomMods []*CustomMod
	Type       *Type
	Params     []*Param
}

func ParsePropertySig(md *Model, data []byte) (*PropertySig, int) {
	if len(data) == 0 {
		return nil, 0
	}
	index := 0
	var cb int
	b := data[0]
	sigKind := SigKind(b & 0b00001111)
	index += 1
	if sigKind != SigKindProperty {
		return nil, 0
	}
	sig := &PropertySig{}
	sig.Data = data
	sig.HasThis = b&0x20 != 0

	sig.CustomMods, cb = ParseCustomMods(md, data[index:])
	index += cb

	paramCount, cb := DecompressUInt(data[index:])
	index += cb

	sig.Type, cb = ParseType(md, data[index:])
	index += cb

	for n := uint32(0); n < paramCount; n++ {
		param, cb := ParseParam(md, data[index:])
		index += cb
		sig.Params = append(sig.Params, param)
	}

	return sig, index
}

//
type LocalVarSig struct {
}

//?
type MethodSpec struct {
}

//
func DecompressUInt(data []byte) (value uint32, cb int) {
	b := data[0]
	if b&0b10000000 == 0 {
		return uint32(b), 1
	} else if b&0b11000000 == 0b10000000 {
		return uint32(uint16(b&0b00111111)<<8 | uint16(data[1])), 2
	} else if b&0b11100000 == 0b11000000 {
		return uint32(b&0b00111111)<<24 | uint32(data[1])<<16 |
			uint32(data[2])<<8 | uint32(data[3]), 4
	} else {
		log.Panic("??")
		return 0, 0
	}
}

func DecompressInt(data []byte) (value int32, cb int) {
	b := data[0]
	if b&0b10000000 == 0 {
		b = (b&0x1)<<7 | (b&0x1)<<6 | (b >> 1 & 0b01111111)
		value = int32(int8(b))
		return value, 1
	} else if b&0b11000000 == 0b10000000 {
		n := uint16(b)<<8 | uint16(data[1])
		n = (n&0x1)<<15 | (n&0x1)<<14 | (n&0x1)<<13 | (n >> 1 & 0b0011111111111111)
		value = int32(int16(n))
		return value, 2
	} else {
		n := uint32(b)<<24 | uint32(data[1])<<16 |
			uint32(data[2])<<8 | uint32(data[3])
		n = n<<31 | (n&0x1)<<30 | (n&0x1)<<29 | (n&0x1)<<28 |
			(n >> 1 & 0b00011111111111111111111111111111)
		value = int32(n)
		return value, 2
	}
}

//

type CustomMod struct {
	Optional       bool
	TypeDefRefSpec interface{}
}

func ParseCustomMod(md *Model, data []byte) (*CustomMod, int) {
	if len(data) == 0 {
		return nil, 0
	}
	cm := &CustomMod{}
	index := 0
	n, cb := DecompressUInt(data)
	index += cb
	switch ElementTypeEnum(n) {
	case ElementTypes.CMod_Opt:
		cm.Optional = true
	case ElementTypes.CMod_Reqd:
		break
	default:
		return nil, 0
	}
	data = data[index:]
	tdrsEncoded, cb := ParseTypeDefOrRefOrSpecEncoded(md, data)
	index += cb
	cm.TypeDefRefSpec = tdrsEncoded.TypeDefRefSpec
	return cm, index
}

func ParseCustomMods(md *Model, data []byte) ([]*CustomMod, int) {
	var customMods []*CustomMod
	index := 0
	var cb int
	var customMod *CustomMod
	for {
		customMod, cb = ParseCustomMod(md, data[index:])
		if cb == 0 {
			break
		}
		index += cb
		customMods = append(customMods, customMod)
	}
	return customMods, index
}

//
type TypeDefOrRefOrSpecEncoded struct {
	TypeDefRefSpec interface{}
}

func ParseTypeDefOrRefOrSpecEncoded(
	metaData *Model, data []byte) (*TypeDefOrRefOrSpecEncoded, int) {

	if len(data) == 0 {
		return nil, 0
	}
	index := 0
	tokenEncoded, cb := DecompressUInt(data)
	index += cb
	tableIndex := byte(tokenEncoded & 0b11)
	rowIndex := tokenEncoded >> 2

	tabs := metaData.Tables
	table := []Table{tabs.TypeDef, tabs.TypeRef, tabs.TypeSpec}[tableIndex]
	row := GetTableRow(table, rowIndex)
	return &TypeDefOrRefOrSpecEncoded{row}, index
}

//

type Param struct {
	CustomMods []*CustomMod
	ByRef      bool
	Type       *Type
	TypedByRef bool //interface{} //?
}

func (this *Param) String() string {
	return this.Type.String()
}

func ParseParam(md *Model, data []byte) (*Param, int) {
	if len(data) == 0 {
		return nil, 0
	}
	p := &Param{}
	index := 0
	var cb int

	p.CustomMods, cb = ParseCustomMods(md, data[index:])
	index += cb

	et := ElementTypeEnum(data[index])
	if et == ElementTypes.TypedByRef {
		p.TypedByRef = true
		index += 1
	} else {
		if et == ElementTypes.ByRef {
			p.ByRef = true
			index += 1
		}
		p.Type, cb = ParseType(md, data[index:])
		index += cb
	}
	if p.Type == nil {
		//println("?4")
	}
	return p, index
}

type RetType struct {
	Param
}

func ParseRetType(md *Model, data []byte) (*RetType, int) {
	param, cb := ParseParam(md, data)
	return &RetType{Param: *param}, cb
}

//
type GenericInst struct {
	Class          bool
	ValueType      bool
	TypeDefRefSpec interface{}
	GenArgTypes    []*Type
}

func ParseGenericInst(md *Model, data []byte) (*GenericInst, int) {
	if len(data) == 0 {
		return nil, 0
	}
	g := &GenericInst{}
	index := 0
	var cb int

	et := ElementTypeEnum(data[index])
	index += 1
	if et == ElementTypes.Class {
		g.Class = true
	} else if et == ElementTypes.ValueType {
		g.ValueType = true
	} else {
		log.Panic("?")
	}
	tdrsEncoded, cb := ParseTypeDefOrRefOrSpecEncoded(md, data[index:])
	index += cb
	g.TypeDefRefSpec = tdrsEncoded.TypeDefRefSpec
	genArgCount, cb := DecompressUInt(data[index:])
	index += cb
	for n := uint32(0); n < genArgCount; n++ {
		typ, cb := ParseType(md, data[index:])
		index += cb
		g.GenArgTypes = append(g.GenArgTypes, typ)
	}
	return g, index
}

//

type ArrayShape struct {
	Rank     uint32
	Sizes    []uint32
	LoBounds []int32
}

func ParseArrayShape(data []byte) (*ArrayShape, int) {
	index := 0
	var cb int
	a := &ArrayShape{}
	a.Rank, cb = DecompressUInt(data[index:])
	index += cb
	sizeCount, cb := DecompressUInt(data[index:])
	index += cb
	for n := uint32(0); n < sizeCount; n++ {
		size, cb := DecompressUInt(data[index:])
		index += cb
		a.Sizes = append(a.Sizes, size)
	}
	loBoundCount, cb := DecompressUInt(data[index:])
	index += cb
	for n := uint32(0); n < loBoundCount; n++ {
		loBound, cb := DecompressInt(data[index:])
		index += cb
		a.LoBounds = append(a.LoBounds, loBound)
	}
	return a, index
}

//
type Type struct {
	Kind ElementTypeEnum

	Primitive bool
	Void      bool
	IsString  bool
	Any       bool

	Array_Type               *Type
	Array_Shape              *ArrayShape
	Class_TypeDefRefSpec     interface{} // TypeDefRow/TypeRefRow/TypeSpecRow
	FnPtr_MethodDefRefSig    interface{} //MethodDefSig/MethodRefSig
	GenericInst              *GenericInst
	Mvar                     int
	Ptr_CustomMods           []*CustomMod
	Ptr_Type                 *Type
	SzArray_CustomMods       []*CustomMod
	SzArray_Type             *Type
	ValueType_TypeDefRefSpec interface{} // TypeDefRow/TypeRefRow/TypeSpecRow
	Var                      uint32
}

func (this *Type) String() string {
	switch this.Kind {
	case ElementTypes.Boolean:
		return "bool"
	case ElementTypes.Char:
		return "uint16"
	case ElementTypes.I1:
		return "int8"
	case ElementTypes.U1:
		return "byte"
	case ElementTypes.I2:
		return "int16"
	case ElementTypes.U2:
		return "uint16"
	case ElementTypes.I4:
		return "int32"
	case ElementTypes.U4:
		return "uint32"
	case ElementTypes.I8:
		return "int64"
	case ElementTypes.U8:
		return "uint64"
	case ElementTypes.R4:
		return "float"
	case ElementTypes.R8:
		return "double"
	case ElementTypes.I:
		return "uintptr" //int?
	case ElementTypes.U:
		return "uintptr"
	case ElementTypes.Array:
		return "[]" + this.Array_Type.String()
	case ElementTypes.Class:
		switch r := this.Class_TypeDefRefSpec.(type) {
		case *TypeDefRow:
			return r.String()
		case *TypeRefRow:
			return r.String()
		case *TypeSpecRow:
			return r.String()
		default:
			return "??"
		}
	case ElementTypes.FnPtr:
		switch s := this.FnPtr_MethodDefRefSig.(type) {
		case *MethodDefSig:
			return s.String()
		case *MethodRefSig:
			return s.String()
		}
	case ElementTypes.GenericInst:
		return "??"
	case ElementTypes.Ptr:
		return "*" + this.Ptr_Type.String()
	case ElementTypes.String:
		return "string"
	case ElementTypes.SzArray:
		return "[]" + this.SzArray_Type.String()
	case ElementTypes.ValueType:
		switch t := this.ValueType_TypeDefRefSpec.(type) {
		case *TypeDefRow:
			return t.String()
		case *TypeRefRow:
			return t.String()
		case *TypeSpecRow:
			return t.String()
		default:
			return "??"
		}
	case ElementTypes.Void: //?
		return ""
	case ElementTypes.Object:
		return "interface{}"
	default:
		return ""
	}
	return ""
}

func ParseType(md *Model, data []byte) (*Type, int) {
	if len(data) == 0 {
		return nil, 0
	}
	typ := &Type{}
	index := 0
	var cb int
	et := ElementTypeEnum(data[0])
	typ.Kind = et
	index += 1
	switch et {
	case ElementTypes.Boolean, ElementTypes.Char,
		ElementTypes.I1, ElementTypes.U1,
		ElementTypes.I2, ElementTypes.U2,
		ElementTypes.I4, ElementTypes.U4,
		ElementTypes.I8, ElementTypes.U8,
		ElementTypes.R4, ElementTypes.R8,
		ElementTypes.I, ElementTypes.U:
		typ.Primitive = true
	case ElementTypes.Array:
		typ.Array_Type, cb = ParseType(md, data[index:])
		index += cb
		typ.Array_Shape, cb = ParseArrayShape(data[index:])
		index += cb
	case ElementTypes.Class:
		tdrsEncoded, cb := ParseTypeDefOrRefOrSpecEncoded(md, data[index:])
		typ.Class_TypeDefRefSpec = tdrsEncoded.TypeDefRefSpec
		index += cb
	case ElementTypes.FnPtr:
		methodDefSig, cb := ParseMethodDefSig(md, data[index:])
		if methodDefSig != nil {
			typ.FnPtr_MethodDefRefSig = methodDefSig
			index += cb
		} else {
			methodRefSig, cb := ParseMethodRefSig(md, data[index:])
			if methodRefSig != nil {
				typ.FnPtr_MethodDefRefSig = methodDefSig
				index += cb
			} else {
				log.Panic("?")
			}
		}
	case ElementTypes.GenericInst:
		typ.GenericInst, cb = ParseGenericInst(md, data[index:])
		index += cb
	case ElementTypes.MVar:
		log.Panic("not implemented")
	case ElementTypes.Ptr:
		typ.Ptr_CustomMods, cb = ParseCustomMods(md, data[index:])
		index += cb
		typ.Ptr_Type, cb = ParseType(md, data[index:])
		index += cb
	case ElementTypes.String:
		typ.IsString = true
	case ElementTypes.SzArray:
		typ.SzArray_CustomMods, cb = ParseCustomMods(md, data[index:])
		index += cb
		typ.SzArray_Type, cb = ParseType(md, data[index:])
		index += cb
	case ElementTypes.ValueType:
		tdrsEncoded, cb := ParseTypeDefOrRefOrSpecEncoded(md, data[index:])
		typ.ValueType_TypeDefRefSpec = tdrsEncoded.TypeDefRefSpec
		index += cb
	case ElementTypes.Var:
		typ.Var, cb = DecompressUInt(data[index:])
		index += cb
	case ElementTypes.Void: //?
		typ.Void = true
	case ElementTypes.Object:
		typ.Any = true
	default:
		return nil, 0
	}
	return typ, index
}

//?
type TypeSpec struct {
	Sig

	Ptr_CustomMods []*CustomMod
	Ptr_Type       *Type

	FnPtr_MethodDefRefSig interface{} //MethodDefSig/MethodRefSig

	Array_Type  *Type
	Array_Shape *ArrayShape

	SzArray_CustomMods []*CustomMod
	SzArray_Type       *Type

	GenericInst *GenericInst
}

func (this *TypeSpec) String() string {
	if this.Ptr_Type != nil {
		return "*" + this.Ptr_Type.String()
	}
	if this.Array_Type != nil {
		return "[]" + this.Array_Type.String()
	}
	if this.SzArray_Type != nil {
		return "[]" + this.SzArray_Type.String()
	}
	return "??"
}

func ParseTypeSpec(md *Model, data []byte) (*TypeSpec, int) {
	if len(data) == 0 {
		return nil, 0
	}
	ts := &TypeSpec{}
	index := 0
	var cb int
	et := ElementTypeEnum(data[0])
	index += 1
	switch et {
	case ElementTypes.Array:
		ts.Array_Type, cb = ParseType(md, data[index:])
		index += cb
		ts.Array_Shape, cb = ParseArrayShape(data[index:])
		index += cb
	case ElementTypes.FnPtr:
		methodDefSig, cb := ParseMethodDefSig(md, data[index:])
		if methodDefSig != nil {
			ts.FnPtr_MethodDefRefSig = methodDefSig
			index += cb
		} else {
			methodRefSig, cb := ParseMethodRefSig(md, data[index:])
			if methodRefSig != nil {
				ts.FnPtr_MethodDefRefSig = methodDefSig
				index += cb
			} else {
				log.Panic("?")
			}
		}
	case ElementTypes.GenericInst:
		ts.GenericInst, cb = ParseGenericInst(md, data[index:])
		index += cb
	case ElementTypes.Ptr:
		ts.Ptr_CustomMods, cb = ParseCustomMods(md, data[index:])
		index += cb
		ts.Ptr_Type, cb = ParseType(md, data[index:])
		index += cb
	case ElementTypes.SzArray:
		log.Panic("not implemented")
	default:
		return nil, 0
	}
	return ts, index
}

//
type CustomAttrib struct {
	FixedArgs []*FixedArg
	NamedArgs []*NamedArg
}

type FixedArg struct {
	Value      ArgValue
	ArrayValue []ArgValue
}

func (this *FixedArg) ToInterface() interface{} {
	if this.ArrayValue != nil {
		var arr []interface{}
		for _, item := range this.ArrayValue {
			arr = append(arr, item.ToInterface())
		}
		return arr
	} else {
		return this.Value.ToInterface()
	}
}

type NamedArg struct {
	FixedArg
	Name string
}

type ArgValue struct {
	SimpleValue interface{}
	String      *string
}

func (this *ArgValue) ToInterface() interface{} {
	var value interface{}
	if this.String != nil {
		value = *this.String
	} else {
		value = this.SimpleValue
	}
	return value
}

func ParseCustomAttrib(md *Model, params []*Param, data []byte) *CustomAttrib {
	index := 0

	prolog := *(*uint16)(unsafe.Pointer(&data[index]))
	if prolog != 0x0001 {
		log.Panic("??")
	}
	index += 2

	attrib := &CustomAttrib{}
	//var cb int
	for _, param := range params {
		fa, cb := parseFixedArg(md, param.Type, data[index:])
		if fa == nil {
			break
		}
		attrib.FixedArgs = append(attrib.FixedArgs, fa)
		index += cb
	}
	numNamed := *(*uint16)(unsafe.Pointer(&data[index]))
	index += 2
	for n := uint16(0); n < numNamed; n++ {
		namedArg, cb := parseNamedArg(md, data[index:])
		index += cb
		attrib.NamedArgs = append(attrib.NamedArgs, namedArg)
	}
	return attrib
}

func parseNamedArg(md *Model, data []byte) (*NamedArg, int) {
	var arg NamedArg
	index := 0
	var cb int
	b := data[index]
	index += 1
	if b == 0x53 { //field
		//
	} else if b == 0x54 { //prop
		//
	} else {
		log.Panic("?")
	}
	fieldOrPropType := ElementTypeEnum(data[index]) //?
	_ = fieldOrPropType
	index += 1
	var argType *Type
	if fieldOrPropType == ElementTypes.X55 {
		ps, cb := parseSerString(data[index:])
		enumTypeName := *ps
		index += cb
		enumTypeRow := md.Tables.TypeDef.RowMap[enumTypeName]
		argType = enumTypeRow.FieldList[0].Signature.Type
	} else {
		argType = &Type{
			Kind: fieldOrPropType,
		}
	}

	fieldOrPropName, cb := parseSerString(data[index:])
	_ = fieldOrPropName
	index += cb

	fixedArg, cb := parseFixedArg(nil, argType, data[index:])
	index += cb
	arg.Name = *fieldOrPropName
	arg.FixedArg = *fixedArg
	return &arg, index
}

func parseFixedArg(md *Model, typ *Type, data []byte) (*FixedArg, int) {
	index := 0
	var cb int
	var arg FixedArg

	if typ.SzArray_Type != nil {
		elemCount, cb := DecompressUInt(data[index:])
		index += cb
		if elemCount == math.MaxUint32 {
			return nil, index //?
		}
		for n := uint32(0); n < elemCount; n++ {
			elem, cb := parseArgElem(md, typ.SzArray_Type, data[index:])
			arg.ArrayValue = append(arg.ArrayValue, elem)
			index += cb
		}
	} else {
		arg.Value, cb = parseArgElem(md, typ, data[index:])
		index += cb
	}
	return &arg, index
}

func parseArgElem(md *Model, typ *Type, data []byte) (ArgValue, int) {
	var value ArgValue
	var cb int
	if typ.Kind == ElementTypes.String || typ.Kind == ElementTypes.Class /*?*/ {
		value.String, cb = parseSerString(data)
	} else if typ.Kind == ElementTypes.X51 { //boxed?
		fieldOrPropType := ElementTypeEnum(data[0]) //?
		if fieldOrPropType == ElementTypes.SzArray {
			log.Panic("???")
		}
		value.SimpleValue, cb = parseSimpleValue(fieldOrPropType, data[1:])
	} else {
		simpleType := typ.Kind
		if simpleType == ElementTypes.ValueType { //enum?
			switch v := typ.ValueType_TypeDefRefSpec.(type) {
			case *TypeDefRow:
				simpleType = v.FieldList[0].Signature.Type.Kind
			case *TypeRefRow:
				if typeDefRow, ok := md.Tables.TypeDef.RowMap[v.TypeNamespace+"."+v.TypeName]; ok {
					simpleType = typeDefRow.FieldList[0].Signature.Type.Kind
				} else {
					simpleType = ElementTypes.I4 //?
				}
			default:
				log.Panic("?")
			}
		}
		value.SimpleValue, cb = parseSimpleValue(simpleType, data)
	}
	return value, cb
}

func parseSerString(data []byte) (*string, int) {
	index := 0
	packedLen, cb := DecompressUInt(data)
	index += cb
	if packedLen == 0xFF {
		return nil, index
	} else if packedLen == 0x00 {
		s := ""
		return &s, index
	} else {
		s := string(data[index : index+int(packedLen)])
		index += int(packedLen)
		return &s, index
	}
}

func parseSimpleValue(elementType ElementTypeEnum, data []byte) (interface{}, int) {
	switch elementType {
	case ElementTypes.Boolean:
		return data[0] != 0, 1
	case ElementTypes.Char:
		return *(*uint16)(unsafe.Pointer(&data[0])), 2
	case ElementTypes.R4:
		return *(*float32)(unsafe.Pointer(&data[0])), 4
	case ElementTypes.R8:
		return *(*float64)(unsafe.Pointer(&data[0])), 8
	case ElementTypes.I1:
		return int8(data[0]), 1
	case ElementTypes.I2:
		return *(*int16)(unsafe.Pointer(&data[0])), 2
	case ElementTypes.I4:
		return *(*int32)(unsafe.Pointer(&data[0])), 4
	case ElementTypes.I8:
		return *(*int64)(unsafe.Pointer(&data[0])), 8
	case ElementTypes.U1:
		return data[0], 1
	case ElementTypes.U2:
		return *(*uint16)(unsafe.Pointer(&data[0])), 2
	case ElementTypes.U4:
		return *(*uint32)(unsafe.Pointer(&data[0])), 4
	case ElementTypes.U8:
		return *(*uint64)(unsafe.Pointer(&data[0])), 8
	}
	log.Panic("??")
	return nil, 0
}
