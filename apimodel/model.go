package apimodel

import (
	"fmt"
	"github.com/zzl/go-winmd/mdmodel"
	"log"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type Model struct {
	RootNamespaces []*Namespace
	AllNamespaces  []*Namespace
}

type ModelParser struct {
	rootNamespaces []*Namespace
	allNamespaces  []*Namespace

	replaceTypeMap map[string]*Type

	nsMap   map[string]*Namespace
	typeMap map[string]*Type
	mdModel *mdmodel.Model
}

func NewModelParser(replaceTypeMap map[string]*Type) *ModelParser {
	return &ModelParser{
		replaceTypeMap: replaceTypeMap,
	}
}

func (this *ModelParser) Parse(mdModel *mdmodel.Model) *Model {
	this.mdModel = mdModel
	this.nsMap = make(map[string]*Namespace)
	this.typeMap = make(map[string]*Type)

	typeDefTable := this.mdModel.Tables.TypeDef
	var rows []*mdmodel.TypeDefRow
	for n := range typeDefTable.Rows {
		row := &typeDefTable.Rows[n]
		ns := this.resolveNamespace(row.TypeNamespace)

		typ := this.parseTypeDef(row)
		if typ == nil { //
			continue
		}
		var ignoreType bool
		rootEnclosingType := this.getRootEnclosingType(row)
		const archAttrName = "Windows.Win32.Interop.SupportedArchitectureAttribute"
		if rootEnclosingType != nil {
			attrRows := this.mdModel.Tables.CustomAttribute.ListByTypeDef(rootEnclosingType)
			attrs := this.ParseAttributes(attrRows)
			for _, attr := range attrs {
				if attr.Type.FullName == archAttrName {
					if attr.Args[0].(int32)&0x02 == 0 {
						ignoreType = true
						break
					}
				}
			}
		}
		if ignoreType {
			continue
		}
		ns.Types = append(ns.Types, typ)
		typ.Namespace = ns
		rows = append(rows, row)
	}
	//

	for n := range rows {
		row := rows[n]
		if row.Flags&mdmodel.TypeAttributes.NestedPublic != 0 {
			enclosingTypeDef := this.mdModel.Tables.NestedClass.GetEnclosingType(row)
			if enclosingTypeDef == nil {
				log.Panic("??")
			}
			typ := this.typeMap[this.makeTypeDefFullName(row)]
			enclosingType := this.typeMap[this.makeTypeDefFullName(enclosingTypeDef)]
			enclosingType.NestedTypes = append(enclosingType.NestedTypes, typ) //?
			typ.EnclosingType = enclosingType
		}
	}

	//
	sort.Slice(this.allNamespaces, func(i, j int) bool {
		return this.allNamespaces[i].FullName < this.allNamespaces[j].FullName
	})
	sort.Slice(this.rootNamespaces, func(i, j int) bool {
		return this.rootNamespaces[i].Name < this.rootNamespaces[j].Name
	})

	//
	for _, ns := range this.allNamespaces {
		for _, typ := range ns.Types {
			if typ.Class {
				for n, st := range typ.ClassDef.StaticInterfaces {
					if st.Placeholder {
						if it, ok := this.typeMap[st.FullName]; ok {
							typ.ClassDef.StaticInterfaces[n] = it
						} else {
							log.Panic("?")
						}
					} else {
						log.Panic("?")
					}
				}
			}
		}
	}

	return &Model{
		RootNamespaces: this.rootNamespaces,
		AllNamespaces:  this.allNamespaces,
	}
}

func (this *ModelParser) resolveNamespace(nsName string) *Namespace {
	parts := strings.Split(nsName, ".")
	var fullName string
	var parentNs *Namespace
	for n, part := range parts {
		if n > 0 {
			fullName += "."
		}
		fullName += part
		ns, ok := this.nsMap[fullName]
		if !ok {
			ns = &Namespace{Name: part, FullName: fullName}
			this.nsMap[fullName] = ns
			this.allNamespaces = append(this.allNamespaces, ns)
			if parentNs != nil {
				parentNs.Children = append(parentNs.Children, ns)
				ns.Parent = parentNs
			} else {
				this.rootNamespaces = append(this.rootNamespaces, ns)
			}
		}
		parentNs = ns
	}
	return parentNs
}

func (this *ModelParser) makeTypeDefFullName(typeDefRow *mdmodel.TypeDefRow) string {
	enclosingTypeDef := this.mdModel.Tables.NestedClass.GetEnclosingType(typeDefRow)
	if enclosingTypeDef != nil {
		return this.makeTypeDefFullName(enclosingTypeDef) + "." + typeDefRow.TypeName
	}
	if typeDefRow.TypeNamespace != "" {
		return typeDefRow.TypeNamespace + "." + typeDefRow.TypeName
	} else {
		return typeDefRow.TypeName
	}
}

func (this *ModelParser) parseTypeDef(typeDefRow *mdmodel.TypeDefRow) *Type {
	typ := &Type{}
	typ.Name = typeDefRow.TypeName

	customAttributeRows := this.mdModel.Tables.CustomAttribute.ListByTypeDef(typeDefRow)
	typ.Attributes = this.ParseAttributes(customAttributeRows)

	attrMap := typ.GetAttributeMap()

	const archAttrName = "Windows.Win32.Interop.SupportedArchitectureAttribute"
	if archAttr, ok := attrMap[archAttrName]; ok {
		arch := archAttr.Args[0].(int32)
		if arch&0x2 == 0 {
			return nil
		}
	}

	//
	typ.FullName = this.makeTypeDefFullName(typeDefRow)

	flags := typeDefRow.Flags
	if flags&mdmodel.TypeAttributes.Interface != 0 {
		typ.Kind = TypeInterface
		typ.Interface = true
		typ.InterfaceDef = this.parseInterfaceDef(typeDefRow)
	} else if typeDefRow.Extends != nil {
		extendsType := typeDefRow.Extends.(mdmodel.TypeRow).GetFullTypeName()
		if extendsType == "System.Enum" {
			typ.Kind = TypeEnum
			typ.Enum = true
			typ.EnumDef = this.parseEnumDef(typ, typeDefRow)
		} else if extendsType == "System.Delegate" ||
			extendsType == "System.MulticastDelegate" {
			typ.Kind = TypeFunction
			typ.Func = true
			typ.FuncDef = this.parseFuncDef(typ, typeDefRow)
		} else if typeDefRow.DerivesFrom("System.ValueType") {
			if typeDefRow.Flags&mdmodel.TypeAttributes.ExplicitLayout != 0 {
				typ.Kind = TypeUnion
				typ.Union = true
				typ.UnionDef = this.parseUnionDef(typ, typeDefRow)
			} else {
				typ.Kind = TypeStruct
				typ.Struct = true
				typ.StructDef = this.parseStructDef(typ, typeDefRow)

				if _, ok := attrMap["Windows.Win32.Interop.NativeTypedefAttribute"]; ok {
					typ.Struct = false
					typ.Kind = TypeAlias
					typ.Alias = true
					typ.AliasType = typ.StructDef.Fields[0].Type
				}

			}
		} else if typeDefRow.DerivesFrom("System.Object") {
			if typ.Name == "<Module>" || typ.Name == "Apis" { //?? public static class Apis
				typ.Kind = TypePseudo
				typ.Pseudo = true //?
				typ.PseudoDef = this.parsePseudoDef(typ, typeDefRow)
			} else {
				typ.Kind = TypeClass
				typ.Class = true
				typ.ClassDef = this.parseClassDef(typ, typeDefRow)
			}
		} else if typeDefRow.DerivesFrom("System.Attribute") {
			//?
			return nil
		} else {
			typ.Kind = TypeClass
			typ.Class = true
			typ.ClassDef = this.parseClassDef(typ, typeDefRow)
		}
	} else {
		typ.PseudoDef = this.parsePseudoDef(typ, typeDefRow)
	}

	if existingType, ok := this.typeMap[typ.FullName]; ok {
		return existingType //?
	}
	if typ.Kind != TypeRef {
		this.typeMap[typ.FullName] = typ
	}
	if typ.Kind == TypeUnknown {
		//println("??")
	}

	//
	genericParamRows := this.mdModel.Tables.GenericParam.ListByTypeDef(typeDefRow)
	for _, gpr := range genericParamRows {
		typ.GenericDefParams = append(typ.GenericDefParams, gpr.Name)
	}
	typ.Generic = typ.GenericDefParams != nil

	return typ
}

func (this *ModelParser) parseInterfaceDef(row *mdmodel.TypeDefRow) *InterfaceDef {
	def := &InterfaceDef{}
	InterfaceImplTable := this.mdModel.Tables.InterfaceImpl
	interfaceImplRows := InterfaceImplTable.GetByTypeDef(row)
	implCount := len(interfaceImplRows)
	for n := 0; n < implCount; n++ {
		extendsType := this.parseTypeRow(interfaceImplRows[n].Interface)
		def.Extends = append(def.Extends, extendsType)
	}
	for _, m := range row.MethodList {
		def.Methods = append(def.Methods, this.parseMethod(m))
	}
	caRows := this.mdModel.Tables.CustomAttribute.ListByTypeDef(row)
	def.Attributes = this.ParseAttributes(caRows)
	
	//
	//?def.Import = row.Flags & mdmodel.TypeAttributes.Import != 0

	return def

}

func (this *ModelParser) parseClassDef(typ *Type, row *mdmodel.TypeDefRow) *ClassDef {
	def := &ClassDef{}

	def.Static = row.Flags&mdmodel.TypeAttributes.Abstract != 0 &&
		row.Flags&mdmodel.TypeAttributes.Sealed != 0

	InterfaceImplTable := this.mdModel.Tables.InterfaceImpl
	interfaceImplRows := InterfaceImplTable.GetByTypeDef(row)
	implCount := len(interfaceImplRows)

	//parse static interfaces
	if true {
		caRows := this.mdModel.Tables.CustomAttribute.ListByTypeDef(row)
		attrs := this.ParseAttributes(caRows)
		for _, attr := range attrs {
			if strings.Contains(attr.Type.FullName, "Windows.Foundation.Metadata.StaticAttribute") {
				staticType := this.newPlaceholderType(attr.Args[0].(string))
				def.StaticInterfaces = append(def.StaticInterfaces, staticType)
			}
		}
	}

	for n := 0; n < implCount; n++ {
		interfaceImplRow := interfaceImplRows[n]
		implementType := this.parseTypeRow(interfaceImplRow.Interface)
		def.Implements = append(def.Implements, implementType)

		//
		caRows := this.mdModel.Tables.CustomAttribute.ListByInterfaceImpl(interfaceImplRow)
		attrs := this.ParseAttributes(caRows)
		for _, attr := range attrs {
			if strings.Contains(attr.Type.FullName, "Windows.Foundation.Metadata.DefaultAttribute") {
				def.DefaultInterface = implementType
				break
			}
		}
	}

	for _, f := range row.FieldList {
		if f.Flags&mdmodel.FieldAttributes.InitOnly != 0 &&
			f.Flags&mdmodel.FieldAttributes.Static != 0 ||
			f.Flags&mdmodel.FieldAttributes.Literal != 0 {
			def.Constants = append(def.Constants, this.parseConstant(f)) //?
		} else {
			def.Fields = append(def.Fields, this.parseField(f))
		}
	}
	for _, m := range row.MethodList {
		def.Methods = append(def.Methods, this.parseMethod(m))
	}
	return def
}

func (this *ModelParser) parseEnumDef(typ *Type, row *mdmodel.TypeDefRow) *EnumDef {
	def := &EnumDef{}
	for _, f := range row.FieldList {
		if f.Flags&mdmodel.FieldAttributes.Static != 0 ||
			f.Flags&mdmodel.FieldAttributes.Literal != 0 {
			def.Values = append(def.Values, this.parseConstant(f))
		} else if f.Name == "value__" {
			def.BaseType = this.parseSigType(f.Signature.Type)
		} else {
			//println("?1")
		}
	}
	//
	cas := this.mdModel.Tables.CustomAttribute.ListByTypeDef(row)
	attrs := this.ParseAttributes(cas)
	for _, attr := range attrs {
		if attr.Type.FullName == "System.FlagsAttribute" {
			def.Flags = true
			break
		}
	}
	return def
}

func (this *ModelParser) parseFuncDef(typ *Type, row *mdmodel.TypeDefRow) *FuncDef {
	fd := &FuncDef{}
	fd.Name = row.TypeName
	for _, mdMethod := range row.MethodList {
		if mdMethod.Name == "Invoke" {
			m := this.parseMethod(mdMethod)
			fd.Params = m.Params
			fd.ReturnType = m.ReturnType
			break
		}
	}

	caRows := this.mdModel.Tables.CustomAttribute.ListByTypeDef(row)
	fd.Attributes = this.ParseAttributes(caRows)
	return fd
}

func (this *ModelParser) parseStructDef(typ *Type, row *mdmodel.TypeDefRow) *StructDef {
	def := &StructDef{}
	for _, f := range row.FieldList {
		if f.Flags&mdmodel.FieldAttributes.InitOnly != 0 &&
			f.Flags&mdmodel.FieldAttributes.Static != 0 ||
			f.Flags&mdmodel.FieldAttributes.Literal != 0 {
			def.Constants = append(def.Constants, this.parseConstant(f))
		} else {
			def.Fields = append(def.Fields, this.parseField(f))
		}
	}
	return def
}

func (this *ModelParser) parseUnionDef(typ *Type, row *mdmodel.TypeDefRow) *UnionDef {
	def := &UnionDef{}
	for _, f := range row.FieldList {
		if f.Flags&mdmodel.FieldAttributes.InitOnly != 0 &&
			f.Flags&mdmodel.FieldAttributes.Static != 0 ||
			f.Flags&mdmodel.FieldAttributes.Literal != 0 {
			def.Constants = append(def.Constants, this.parseConstant(f))
		} else {
			def.Fields = append(def.Fields, this.parseUnionField(f))
		}
	}
	return def
}

func (this *ModelParser) parsePseudoDef(typ *Type, row *mdmodel.TypeDefRow) *PseudoDef {
	def := &PseudoDef{}
	for _, f := range row.FieldList {
		if f.Flags&mdmodel.FieldAttributes.InitOnly != 0 &&
			f.Flags&mdmodel.FieldAttributes.Static != 0 ||
			f.Flags&mdmodel.FieldAttributes.Literal != 0 {
			def.Constants = append(def.Constants, this.parseConstant(f))
		} else {
			def.Fields = append(def.Fields, this.parseField(f))
		}
	}
	for _, m := range row.MethodList {
		method := this.parseMethod(m)
		if method == nil { //
			continue
		}
		def.Methods = append(def.Methods, method)
	}
	return def
}

func (this *ModelParser) parseSigType(sigType *mdmodel.Type) *Type {
	var typ *Type
	if sigType.Primitive {
		var t Type
		t.Kind = TypePrimitive
		t.Primitive = true
		types := mdmodel.ElementTypes
		switch sigType.Kind {
		case types.I1:
			t.Name = "int8"
			t.Size = 1
		case types.U1:
			t.Name = "byte"
			t.Unsigned = true
			t.Size = 1
		case types.I2:
			t.Name = "int16"
			t.Size = 2
		case types.U2:
			t.Name = "uint16"
			t.Unsigned = true
			t.Size = 2
		case types.I4:
			t.Name = "int32"
			t.Size = 4
		case types.U4:
			t.Name = "uint32"
			t.Unsigned = true
			t.Size = 4
		case types.I8:
			t.Name = "int64"
			t.Size = 8
		case types.U8:
			t.Name = "uint64"
			t.Unsigned = true
			t.Size = 8
		case types.R4:
			t.Name = "float32"
			t.Size = 4
		case types.R8:
			t.Name = "float64"
			t.Size = 8
		case types.I:
			t.Name = "uintptr" //int?
			t.Unsigned = true
			t.Size = int(unsafe.Sizeof(uintptr(0)))
		case types.U:
			t.Name = "uintptr"
			t.Unsigned = true
			t.Size = int(unsafe.Sizeof(uintptr(0)))
		case types.Char:
			t.Name = "uint16" //xx
			t.Unsigned = true
			t.Size = 2
		case types.Boolean:
			t.Name = "bool"
			t.Size = 1
		default:
			log.Panic("??")
		}
		typ = &t
	} else if sigType.IsString {
		typ = &Type{
			Kind: TypeString,
			Name: "string",
		}
	} else if sigType.ValueType_TypeDefRefSpec != nil {
		anyTypeRow := sigType.ValueType_TypeDefRefSpec
		typ = this.parseTypeRow(anyTypeRow)
	} else if sigType.Class_TypeDefRefSpec != nil {
		typ = this.parseTypeRow(sigType.Class_TypeDefRefSpec)
	} else if sigType.Ptr_Type != nil {
		typ = this.parsePtrType(sigType.Ptr_Type)
	} else if sigType.Void {
		typ = this.voidType()
	} else if sigType.Array_Type != nil {
		typ = this.parseArray(sigType)
	} else if sigType.GenericInst != nil {
		typ = this.ParseSigGenericInst(sigType.GenericInst)
	} else if sigType.Any {
		typ = &Type{
			Kind: TypeAny,
			Any:  true,
			Name: "interface{}",
		} //?
	} else if sigType.SzArray_Type != nil {
		typ = this.parseSzArray(sigType)
	} else if sigType.Kind == mdmodel.ElementTypes.Var {
		typ = &Type{
			Kind:              TypeGenericParam,
			GenericParam:      true,
			GenericParamIndex: sigType.Var,
			Name:              "`" + strconv.Itoa(int(sigType.Var)+1),
		}
	} else {
		log.Panic("not implemented")
	}
	if typ.FullName == "" {
		typ.FullName = typ.Name
	}
	if typ.FullName == "" {
		//void?
		if !typ.Void {
			log.Panic("?")
		}
	} else if existingType, ok := this.typeMap[typ.FullName]; ok {
		typ = existingType //?
	}

	//
	if replaceType, ok := this.replaceTypeMap[typ.FullName]; ok {
		return replaceType
	}

	//
	if typ.Kind != TypeRef {
		this.typeMap[typ.FullName] = typ
	}
	if typ.Kind == TypeUnknown {
		//println("?")
	}
	return typ
}

func (this *ModelParser) parseTypeRow(anyTypeRow interface{}) *Type {
	var typ *Type
	switch v := anyTypeRow.(type) {
	case *mdmodel.TypeDefRow:
		typ = this.parseTypeDef(v)
	case *mdmodel.TypeRefRow:
		typ = this.parseTypeRef(v)
	case *mdmodel.TypeSpecRow:
		typ = this.parseTypeSpec(v)
	default:
		panic("?")
	}
	return typ
}

func (this *ModelParser) parseFieldType(fieldRow *mdmodel.FieldRow) *Type {
	return this.parseSigType(fieldRow.Signature.Type)
}

func (this *ModelParser) parseConstValue(constantRow *mdmodel.ConstantRow) interface{} {
	types := mdmodel.ElementTypes
	valueBytes := constantRow.Value

	var pValue unsafe.Pointer
	if len(valueBytes) != 0 {
		pValue = unsafe.Pointer(&valueBytes[0])
	}
	var value interface{}
	switch constantRow.Type {
	case types.I1:
		value = int8(valueBytes[0])
	case types.U1:
		value = valueBytes[0]
	case types.I2:
		value = *(*int16)(pValue)
	case types.U2:
		value = *(*uint16)(pValue)
	case types.I4:
		value = *(*int32)(pValue)
	case types.U4:
		value = *(*uint32)(pValue)
	case types.I8:
		value = *(*int64)(pValue)
	case types.U8:
		value = *(*uint64)(pValue)
	case types.R4:
		value = *(*float32)(pValue)
	case types.R8:
		value = *(*float64)(pValue)
	case types.String:
		ws := unsafe.Slice((*uint16)(pValue), len(valueBytes)/2)
		value = syscall.UTF16ToString(ws)
	case types.I:
		value = *(*int)(pValue) //?
	case types.U:
		value = *(*uintptr)(pValue) //?
	default:
		//println("??")
	}
	return value
}

func (this *ModelParser) parseConstant(fieldRow *mdmodel.FieldRow) *Constant {
	constantRow := this.mdModel.Tables.Constant.GetByField(fieldRow)
	c := &Constant{
		Name:  fieldRow.Name,
		Type:  this.parseFieldType(fieldRow),
		Value: this.parseConstValue(constantRow),
	}
	return c
}

func (this *ModelParser) parseTypeRef(typeRefRow *mdmodel.TypeRefRow) *Type {
	var refScope string
	switch v := typeRefRow.ResolutionScope.(type) {
	case *mdmodel.ModuleRow:
		break
	case *mdmodel.ModuleRefRow:
		log.Panic("?")
	case *mdmodel.AssemblyRefRow:
		refScope = "assembly:" + v.Name
	case *mdmodel.TypeRefRow:
		parentType := this.parseTypeRef(v)
		refScope = "type:" + parentType.FullName //?
	}
	ns := &Namespace{}
	ns.FullName = typeRefRow.TypeNamespace
	if len(ns.FullName) > 0 {
		pos := strings.LastIndexByte(ns.FullName, '.')
		ns.Name = ns.FullName[pos+1:]
	}
	typ := &Type{
		Kind:      TypeRef,
		RefScope:  refScope,
		Name:      typeRefRow.TypeName,
		Namespace: ns,
	}
	if strings.HasPrefix(typ.RefScope, "type:") {
		typ.FullName = typ.RefScope[5:] + "." + typ.Name
	} else {
		typ.FullName = typeRefRow.GetFullTypeName()
	}
	if existingType, ok := this.typeMap[typ.FullName]; ok {
		return existingType //?
	}
	if typ.Kind != TypeRef {
		this.typeMap[typ.FullName] = typ
	}
	return typ
}

func (this *ModelParser) parseTypeSpec(typeSpecRow *mdmodel.TypeSpecRow) *Type {
	var typ *Type
	typeSpec := typeSpecRow.Signature
	if typeSpec.Ptr_Type != nil {
		toType := this.parseSigType(typeSpec.Ptr_Type)
		typ = &Type{
			Kind:      TypePointer,
			Pointer:   true,
			PointerTo: toType,
			Name:      "*" + toType.FullName,
		}
	} else if typeSpec.Array_Type != nil {
		arrType := this.parseSigType(typeSpec.Array_Type)
		typ = &Type{
			Kind:     TypeArray,
			Array:    true,
			ArrayDef: &ArrayDef{ElementType: arrType},
			Name:     "[]" + arrType.FullName,
		}
	} else if typeSpec.SzArray_Type != nil {
		arrType := this.parseSigType(typeSpec.SzArray_Type)
		typ = &Type{
			Kind:     TypeArray,
			Array:    true,
			ArrayDef: &ArrayDef{ElementType: arrType},
			Name:     "[]" + arrType.FullName,
		}
	} else if typeSpec.GenericInst != nil {
		typ = this.ParseSigGenericInst(typeSpec.GenericInst)
	} else {
		log.Panic("?")
	}
	if typ.FullName == "" {
		typ.FullName = typ.Name //?
	}
	if existingType, ok := this.typeMap[typ.FullName]; ok {
		return existingType //?
	}
	if typ.Kind != TypeRef {
		this.typeMap[typ.FullName] = typ
	}
	return typ
}

func (this *ModelParser) parseField(fieldRow *mdmodel.FieldRow) *Field {
	f := &Field{}
	if fieldRow.Flags&mdmodel.FieldAttributes.Static != 0 {
		f.Static = true //?
	}
	f.Name = fieldRow.Name
	f.Type = this.parseFieldType(fieldRow)
	//?f.Value
	caRows := this.mdModel.Tables.CustomAttribute.ListByField(fieldRow)
	f.Attributes = this.ParseAttributes(caRows)
	return f
}

func (this *ModelParser) parseUnionField(fieldRow *mdmodel.FieldRow) *Field {
	f := &Field{}
	if fieldRow.Flags&mdmodel.FieldAttributes.Static != 0 {
		f.Static = true //?
	}
	f.Name = fieldRow.Name
	f.Type = this.parseFieldType(fieldRow)
	fieldLayout := this.mdModel.Tables.FieldLayout.GetByField(fieldRow)
	if fieldLayout.Offset != 0 {
		//println("??")
	}
	return f
}

func (this *ModelParser) parseMethod(methodDefRow *mdmodel.MethodDefRow) *Method {

	const archAttrName = "Windows.Win32.Interop.SupportedArchitectureAttribute"
	const overloadAttrName = "Windows.Foundation.Metadata.OverloadAttribute"
	attrRows := this.mdModel.Tables.CustomAttribute.ListByMethodDef(methodDefRow)
	attrs := this.ParseAttributes(attrRows)

	m := &Method{}
	for _, attr := range attrs {
		if attr.Type.FullName == archAttrName {
			if attr.Args[0].(int32)&0x02 == 0 {
				return nil
			}
		} else if attr.Type.FullName == overloadAttrName {
			m.OverloadName = attr.Args[0].(string)
		}
	}

	m.Name = methodDefRow.Name

	m.Static = methodDefRow.Flags&mdmodel.MethodAttributes.Static != 0
	m.SysCall = methodDefRow.Flags&mdmodel.MethodAttributes.PInvokeImpl != 0
	if m.SysCall {
		implMapTable := this.mdModel.Tables.ImplMap
		implMap := implMapTable.MemberForwardImplRowMap[methodDefRow.RowIndex()]
		//
		m.SysCallName = implMap.ImportName
		m.SysCallDll = implMap.ImportScope.Name
		m.SysCallSetLastError = implMap.MappingFlags&
			mdmodel.PInvokeAttributes.SupportsLastError != 0
		//
		caRows := this.mdModel.Tables.CustomAttribute.ListByMethodDef(methodDefRow)
		m.Attributes = this.ParseAttributes(caRows)
		attrMap := make(map[string]*Attribute)
		for _, a := range m.Attributes {
			attrMap[a.Type.Name] = a
		}
		if attr, ok := attrMap["Windows.Win32.Interop.SupportedOSPlatformAttribute"]; ok {
			m.SupportedOS = attr.Args[0].(string)
		}
	}
	var retTypeResolved bool
	for _, p := range methodDefRow.ParamList {
		if p.Sequence == 0 {
			m.ReturnType = this.parseReturnType(methodDefRow.Signature, p) //?
			retTypeResolved = true
			continue
		}
		m.Params = append(m.Params, this.parseParam(methodDefRow.Signature, p))
	}
	if !retTypeResolved {
		m.ReturnType = this.parseReturnType(methodDefRow.Signature, nil) //?
	}

	//
	genericParamRows := this.mdModel.Tables.GenericParam.ListByMethodDef(methodDefRow)
	for _, gpr := range genericParamRows {
		m.GenericParams = append(m.GenericParams, gpr.Name)
	}
	m.GenericParamCount = len(m.GenericParams)
	m.Generic = m.GenericParamCount != 0

	return m
}

func (this *ModelParser) parseParam(methodDefSig *mdmodel.MethodDefSig,
	paramRow *mdmodel.ParamRow) *Param {
	p := &Param{}
	p.Name = paramRow.Name
	sigParam := methodDefSig.Params[paramRow.Sequence-1]
	p.Type = this.parseSigType(sigParam.Type)

	p.In = paramRow.Flags&mdmodel.ParamAttributes.In != 0
	p.Out = paramRow.Flags&mdmodel.ParamAttributes.Out != 0
	p.Optional = paramRow.Flags&mdmodel.ParamAttributes.Optional != 0
	//default?

	return p
}

func (this *ModelParser) parseReturnType(methodDefSig *mdmodel.MethodDefSig,
	paramRow *mdmodel.ParamRow) *Type {
	return this.parseSigType(methodDefSig.RetType.Param.Type)
}

func (this *ModelParser) parsePtrType(targetSigType *mdmodel.Type) *Type {
	t := &Type{}
	t.Kind = TypePointer
	t.Pointer = true
	t.PointerTo = this.parseSigType(targetSigType)
	t.Name = "*" + t.PointerTo.Name
	t.FullName = "*" + t.PointerTo.FullName
	return t
}

func (this *ModelParser) voidType() *Type {
	return &Type{
		Kind: TypeVoid,
		Void: true,
	}
}

func (this *ModelParser) parseArray(arrayType *mdmodel.Type) *Type {
	typ := &Type{}
	typ.Kind = TypeArray
	typ.Array = true
	elemType := this.parseSigType(arrayType.Array_Type)
	var dimSizes []uint32
	if arrayType.Array_Shape != nil {
		dimSizes = arrayType.Array_Shape.Sizes
	}
	typ.ArrayDef = &ArrayDef{
		ElementType: elemType,
		DimSizes:    dimSizes,
	}
	if len(dimSizes) == 1 {
		typ.Name = fmt.Sprintf("[%d]", dimSizes[0]) + elemType.FullName
	} else {
		typ.Name = "[]" + elemType.FullName
	}

	return typ
}

func (this *ModelParser) parseSzArray(arrayType *mdmodel.Type) *Type {
	elemType := this.parseSigType(arrayType.SzArray_Type)
	typ := &Type{
		Kind:  TypeArray,
		Array: true,
		ArrayDef: &ArrayDef{
			ElementType: elemType,
		},
		Name: "[]" + elemType.FullName,
	}
	return typ
}

func (this *ModelParser) ParseSigGenericInst(genericInst *mdmodel.GenericInst) *Type {

	typ := &Type{
		GenericInst: true,
	}
	genType := this.parseTypeRow(genericInst.TypeDefRefSpec)
	typ.GenericType = genType
	var genArgTypes []*Type
	for _, t := range genericInst.GenArgTypes {
		genArgTypes = append(genArgTypes, this.parseSigType(t))
	}
	typ.GenericArgTypes = genArgTypes

	typ.FullName, typ.Name = this.makeGenericInstName(genType, genArgTypes)
	if genericInst.Class {
		typ.Class = true //?
	} else {
		typ.Struct = true //?
	}
	if genType.Kind == TypeRef {
		if typ.Name[0] == 'I' {
			typ.Kind = TypeInterface
		} else {
			typ.Kind = TypeClass //?
		}
	} else {
		typ.Kind = genType.Kind
		switch typ.Kind {
		case TypeFunction:
			typ.Func = true
			typ.FuncDef = genType.FuncDef
		case TypeStruct:
			typ.Struct = true
			typ.StructDef = genType.StructDef
		case TypeUnion:
			typ.Union = true
			typ.UnionDef = genType.UnionDef
		case TypeArray:
			typ.Array = true
			typ.ArrayDef = genType.ArrayDef
		case TypeInterface:
			typ.Interface = true
			typ.InterfaceDef = genType.InterfaceDef
		case TypeClass:
			typ.Class = true
			typ.ClassDef = genType.ClassDef
		default:
			log.Panic("?")
		}
	}
	return typ
}

func (this *ModelParser) ParseAttributes(caRows []*mdmodel.CustomAttributeRow) []*Attribute {
	var attrs []*Attribute
	for _, caRow := range caRows {
		attr := &Attribute{}
		switch v := caRow.Type.(type) {
		case *mdmodel.MethodDefRow:
			attr.Type = this.parseTypeDef(v.OwnerType)
		case *mdmodel.MemberRefRow:
			switch v2 := v.Class.(type) {
			case *mdmodel.TypeRefRow:
				attr.Type = this.parseTypeRef(v2)
			case *mdmodel.TypeSpecRow:
				attr.Type = this.parseTypeSpec(v2)
			}
		}
		for _, fa := range caRow.ValueSig.FixedArgs {
			attr.Args = append(attr.Args, fa.ToInterface())
		}
		if len(caRow.ValueSig.NamedArgs) > 0 {
			attr.NamedArgs = make(map[string]interface{})
		}
		for _, na := range caRow.ValueSig.NamedArgs {
			attr.NamedArgs[na.Name] = na.ToInterface()
		}
		//
		attrs = append(attrs, attr)
	}
	return attrs
}

func (this *ModelParser) makeGenericInstName(genType *Type, genArgTypes []*Type) (string, string) {
	name := genType.Name
	pos := strings.IndexByte(name, '`')
	name = name[:pos]
	name += "["
	for n, gat := range genArgTypes {
		if n > 0 {
			name += ", "
		}
		name += gat.FullName
	}
	name += "]"
	fullName := genType.FullName[:len(genType.FullName)-len(genType.Name)] + name
	return fullName, name
}

func (this *ModelParser) _getRootEnclosingType(row *mdmodel.TypeDefRow) *mdmodel.TypeDefRow {
	enclosingType := this.mdModel.Tables.NestedClass.GetEnclosingType(row)
	if enclosingType == nil {
		return row
	}
	return this._getRootEnclosingType(enclosingType)
}

func (this *ModelParser) getRootEnclosingType(row *mdmodel.TypeDefRow) *mdmodel.TypeDefRow {
	rootEnclosingType := this._getRootEnclosingType(row)
	if rootEnclosingType == row {
		return nil
	}
	return rootEnclosingType
}

func (this *ModelParser) newPlaceholderType(fullName string) *Type {
	return &Type{
		Kind:        TypeUnknown,
		FullName:    fullName,
		Placeholder: true,
	}
}
