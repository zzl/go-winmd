package apimodel

type Namespace struct {
	Name     string
	FullName string
	Parent   *Namespace
	Children []*Namespace

	Types []*Type //struct?
}
