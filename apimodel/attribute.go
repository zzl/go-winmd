package apimodel

type Attribute struct {
	Type      *Type
	Args      []interface{}
	NamedArgs map[string]interface{}
}
