package mdmodel

type AssemblyOSRow struct {
	BaseRow
	OSPlatformID   uint32
	OSMajorVersion uint32
	OSMinorVersion uint32
}

type AssemblyOSTable struct {
	BaseTable[AssemblyOSRow]
}

func (this *AssemblyOSTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.OSPlatformID)
		readAddr(&addr, &row.OSMajorVersion)
		readAddr(&addr, &row.OSMinorVersion)
	}
	return addr
}
