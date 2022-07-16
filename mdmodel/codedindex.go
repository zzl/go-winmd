package mdmodel

import (
	"math"
	"unsafe"
)

type CodedIndex struct {
	Tables  []Table
	TagBits byte
	Size    uintptr
}

func NewCodedIndex(tables ...Table) *CodedIndex {
	codedIndex := &CodedIndex{Tables: tables}
	codedIndex.init()
	return codedIndex
}

func (this *CodedIndex) init() {
	maxRowCount := GetMaxRowCount(this.Tables...)
	this.TagBits = byte(math.Ceil(math.Log2(float64(len(this.Tables)))))
	if 1<<(16-this.TagBits) < maxRowCount {
		this.Size = 4
	} else {
		this.Size = 2
	}

}

func (this *CodedIndex) ReadAddr(paddr *uintptr, pi *Row) {
	*pi = this.GetByAddr(*paddr)
	*paddr += this.Size
}

func (this *CodedIndex) GetByAddr(addr uintptr) Row {
	var tableIndex int
	var rowIndex int

	if this.Size == 2 {
		index := *(*uint16)(unsafe.Pointer(addr))
		nonTagBits := 16 - this.TagBits
		tableIndex = int((index << nonTagBits) >> nonTagBits)
		rowIndex = int(index >> this.TagBits)

	} else {
		index := *(*uint32)(unsafe.Pointer(addr))
		nonTagBits := 32 - this.TagBits
		tableIndex = int((index << nonTagBits) >> nonTagBits)
		rowIndex = int(index >> this.TagBits)
	}
	table := this.Tables[tableIndex]
	if rowIndex == 0 {
		return nil //??
	}
	return table.GetRow(rowIndex)
}

type CodedIndexes struct {
	metadata *Model

	TypeDefOrRef        *CodedIndex
	HasConstant         *CodedIndex
	HasCustomAttribute  *CodedIndex
	HasFieldMarshall    *CodedIndex
	HasDeclSecurity     *CodedIndex
	MemberRefParent     *CodedIndex
	HasSemantics        *CodedIndex
	MethodDefOrRef      *CodedIndex
	MemberForwarded     *CodedIndex
	Implementation      *CodedIndex
	CustomAttributeType *CodedIndex
	ResolutionScope     *CodedIndex
	TypeOrMethodDef     *CodedIndex
}

func NewCodedIndexes(metaData *Model) *CodedIndexes {
	codedIndexes := &CodedIndexes{metadata: metaData}
	codedIndexes.init()
	return codedIndexes
}

func (this *CodedIndexes) init() {
	tabs := this.metadata.Tables
	this.TypeDefOrRef = NewCodedIndex(tabs.TypeDef, tabs.TypeRef, tabs.TypeSpec)
	this.HasConstant = NewCodedIndex(tabs.Field, tabs.Param, tabs.Property)

	this.HasCustomAttribute = NewCodedIndex(
		tabs.MethodDef, tabs.Field, tabs.TypeRef, tabs.TypeDef,
		tabs.Param, tabs.InterfaceImpl, tabs.MemberRef, tabs.Module,
		nil /*?Permission*/, tabs.Property, tabs.Event, tabs.StandAloneSig,
		tabs.ModuleRef, tabs.TypeSpec, tabs.Assembly, tabs.AssemblyRef,
		tabs.File, tabs.ExportedType, tabs.ManifestResource, tabs.GenericParam,
		tabs.GenericParamConstraint, tabs.MethodSpec)

	this.HasFieldMarshall = NewCodedIndex(tabs.Field, tabs.Param)
	this.HasDeclSecurity = NewCodedIndex(tabs.TypeDef, tabs.MethodDef, tabs.Assembly)

	this.MemberRefParent = NewCodedIndex(tabs.TypeDef,
		tabs.TypeRef, tabs.ModuleRef, tabs.MethodDef, tabs.TypeSpec)

	this.HasSemantics = NewCodedIndex(tabs.Event, tabs.Property)
	this.MethodDefOrRef = NewCodedIndex(tabs.MethodDef, tabs.MemberRef)
	this.MemberForwarded = NewCodedIndex(tabs.Field, tabs.MethodDef)
	this.Implementation = NewCodedIndex(tabs.File, tabs.AssemblyRef, tabs.ExportedType)

	this.CustomAttributeType = NewCodedIndex(
		nil, nil, tabs.MethodDef, tabs.MemberRef, nil)

	this.ResolutionScope = NewCodedIndex(tabs.Module,
		tabs.ModuleRef, tabs.AssemblyRef, tabs.TypeRef)

	this.TypeOrMethodDef = NewCodedIndex(tabs.TypeDef, tabs.MethodDef)
}
