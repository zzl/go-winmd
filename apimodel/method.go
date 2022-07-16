package apimodel

type Method struct {
	Static  bool
	SysCall bool

	Name                string
	SysCallName         string
	SysCallDll          string
	SysCallSetLastError bool
	SupportedOS         string

	Params     []*Param
	ReturnType *Type

	Generic           bool
	GenericParamCount int

	GenericParams []string

	//
	Attributes []*Attribute

	//
	OverloadName string
}
