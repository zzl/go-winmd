package main

import (
	"github.com/zzl/go-winmd/apimodel"
	"github.com/zzl/go-winmd/mdmodel"
	"log"
)

func main() {
	mdFilePath := `C:\Windows\System32\WinMetadata\Windows.Foundation.winmd`

	mdModelParser := mdmodel.NewModelParser()
	mdModel, err := mdModelParser.Parse(mdFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer mdModel.Close()

	apiModelParser := apimodel.NewModelParser(nil)
	apiModel := apiModelParser.Parse(mdModel)

	for _, ns := range apiModel.AllNamespaces {
		if ns.Name == "" || len(ns.Types) == 0 {
			continue
		}
		println(ns.FullName)
		for _, typ := range ns.Types {
			println("\t" + typ.Name)
			if typ.Struct {
				for _, f := range typ.StructDef.Fields {
					println("\t\t" + f.Name + " " + f.Type.Name)
				}
			} else if typ.Interface {
				for _, m := range typ.InterfaceDef.Methods {
					print("\t\tfunc " + m.Name)
					print("(")
					for i, p := range m.Params {
						if i > 0 {
							print(", ")
						}
						print(p.Name + " " + p.Type.Name)
					}
					print(") ")
					println(m.ReturnType.Name)
				}
			}
		}
	}
}
