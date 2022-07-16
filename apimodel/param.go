package apimodel

type Param struct {
	Name     string
	Type     *Type
	In       bool
	Out      bool
	Optional bool
}
