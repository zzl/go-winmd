package apimodel

type Field struct {
	Static bool
	Name   string
	Type   *Type

	Value interface{}

	FieldOffset uint32 //?

	Attributes []*Attribute
}
