package mdmodel

import (
	"log"
)

type MethodDefRow struct {
	BaseRow
	Rva       uint32
	ImplFlags MethodImplAttributesEnum
	Flags     MethodAttributesEnum
	Name      string
	Signature *MethodDefSig
	ParamList []*ParamRow

	OwnerType *TypeDefRow
}

func (this *MethodDefRow) String() string {
	sig := this.Signature
	var s string
	s += "func " + this.Name + "("
	var retParam *ParamRow
	index := 0
	for n, param := range this.ParamList {
		if param.Sequence == 0 {
			if n != 0 {
				log.Panic("?")
			}
			retParam = param
			continue
		}
		if index > 0 {
			s += ", "
		}
		s += param.Name + " " + sig.Params[index].String()
		index += 1
	}
	s += ")"
	_ = retParam
	if retParam != nil {
		s += " " + sig.RetType.String()
	}
	return s
}

type MethodDefTable struct {
	BaseTable[MethodDefRow]
}

func (this *MethodDefTable) ParseRows(baseAddr uintptr) uintptr {
	addr := baseAddr
	ParamTable := this.md.Tables.Param

	var prevParamIndex uint32
	var prevRow *MethodDefRow

	for n := range this.Rows {
		row := &this.Rows[n]

		readAddr(&addr, &row.Rva)
		readAddr(&addr, &row.ImplFlags)
		readAddr(&addr, &row.Flags)
		this.md.StringHeap.ReadAddr(&addr, &row.Name)

		blob := this.md.BlobHeap.ReadAddr(&addr)
		row.Signature, _ = ParseMethodDefSig(this.md, blob.Data)

		//
		var paramIndex uint32
		readRowIndexByAddr(ParamTable, &addr, &paramIndex)
		if prevRow != nil {
			GetTableRows(ParamTable, prevParamIndex, paramIndex, &prevRow.ParamList)
		}

		prevParamIndex = paramIndex
		prevRow = row
	}
	if prevRow != nil {
		GetTableRows(ParamTable, prevParamIndex, 0, &prevRow.ParamList)
	}
	return addr
}
