package apimodel

type EnumDef struct {
	BaseType *Type
	Flags    bool
	Values   []*Constant
}
