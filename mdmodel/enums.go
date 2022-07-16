package mdmodel

import "reflect"

type AssemblyHashAlgorithmEnum uint32

var AssemblyHashAlgorithm = struct {
	None     AssemblyHashAlgorithmEnum
	Reserved AssemblyHashAlgorithmEnum
	SHA1     AssemblyHashAlgorithmEnum
}{
	None:     0,
	Reserved: 0x8003,
	SHA1:     0x8004,
}

type AssemblyFlagsEnum uint32

var AssemblyFlags = struct {
	PublicKey                  AssemblyFlagsEnum
	Retargetable               AssemblyFlagsEnum
	DisableJITcompileOptimizer AssemblyFlagsEnum
	EnableJITcompileTracking   AssemblyFlagsEnum
}{
	PublicKey:                  0x0001,
	Retargetable:               0x0100,
	DisableJITcompileOptimizer: 0x4000,
	EnableJITcompileTracking:   0x8000,
}

type EventAttributesEnum uint16

var EventAttributes = struct {
	SpecialName   EventAttributesEnum
	RTSpecialName EventAttributesEnum
}{
	SpecialName:   0x0200,
	RTSpecialName: 0x0400,
}

type FieldAttributesEnum uint16

func (me FieldAttributesEnum) String() string {
	return flagsEnumToString(FieldAttributes, uint16(me))
}

var FieldAttributes = struct {
	FieldAccessMask    FieldAttributesEnum
	CompilerControlled FieldAttributesEnum
	Private            FieldAttributesEnum
	FamANDAssem        FieldAttributesEnum
	Assembly           FieldAttributesEnum
	Family             FieldAttributesEnum
	FamORAssem         FieldAttributesEnum
	Public             FieldAttributesEnum
	Static             FieldAttributesEnum
	InitOnly           FieldAttributesEnum
	Literal            FieldAttributesEnum
	NotSerialized      FieldAttributesEnum
	SpecialName        FieldAttributesEnum
	PInvokeImpl        FieldAttributesEnum
	RTSpecialName      FieldAttributesEnum
	HasFieldMarshal    FieldAttributesEnum
	HasDefault         FieldAttributesEnum
	HasFieldRVA        FieldAttributesEnum
}{
	FieldAccessMask:    0x0007,
	CompilerControlled: 0x0000,
	Private:            0x0001,
	FamANDAssem:        0x0002,
	Assembly:           0x0003,
	Family:             0x0004,
	FamORAssem:         0x0005,
	Public:             0x0006,
	Static:             0x0010,
	InitOnly:           0x0020,
	Literal:            0x0040,
	NotSerialized:      0x0080,
	SpecialName:        0x0200,
	PInvokeImpl:        0x2000,
	RTSpecialName:      0x0400,
	HasFieldMarshal:    0x1000,
	HasDefault:         0x8000,
	HasFieldRVA:        0x0100,
}

type FileAttributesEnum uint32

var FileAttributes = struct {
	ContainsMetaData   FileAttributesEnum
	ContainsNoMetaData FileAttributesEnum
}{
	ContainsMetaData:   0x0000,
	ContainsNoMetaData: 0x0001,
}

type GenericParamAttributesEnum uint16

var GenericParamAttributes = struct {
	VarianceMask                   GenericParamAttributesEnum
	None                           GenericParamAttributesEnum
	Covariant                      GenericParamAttributesEnum
	Contravariant                  GenericParamAttributesEnum
	SpecialConstraintMask          GenericParamAttributesEnum
	ReferenceTypeConstraint        GenericParamAttributesEnum
	NotNullableValueTypeConstraint GenericParamAttributesEnum
	DefaultConstructorConstraint   GenericParamAttributesEnum
}{
	VarianceMask:                   0x0003,
	None:                           0x0000,
	Covariant:                      0x0001,
	Contravariant:                  0x0002,
	SpecialConstraintMask:          0x001C,
	ReferenceTypeConstraint:        0x0004,
	NotNullableValueTypeConstraint: 0x0008,
	DefaultConstructorConstraint:   0x0010,
}

type PInvokeAttributesEnum uint16

