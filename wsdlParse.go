package WSDLCompiler_Go

import (
	"github.com/beevik/etree"
	"os"
)

/* Parse WSDL and generate datatypes
using:
doc := WSDLCompiler_Go.Reader("wsdl/ptz.wsdl")
operations := WSDLCompiler_Go.GetAllOperations(doc)
for _, element := range *operations {
	funcMap := WSDLCompiler_Go.GetWSDLElements(doc, element)
	WSDLCompiler_Go.MakeWSDLStructs(funcMap)
}
for _, element := range *operations {
	funcMap := WSDLCompiler_Go.GetWSDLElements(doc, element)
	WSDLCompiler_Go.MakeWSDLFuncs(funcMap)
}
*/

type data struct {
	FuncName string
	FuncArgs []map[string] string
}

func GetWSDLElements( doc *etree.Document, element string ) (* data) {
	mainElements := doc.FindElements("//element[@name='"+ element +"']")
	dataCollection := new(data)
	for _, wsdlOperation := range mainElements {
		dataCollection.FuncName = wsdlOperation.SelectAttr("name").Value//Add Name (e.g. AbsoluteMove)
		getChildElements := doc.FindElements("//element[@name='"+ element +"']//*")
		var elementAttr []map[string]string
		for _, childElement := range getChildElements {
			if childElement.Tag == "element" {
				tmpMap := make(map[string]string)
				for _, funcArgs := range childElement.Attr {
					tmpMap[funcArgs.Key] = funcArgs.Value
				}
				elementAttr = append(elementAttr, tmpMap)
				dataCollection.FuncArgs = elementAttr
			}
		}
	}
	return dataCollection
}


func MakeWSDLStructs( data *data )  {
	tmpStr := "\n\ntype " + data.FuncName + " struct {\n"
	for _, param := range data.FuncArgs {
		tmpStr += "\t" + param["name"] + " " + param["type"] + "\n"
	}
	tmpStr += "\n}\n"
	write(tmpStr, "types/funcs.go")
}

//Generate Go functions from WSDL
func MakeWSDLFuncs( data *data )  {
	tmpStr := "\nfunc Fn" + data.FuncName + " (arg " + data.FuncName +") {\n"
	tmpStr += "\n}\n"
	write(tmpStr, "types/funcs.go")
}

func write(data string, fileName string)  {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(data); err != nil {
		panic(err)
	}
}

func GetAllOperations( doc *etree.Document ) *[]string {
	var allOperations []string
	portType := doc.FindElement("//portType")
	//fmt.Println(portType)
	for _, operation := range portType.SelectElements("wsdl:operation") {
		operationName := operation.SelectAttr("name").Value
		allOperations = append(allOperations, operationName)
	}
	return &allOperations
}