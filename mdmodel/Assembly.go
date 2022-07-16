package mdmodel

type AssemblyRow struct {
	BaseRow
	HashAlgId      AssemblyHashAlgorithmEnum
	MajorVersion   uint16
	MinorVersion   uint16
	BuildNumber    uint16
	RevisionNumber uint16
	Flags          AssemblyFlagsEnum
	PublicKey      []byte //blob
	Name           string
	Culture        string
}

type AssemblyTable struct {
	BaseTable[AssemblyRow]
}

func (this *AssemblyTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.HashAlgId)
		readAddr(&addr, &row.MajorVersion)
		readAddr(&addr, &row.MinorVersion)
		readAddr(&addr, &row.BuildNumber)
		readAddr(&addr, &row.RevisionNumber)
		readAddr(&addr, &row.Flags)

		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.PublicKey = blob.Data

		this.md.StringHeap.ReadAddr(&addr, &row.Name)
		this.md.StringHeap.ReadAddr(&addr, &row.Culture)
	}
	return addr
}