var PInvokeAttributes = struct {
	NoMangle            PInvokeAttributesEnum
	CharSetMask         PInvokeAttributesEnum
	CharSetNotSpec      PInvokeAttributesEnum
	CharSetAnsi         PInvokeAttributesEnum
	CharSetUnicode      PInvokeAttributesEnum
	CharSetAuto         PInvokeAttributesEnum
	SupportsLastError   PInvokeAttributesEnum
	CallConvMask        PInvokeAttributesEnum
	CallConvPlatformapi PInvokeAttributesEnum
	CallConvCdecl       PInvokeAttributesEnum
	CallConvStdcall     PInvokeAttributesEnum
	CallConvThiscall    PInvokeAttributesEnum
	CallConvFastcall    PInvokeAttributesEnum
}{
	NoMangle:            0x0001,
	CharSetMask:         0x0006,
	CharSetNotSpec:      0x0000,
	CharSetAnsi:         0x0002,
	CharSetUnicode:      0x0004,
	CharSetAuto:         0x0006,
	SupportsLastError:   0x0040,
	CallConvMask:        0x0700,
	CallConvPlatformapi: 0x0100,
	CallConvCdecl:       0x0200,
	CallConvStdcall:     0x0300,
	CallConvThiscall:    0x0400,
	CallConvFastcall:    0x0500,
}

type ManifestResourceAttributesEnum uint32

var ManifestResourceAttributes = struct {
	VisibilityMask ManifestResourceAttributesEnum
	Public         ManifestResourceAttributesEnum
	Private        ManifestResourceAttributesEnum
}{
	VisibilityMask: 0x0007,
	Public:         0x0001,
	Private:        0x0002,
}

type MethodAttributesEnum uint16

func (me MethodAttributesEnum) String() string {
	return flagsEnumToString(MethodAttributes, uint16(me))
}

var MethodAttributes = struct {
	MemberAccessMask   MethodAttributesEnum
	CompilerControlled MethodAttributesEnum
	Private            MethodAttributesEnum
	FamANDAssem        MethodAttributesEnum
	Assem              MethodAttributesEnum
	Family             MethodAttributesEnum
	FamORAssem         MethodAttributesEnum
	Public             MethodAttributesEnum
	Static             MethodAttributesEnum
	Final              MethodAttributesEnum
	Virtual            MethodAttributesEnum
	HideBySig          MethodAttributesEnum
	VtableLayoutMask   MethodAttributesEnum
	ReuseSlot          MethodAttributesEnum
	NewSlot            MethodAttributesEnum
	Strict             MethodAttributesEnum
	Abstract           MethodAttributesEnum
	SpecialName        MethodAttributesEnum
	PInvokeImpl        MethodAttributesEnum
	UnmanagedExport    MethodAttributesEnum
	RTSpecialName      MethodAttributesEnum
	HasSecurity        MethodAttributesEnum
	RequieSecObject    MethodAttributesEnum
}{
	MemberAccessMask:   0x0007,
	CompilerControlled: 0x0000,
	Private:            0x0001,
	FamANDAssem:        0x0002,
	Assem:              0x0003,
	Family:             0x0004,
	FamORAssem:         0x0005,
	Public:             0x0006,
	Static:             0x0010,
	Final:              0x0020,
	Virtual:            0x0040,
	HideBySig:          0x0080,
	VtableLayoutMask:   0x0100,
	ReuseSlot:          0x0000,
	NewSlot:            0x0100,
	Strict:             0x0200,
	Abstract:           0x0400,
	SpecialName:        0x0800,
	PInvokeImpl:        0x2000,
	UnmanagedExport:    0x0008,
	RTSpecialName:      0x1000,
	HasSecurity:        0x4000,
	RequieSecObject:    0x8000,
}

type MethodImplAttributesEnum uint16

var MethodImplAttributes = struct {
	CodeTypeMask     MethodImplAttributesEnum
	IL               MethodImplAttributesEnum
	Native           MethodImplAttributesEnum
	OPTIL            MethodImplAttributesEnum
	Runtime          MethodImplAttributesEnum
	ManagedMask      MethodImplAttributesEnum
	Unmanaged        MethodImplAttributesEnum
	Managed          MethodImplAttributesEnum
	ForwardRef       MethodImplAttributesEnum
	PreserveSig      MethodImplAttributesEnum
	InternalCall     MethodImplAttributesEnum
	Synchronized     MethodImplAttributesEnum
	NoInlining       MethodImplAttributesEnum
	MaxMethodImplVal MethodImplAttributesEnum
	NoOptimization   MethodImplAttributesEnum
}{
	CodeTypeMask:     0x0003,
	IL:               0x0000,
	Native:           0x0001,
	OPTIL:            0x0002,
	Runtime:          0x0003,
	ManagedMask:      0x0004,
	Unmanaged:        0x0004,
	Managed:          0x0000,
	ForwardRef:       0x0010,
	PreserveSig:      0x0080,
	InternalCall:     0x1000,
	Synchronized:     0x0020,
	NoInlining:       0x0008,
	MaxMethodImplVal: 0xffff,
	NoOptimization:   0x0040,
}

