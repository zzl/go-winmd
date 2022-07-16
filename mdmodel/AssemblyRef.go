package mdmodel

type AssemblyRefRow struct {
	BaseRow
	MajorVersion     uint16
	MinorVersion     uint16
	BuildNumber      uint16
	RevisionNumber   uint16
	Flags            AssemblyFlagsEnum
	PublicKeyOrToken []byte
	Name             string
	Culture          string
	HashValue        []byte
}

type AssemblyRefTable struct {
	BaseTable[AssemblyRefRow]
}

func (this *AssemblyRefTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]

		readAddr(&addr, &row.MajorVersion)
		readAddr(&addr, &row.MinorVersion)
		readAddr(&addr, &row.BuildNumber)
		readAddr(&addr, &row.RevisionNumber)
		readAddr(&addr, &row.Flags)

		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.PublicKeyOrToken = blob.Data

		this.md.StringHeap.ReadAddr(&addr, &row.Name)
		this.md.StringHeap.ReadAddr(&addr, &row.Culture)

		blob = this.md.BlobHeap.ReadAddr(&addr)
		row.HashValue = blob.Data
	}
	return addr
}
