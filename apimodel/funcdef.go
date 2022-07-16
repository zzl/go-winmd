package apimodel

type FuncDef struct {
	Name       string
	Params     []*Param
	ReturnType *Type

	Attributes []*Attribute
}