type MethodSemanticsAttributesEnum uint16

var MethodSemanticsAttributes = struct {
	Setter   MethodSemanticsAttributesEnum
	Getter   MethodSemanticsAttributesEnum
	Other    MethodSemanticsAttributesEnum
	AddOn    MethodSemanticsAttributesEnum
	RemoveOn MethodSemanticsAttributesEnum
	Fire     MethodSemanticsAttributesEnum
}{
	Setter:   0x0001,
	Getter:   0x0002,
	Other:    0x0004,
	AddOn:    0x0008,
	RemoveOn: 0x0010,
	Fire:     0x0020,
}

type ParamAttributesEnum uint16

func (me ParamAttributesEnum) String() string {
	return flagsEnumToString(ParamAttributes, uint16(me))
}

var ParamAttributes = struct {
	In            ParamAttributesEnum
	Out           ParamAttributesEnum
	Optional      ParamAttributesEnum
	HasDefault    ParamAttributesEnum
	HasFieldMarsh ParamAttributesEnum
	Unused        ParamAttributesEnum
}{
	In:            0x0001,
	Out:           0x0002,
	Optional:      0x0010,
	HasDefault:    0x1000,
	HasFieldMarsh: 0x2000,
	Unused:        0xcfe0,
}

type PropertyAttributesEnum uint16

var PropertyAttributes = struct {
	SpecialName   PropertyAttributesEnum
	RTSpecialName PropertyAttributesEnum
	HasDefault    PropertyAttributesEnum
	Unused        PropertyAttributesEnum
}{
	SpecialName:   0x0200,
	RTSpecialName: 0x0400,
	HasDefault:    0x1000,
	Unused:        0xe9ff,
}

type TypeAttributesEnum uint32

func (me TypeAttributesEnum) String() string {
	return flagsEnumToString(TypeAttributes, uint32(me))
}

var TypeAttributes = struct {
	VisibilityMask         TypeAttributesEnum
	NotPublic              TypeAttributesEnum
	Public                 TypeAttributesEnum
	NestedPublic           TypeAttributesEnum
	NestedPrivate          TypeAttributesEnum
	NestedFamily           TypeAttributesEnum
	NestedAssembly         TypeAttributesEnum
	NestedFamANDAssem      TypeAttributesEnum
	NestedFamORAssem       TypeAttributesEnum
	LayoutMask             TypeAttributesEnum
	AutoLayout             TypeAttributesEnum
	SequentialLayout       TypeAttributesEnum
	ExplicitLayout         TypeAttributesEnum
	ClassSemanticsMask     TypeAttributesEnum
	Class                  TypeAttributesEnum
	Interface              TypeAttributesEnum
	Abstract               TypeAttributesEnum
	Sealed                 TypeAttributesEnum
	SpecialName            TypeAttributesEnum
	Import                 TypeAttributesEnum
	Serializable           TypeAttributesEnum
	StringFormatMask       TypeAttributesEnum
	AnsiClass              TypeAttributesEnum
	UnicodeClass           TypeAttributesEnum
	AutoClass              TypeAttributesEnum
	CustomFormatClass      TypeAttributesEnum
	CustomStringFormatMask TypeAttributesEnum
	BeforeFieldInit        TypeAttributesEnum
	RTSpecialName          TypeAttributesEnum
	HasSecurity            TypeAttributesEnum
	IsTypeForwarder        TypeAttributesEnum
}{
	VisibilityMask:         0x00000007,
	NotPublic:              0x00000000,
	Public:                 0x00000001,
	NestedPublic:           0x00000002,
	NestedPrivate:          0x00000003,
	NestedFamily:           0x00000004,
	NestedAssembly:         0x00000005,
	NestedFamANDAssem:      0x00000006,
	NestedFamORAssem:       0x00000007,
	LayoutMask:             0x00000018,
	AutoLayout:             0x00000000,
	SequentialLayout:       0x00000008,
	ExplicitLayout:         0x00000010,
	ClassSemanticsMask:     0x00000020,
	Class:                  0x00000000,
	Interface:              0x00000020,
	Abstract:               0x00000080,
	Sealed:                 0x00000100,
	SpecialName:            0x00000400,
	Import:                 0x00001000,
	Serializable:           0x00002000,
	StringFormatMask:       0x00030000,
	AnsiClass:              0x00000000,
	UnicodeClass:           0x00010000,
	AutoClass:              0x00020000,
	CustomFormatClass:      0x00030000,
	CustomStringFormatMask: 0x00c00000,
	BeforeFieldInit:        0x00100000,
	RTSpecialName:          0x00000800,
	HasSecurity:            0x00040000,
	IsTypeForwarder:        0x00200000,
}

