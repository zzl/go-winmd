package mdmodel

type ImplMapRow struct {
	BaseRow
	MappingFlags    PInvokeAttributesEnum
	MemberForwarded Row //FieldRow/MethodDefRow
	ImportName      string
	ImportScope     *ModuleRefRow
}

type ImplMapTable struct {
	BaseTable[ImplMapRow]

	MemberForwardImplRowMap map[uint32]*ImplMapRow
}

func (this *ImplMapTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	MemberForwarded := this.md.CodedIndexes.MemberForwarded
	ModuleRefTable := this.md.Tables.ModuleRef

	this.MemberForwardImplRowMap = make(map[uint32]*ImplMapRow)
	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.MappingFlags)
		MemberForwarded.ReadAddr(&addr, &row.MemberForwarded)
		this.md.StringHeap.ReadAddr(&addr, &row.ImportName)
		ReadTableRowByAddr(&addr, ModuleRefTable, &row.ImportScope)
		this.MemberForwardImplRowMap[row.MemberForwarded.RowIndex()] = row
	}
	return addr
}
