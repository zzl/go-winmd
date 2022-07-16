package apimodel

type ClassDef struct {
	Static           bool
	Extends          *Type
	Implements       []*Type
	DefaultInterface *Type
	StaticInterfaces []*Type
	Fields           []*Field
	Methods          []*Method
	Constants        []*Constant
}
