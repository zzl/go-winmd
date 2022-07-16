package apimodel

type InterfaceDef struct {
	Attributes []*Attribute
	Extends    []*Type
	Methods    []*Method

	Import bool
}
