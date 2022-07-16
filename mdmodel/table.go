package mdmodel

import (
	"log"
	"math"
	"reflect"
	"unsafe"
)

type Table interface {
	SetCode(code byte)
	GetCode() byte
	SetMetaData(metaData *Model)
	GetMetaData() *Model

	SetRowCount(rowCount int)
	ParseRows(baseAddr uintptr) uintptr
	GetRow(index int) Row
}

type BaseTable[TRow any] struct {
	md   *Model
	Code byte
	Rows []TRow
}

func (this *BaseTable[TRow]) GetCode() byte {
	return this.Code
}

func (this *BaseTable[TRow]) SetCode(code byte) {
	this.Code = code
}

func (this *BaseTable[TRow]) SetMetaData(metaData *Model) {
	this.md = metaData
}

func (this *BaseTable[TRow]) GetMetaData() *Model {
	return this.md
}

func (this *BaseTable[TRow]) SetRowCount(rowCount int) {
	this.Rows = make([]TRow, rowCount)
	for n := range this.Rows {
		row := &this.Rows[n]
		any(row).(rowIndexSetter).setRowIndex(uint32(n + 1))
	}
}

func (this *BaseTable[TRow]) ParseRows(baseAddr uintptr) uintptr {
	if len(this.Rows) != 0 {
		log.Panic("not implemented")
	}
	return baseAddr
}

func (this *BaseTable[TRow]) GetRow(index int) Row {
	n := index - 1
	if n > len(this.Rows)-1 {
		return nil //?
	}
	return any(&this.Rows[n]).(Row)
}

func readRowIndexByAddr(table Table, paddr *uintptr, pindex *uint32) {
	tableValue := reflect.ValueOf(table).Elem()
	rowsFieldValue := tableValue.FieldByName("Rows")
	count := rowsFieldValue.Len()
	if count > math.MaxUint16 {
		*pindex = *(*uint32)(unsafe.Pointer(*paddr))
		*paddr += 4
	} else {
		*pindex = uint32(*(*uint16)(unsafe.Pointer(*paddr)))
		*paddr += 2
	}
}

func GetTableRows[T any](table Table, fromIndex uint32, toIndex uint32, pRows *[]*T) {
	var rows []*T
	tableValue := reflect.ValueOf(table).Elem()
	rowsFieldValue := tableValue.FieldByName("Rows")
	if toIndex == 0 {
		toIndex = uint32(rowsFieldValue.Len() + 1)
	}
	for n := fromIndex; n < toIndex; n++ {
		row := rowsFieldValue.Index(int(n - 1))
		rows = append(rows, row.Addr().Interface().(*T))
	}
	*pRows = rows
}

func GetTableRow(table Table, rowIndex uint32) interface{} {
	tableValue := reflect.ValueOf(table).Elem()
	rowsFieldValue := tableValue.FieldByName("Rows")
	row := rowsFieldValue.Index(int(rowIndex - 1))
	return row.Addr().Interface()
}

func ReadTableRowByAddr[T any](paddr *uintptr, table Table, ppRow **T) {
	var index uint32
	readRowIndexByAddr(table, paddr, &index)
	*ppRow = GetTableRow(table, index).(*T)
}
