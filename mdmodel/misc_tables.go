package mdmodel

type AssemblyProcessorRow struct {
	BaseRow
	Processor uint32
}

type AssemblyProcessorTable struct {
	BaseTable[AssemblyProcessorRow]
}

//
type AssemblyRefOSRow struct {
	BaseRow
	OSPlatformID   uint32
	OSMajorVersion uint32
	OSMinorVersion uint32
	AssemblyRef    *AssemblyRefRow
}

type AssemblyRefOSTable struct {
	BaseTable[AssemblyRefOSRow]
}

//
type AssemblyRefProcessorRow struct {
	BaseRow
	Processor   uint32
	AssemblyRef *AssemblyRefRow
}

type AssemblyRefProcessorTable struct {
	BaseTable[AssemblyRefProcessorRow]
}

//
type DeclSecurityRow struct {
	BaseRow
	Action        uint16
	Parent        interface{} //TypeDefRow/MethodDefRow/AssemblyRow
	PermissionSet []byte
}

type DeclSecurityTable struct {
	BaseTable[DeclSecurityRow]
}

//
type ExportedTypeRow struct {
	BaseRow
	Flags          TypeAttributesEnum
	TypeDefId      uint32
	TypeName       string
	TypeNamespace  string
	Implementation interface{} //FileRow/ExportedTypeRow/AssemblyRefRow
}

type ExportedTypeTable struct {
	BaseTable[ExportedTypeRow]
}

//
type FieldRvaRow struct {
	BaseRow
	Rva   uint32
	Field *FieldRow
}

type FieldRvaTable struct {
	BaseTable[FieldRvaRow]
}

//
type FileRow struct {
	BaseRow
	Flags     FileAttributesEnum
	Name      string
	HashValue []byte
}

type FileTable struct {
	BaseTable[FileRow]
}

//
type ManifestResourceRow struct {
	BaseRow
	Offset         uint32
	Flags          ManifestResourceAttributesEnum
	Name           string
	Implementation interface{} //FileRow/AssemblyRefRow/null
}

type ManifestResourceTable struct {
	BaseTable[ManifestResourceRow]
}

//
type MethodSpecRow struct {
	BaseRow
	Method        interface{} //MethodDefRow/MemberRefRow
	Instantiation *Sig
}

type MethodSpecTable struct {
	BaseTable[MethodSpecRow]
}

//
type StandAloneSigRow struct {
	BaseRow
	Signature *Sig
}

type StandAloneSigTable struct {
	BaseTable[StandAloneSigRow]
}
