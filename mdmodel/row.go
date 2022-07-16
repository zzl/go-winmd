package mdmodel

type Row interface {
	RowIndex() uint32
}

type BaseRow struct {
	index uint32
}

func (this *BaseRow) RowIndex() uint32 {
	return this.index
}

type rowIndexSetter interface {
	setRowIndex(index uint32)
}

func (this *BaseRow) setRowIndex(index uint32) {
	this.index = index
}

type TypeRow interface {
	GetFullTypeName() string
}
