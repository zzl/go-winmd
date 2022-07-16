package mdmodel

type TypeRefRow struct {
	BaseRow
	ResolutionScope interface{} //ModuleRow/ModuleRefRow/AssemblyRefRow/TypeRefRow
	TypeName        string
	TypeNamespace   string
}

func (this *TypeRefRow) GetFullTypeName() string {
	if this.TypeNamespace != "" {
		return this.TypeNamespace + "." + this.TypeName
	} else {
		return this.TypeName
	}
}

func (this *TypeRefRow) String() string {
	return this.TypeNamespace + "." + this.TypeName
}

type TypeRefTable struct {
	BaseTable[TypeRefRow]
}

func (this *TypeRefTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	ResolutionScope := this.md.CodedIndexes.ResolutionScope

	for n := range this.Rows {
		row := &this.Rows[n]
		row.ResolutionScope = ResolutionScope.GetByAddr(addr)
		addr += ResolutionScope.Size
		row.TypeName = this.md.StringHeap.GetByIndexAddr(addr)
		addr += this.md.StringHeap.IndexSize
		row.TypeNamespace = this.md.StringHeap.GetByIndexAddr(addr)
		addr += this.md.StringHeap.IndexSize
	}
	return addr

}
