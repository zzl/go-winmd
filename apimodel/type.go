package apimodel

import "fmt"

type TypeKind byte

const (
	TypeUnknown      TypeKind = 0
	TypeAny          TypeKind = 1
	TypePrimitive    TypeKind = 2
	TypePointer      TypeKind = 3
	TypeEnum         TypeKind = 4
	TypeStruct       TypeKind = 5
	TypeUnion        TypeKind = 6
	TypeInterface    TypeKind = 7
	TypeClass        TypeKind = 8
	TypeAlias        TypeKind = 9
	TypeVoid         TypeKind = 10
	TypeRef          TypeKind = 11
	TypeString       TypeKind = 12
	TypeFunction     TypeKind = 13
	TypeArray        TypeKind = 14
	TypeGenericParam TypeKind = 15
	TypePseudo       TypeKind = 16
)

type Type struct {
	Namespace *Namespace
	Kind      TypeKind

	EnclosingType *Type
	NestedTypes   []*Type

	Primitive bool //i1-8/ui1-8,r4,r8,i/u,bool
	Unsigned  bool
	Size      int

	//
	Pointer   bool
	PointerTo *Type

	Array    bool
	ArrayDef *ArrayDef

	Enum    bool
	EnumDef *EnumDef

	Struct    bool
	StructDef *StructDef

	Union    bool
	UnionDef *UnionDef

	Interface    bool
	InterfaceDef *InterfaceDef

	Class    bool
	ClassDef *ClassDef

	Func    bool
	FuncDef *FuncDef

	Alias     bool
	AliasType *Type

	Pseudo    bool
	PseudoDef *PseudoDef

	Void bool
	Any  bool

	RefScope string

	Name     string
	FullName string

	//
	Generic          bool
	GenericDefParams []string

	//
	GenericParam      bool //enclosing type generic param
	GenericParamIndex uint32

	//
	GenericInst     bool
	GenericType     *Type
	GenericArgTypes []*Type

	//
	Attributes []*Attribute

	SiezInfo *SizeInfo

	Placeholder bool
}

func (this *Type) String() string {
	if this.Pointer {
		return "*" + this.PointerTo.String()
	}
	if this.Array {
		if len(this.ArrayDef.DimSizes) == 1 {
			return fmt.Sprintf("[%d]", this.ArrayDef.DimSizes[0]) +
				this.ArrayDef.ElementType.String()
		} else {
			return "[]" + this.ArrayDef.ElementType.String() //?
		}
	}
	if this.Class || this.Interface {
		return "*" + this.FullName
	}
	return this.FullName
}

func (this *Type) GetAttributeMap() map[string]*Attribute {
	m := make(map[string]*Attribute)
	for _, a := range this.Attributes {
		m[a.Type.FullName] = a
	}
	return m
}

func (this *Type) HasAttribute(name string) bool {
	for _, a := range this.Attributes {
		if a.Type.FullName == name {
			return true
		}
	}
	return false
}
