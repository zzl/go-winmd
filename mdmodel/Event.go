package mdmodel

type EventRow struct {
	BaseRow
	EventFlags EventAttributesEnum
	Name       string
	EventType  Row //TypeDefRow/TypeRefRow/TypeSpecRow
}

type EventTable struct {
	BaseTable[EventRow]
}

func (this *EventTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	TypeDefOrRef := this.md.CodedIndexes.TypeDefOrRef
	for n := range this.Rows {
		row := &this.Rows[n]
		readAddr(&addr, &row.EventFlags)
		this.md.StringHeap.ReadAddr(&addr, &row.Name)
		TypeDefOrRef.ReadAddr(&addr, &row.EventType)
	}
	return addr
}
