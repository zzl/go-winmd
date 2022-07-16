package mdmodel

type ClassLayoutRow struct {
	BaseRow
	PackingSize uint16
	ClassSize   uint32
	Parent      *TypeDefRow
}

type ClassLayoutTable struct {
	BaseTable[ClassLayoutRow]
}

func (this *ClassLayoutTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	TypeDefTable := this.md.Tables.TypeDef
	for n := range this.Rows {
		row := &this.Rows[n]

		readAddr(&addr, &row.PackingSize)
		readAddr(&addr, &row.ClassSize)

		ReadTableRowByAddr(&addr, TypeDefTable, &row.Parent)
	}
	return addr
}
