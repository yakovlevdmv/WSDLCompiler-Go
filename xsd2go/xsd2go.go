package xsd2go

import (
	"net/http"
	"io/ioutil"
	"errors"
	"strconv"
	"github.com/beevik/etree"
	"fmt"
	"os"
)

var namespaces = make(map[string]etree.Document)

/*
	Функция получает XSD документ по URL адресу
 */
func getXsdByURL(url string) ([]byte, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("Got response status code: " + strconv.Itoa(resp.StatusCode))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return respBody, string(respBody), nil
}

func getXsdByPath(path string) ([]byte, string,error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	return dat, string(dat), nil
}

func processSimpleType(element etree.Element) string {
	var result = "type " + element.SelectAttr("name").Value + " struct {\n"
	fmt.Println(element.SelectAttr("name").Value)
	for _, i := range element.FindElements("./*") {
		switch i.Tag {
			case "restriction":
				if i.FindElement("./*") != nil && i.FindElement("./*").Tag == "enumeration" {
					var enum = "type " + element.SelectAttr("name").Value + "_enum" + " int\n" + "const (\n"
					var first = true
					for _, j := range i.FindElements("./*") {
						enum += "\t" + j.SelectAttr("value").Value
						if first {
							enum += " " + element.SelectAttr("name").Value + "_enum" + " = iota"
							first = false
						}
						enum += "\n"
					}
					enum += ")\n"
					//var enum = "var " + element.SelectAttr("name").Value + "_enum" + " = [...]string{\n"
					//for _, j := range i.FindElements("./*") {
					//	enum += "\t" + "\"" + j.SelectAttr("value").Value + "\",\n"
					//}
					//enum += "}\n"
					result = enum + result
				} else {
					_type := i.SelectAttr("base").Value[3:]
					if _type == "float" {
						_type = "float64"
					} else if _type == "integer" {
						_type = "int"
					}
					result = result + "\t" + "value " + _type
					fmt.Println("\t", i)
				}

				result += "\n}"

				break
			case "union":
				fmt.Println("\t", i)
				result += "\n}"
				break
			case "list":
				fmt.Println("\t", i)
				_type := i.SelectAttr("itemType").Value[3:]
				if _type == "float" {
					_type = "float64"
				} else if _type == "integer" {
					_type = "int"
				}
				result += "\t" + "value " + "[]" + _type + "\n}"
				break
		}
	}
	//var result = "type " + element.SelectAttr("name").Value + " struct {\n"
	//fmt.Println(element.SelectAttr("name").Value)
	//restrictions := element.SelectElements("restriction")
	//for _, i := range restrictions {
	//	result = result + "\t" + "value " + i.SelectAttr("base").Value + "\n"
	//	for _, e := range i.FindElements("./*") {
	//		fmt.Print(e)
	//	}
	//}
	//result = result + "}"
	////fmt.Println("RESULT: \n" + result)

	//return result
	return result
}

func processComplexType(element etree.Element) {
	attributes := element.SelectElements("attribute")
	fmt.Println(attributes)
	fmt.Println(attributes[0])

}

func ProcessXSD() {
	_, xsd_data, err := getXsdByURL("https://www.onvif.org/ver10/schema/onvif.xsd")
	if err != nil {
		panic(err)
	}

	xsd_document := etree.NewDocument()
	if err := xsd_document.ReadFromString(xsd_data); err != nil {
		panic(err)
	}
	root := xsd_document.SelectElement("schema")

	//includes := root.SelectElements("include")
	//imports  := root.SelectElements("import")
	simpleTypes := root.SelectElements("simpleType")
	//complexTypes := root.SelectElements("complexType")

	//for _, _import := range imports {
	//	for _, attr := range root.Attr {
	//		if attr.Space == "xmlns" && _import.SelectAttr("namespace").Value == attr.Value{
	//			_, xsd, err := getXsdByURL(_import.SelectAttr("schemaLocation").Value)
	//			if err != nil {
	//				log.Println(_import.SelectAttr("schemaLocation").Value)
	//				panic(err)
	//			}
	//			doc := etree.NewDocument()
	//			if err := doc.ReadFromString(xsd); err != nil {
	//				panic(err)
	//			}
	//			namespaces[attr.Key] = *doc
	//		}
	//	}
	//}
	//fmt.Println(namespaces)

	f, err := os.Create("onvif.go")
	if err != nil {
		panic(err)
	}

	f.WriteString("package xsd2go\n\n")
	for _, i := range simpleTypes {
		f.WriteString("\n")
		f.WriteString(processSimpleType(*i))
		f.WriteString("\n")
	}

	f.Close()


}