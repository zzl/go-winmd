package mdmodel

type EventMapRow struct {
	BaseRow
	Parent    *TypeDefRow
	EventList []*EventRow
}

type EventMapTable struct {
	BaseTable[EventMapRow]
}

func (this *EventMapTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr

	TypeDefTable := this.md.Tables.TypeDef
	EventTable := this.md.Tables.Event

	var prevEventIndex uint32
	var prevRow *EventMapRow

	for n := range this.Rows {
		row := &this.Rows[n]

		ReadTableRowByAddr(&addr, TypeDefTable, &row.Parent)

		var eventIndex uint32
		readRowIndexByAddr(EventTable, &addr, &eventIndex)
		if prevRow != nil {
			GetTableRows(EventTable, prevEventIndex, eventIndex, &prevRow.EventList)
		}
		prevEventIndex = eventIndex
		prevRow = row
	}
	if prevRow != nil {
		GetTableRows(EventTable, prevEventIndex, 0, &prevRow.EventList)
	}
	return addr
}
