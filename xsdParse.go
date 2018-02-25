package WSDLCompiler_Go

import (
	"github.com/beevik/etree"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)


/*
Parse all namespaces dependencies

Using:
doc := WSDLCompiler_Go.Reader("wsdl/ptz.wsdl")
tmp := make(map[string]string)
WSDLCompiler_Go.ParseAllSpaces(doc, &tmp, "")
*/

func ParseAllSpaces( doc *etree.Document,spaces *map[string]string, nmsp string)  {
	namespaces := searchNamespaces(doc)
	prefPath := findImportsNamespaces(doc, *namespaces)
	fileIncludes := fileIncludes(doc)
	for _, value := range fileIncludes {
		tmp := (*spaces)[nmsp]
		(*spaces)[nmsp] = tmp + "|" + value
	}
	for key, value := range *prefPath {
		strProcessing := strings.Split(value, "/")
		//if value[len(value)-4:] == ".xsd" {
			if strProcessing[0] == "http:" {
				data, _ := UrlReader(value)
				(*spaces)[key] = value
				ParseAllSpaces(data, spaces, key)
			} else {
				data := Reader("scheme/"+ strProcessing[len(strProcessing)-1])
				(*spaces)[key] = value
				ParseAllSpaces(data, spaces, key)
			}
		//}
	}
}



func searchNamespaces ( doc *etree.Document ) (*map[string]string) {
	prefixMap := make(map[string]string)
	definitions := doc.FindElement("//definitions")
	schema := doc.FindElement("//schema")
	if definitions != nil {
		for _, param := range definitions.Attr {
			if param.Space == "xmlns" {
				if _, ok := prefixMap[param.Key]; ok == false {
					prefixMap[param.Key] = param.Value
				}
			}
		}
	}
	for _, param := range schema.Attr {
		if param.Space == "xmlns" {
			if _, ok := prefixMap[param.Key]; ok == false {
				prefixMap[param.Key] = param.Value
			}
		}
	}
	return &prefixMap
}

func UrlReader( url string ) (doc *etree.Document, err error) {
	var client http.Client
	newDoc := etree.NewDocument()
	resp, err := client.Get(url)
	if err != nil {
		return newDoc, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		if err := newDoc.ReadFromBytes(bodyBytes); err != nil {
			panic(err)
		}
	}
	return newDoc, nil
}

func Reader( path string ) (doc *etree.Document) {
	newDoc := etree.NewDocument()
	if err := newDoc.ReadFromFile(path); err != nil {
		panic(err)
	}
	return newDoc
}

func findImportsNamespaces( doc *etree.Document, prefixMap map[string]string ) (*map[string]string) {
	bool := false
	prefix := ""
	path := ""
	prefixPath := make(map[string]string)
	imports := doc.FindElements("//import")
	for _, imp := range imports {
		for _, elem := range imp.Attr {
			if elem.Key == "namespace" {
				for key, value := range prefixMap {
					if elem.Value == value {
						prefix = key
						bool = true
					}
				}
			}
			if elem.Key == "schemaLocation" && bool {
				path = elem.Value
				prefixPath[prefix] = path
				bool = false
			}
		}
	}
	return &prefixPath
}

func fileIncludes( doc *etree.Document ) ([]string) {
	var includeCollect []string
	includes := doc.FindElements("//include")
	for _, include := range includes {
		includeCollect = append(includeCollect, include.SelectAttr("schemaLocation").Value)
	}
	return includeCollect
}

// Types realization

func GetXSDElements( doc *etree.Document ) {
	complexTypes := doc.FindElements("//complexType")
	simpleType := doc.FindElements("//simpleType")

	for _, xsdComplexType := range complexTypes {
		fmt.Println(xsdComplexType)
	}

	for _, xsdSimpleType := range simpleType {
		fmt.Println(xsdSimpleType)
	}
}