type ElementTypeEnum uint16

func (me ElementTypeEnum) String() string {
	return enumToString(ElementTypes, uint16(me))
}

var ElementTypes = struct {
	End         ElementTypeEnum
	Void        ElementTypeEnum
	Boolean     ElementTypeEnum
	Char        ElementTypeEnum
	I1          ElementTypeEnum
	U1          ElementTypeEnum
	I2          ElementTypeEnum
	U2          ElementTypeEnum
	I4          ElementTypeEnum
	U4          ElementTypeEnum
	I8          ElementTypeEnum
	U8          ElementTypeEnum
	R4          ElementTypeEnum
	R8          ElementTypeEnum
	String      ElementTypeEnum
	Ptr         ElementTypeEnum
	ByRef       ElementTypeEnum
	ValueType   ElementTypeEnum
	Class       ElementTypeEnum
	Var         ElementTypeEnum
	Array       ElementTypeEnum
	GenericInst ElementTypeEnum
	TypedByRef  ElementTypeEnum
	I           ElementTypeEnum
	U           ElementTypeEnum
	FnPtr       ElementTypeEnum
	Object      ElementTypeEnum
	SzArray     ElementTypeEnum
	MVar        ElementTypeEnum
	CMod_Reqd   ElementTypeEnum
	CMod_Opt    ElementTypeEnum
	Internal    ElementTypeEnum
	Modifier    ElementTypeEnum
	Sentinel    ElementTypeEnum
	Pinned      ElementTypeEnum
	X50         ElementTypeEnum
	X51         ElementTypeEnum
	X52         ElementTypeEnum
	X53         ElementTypeEnum
	X54         ElementTypeEnum
	X55         ElementTypeEnum
}{
	End:         0x00,
	Void:        0x01,
	Boolean:     0x02,
	Char:        0x03,
	I1:          0x04,
	U1:          0x05,
	I2:          0x06,
	U2:          0x07,
	I4:          0x08,
	U4:          0x09,
	I8:          0x0a,
	U8:          0x0b,
	R4:          0x0c,
	R8:          0x0d,
	String:      0x0e,
	Ptr:         0x0f,
	ByRef:       0x10,
	ValueType:   0x11,
	Class:       0x12,
	Var:         0x13,
	Array:       0x14,
	GenericInst: 0x15,
	TypedByRef:  0x16,
	I:           0x18,
	U:           0x19,
	FnPtr:       0x1b,
	Object:      0x1c,
	SzArray:     0x1d,
	MVar:        0x1e,
	CMod_Reqd:   0x1f,
	CMod_Opt:    0x20,
	Internal:    0x21,
	Modifier:    0x40,
	Sentinel:    0x41,
	Pinned:      0x45,
	X50:         0x50,
	X51:         0x51,
	X52:         0x52,
	X53:         0x53,
	X54:         0x54,
	X55:         0x55,
}

//
func flagsEnumToString[TBase uint16 | uint32](enumValuesStruct any, value TBase) string {
	var s string
	structValue := reflect.ValueOf(enumValuesStruct)
	structType := structValue.Type()
	for n := 0; n < structValue.NumField(); n++ {
		field := structValue.Field(n)
		if TBase(field.Uint())&value != 0 {
			if s != "" {
				s += "|"
			}
			s += structType.Field(n).Name
		}
	}
	return s
}

func enumToString[TBase uint16 | uint32](enumValuesStruct any, value TBase) string {
	structValue := reflect.ValueOf(enumValuesStruct)
	structType := structValue.Type()
	for n := 0; n < structValue.NumField(); n++ {
		field := structValue.Field(n)
		if TBase(field.Uint()) == value {
			return structType.Field(n).Name
		}
	}
	return "?"
}
