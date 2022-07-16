package mdmodel

type ParamRow struct { //P267
	BaseRow
	Flags    ParamAttributesEnum
	Sequence uint16
	Name     string
}

type ParamTable struct {
	BaseTable[ParamRow]
}

func (this *ParamTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.Flags)
		readAddr(&addr, &row.Sequence)
		this.md.StringHeap.ReadAddr(&addr, &row.Name)
	}
	return addr
}
