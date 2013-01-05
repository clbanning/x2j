package x2j

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestX2j(t *testing.T) {
	fi, fierr := os.Stat("x2j_test.xml")
	if fierr != nil {
		fmt.Println("fierr:",fierr.Error())
		return
	}
	fh, fherr := os.Open("x2j_test.xml")
	if fherr != nil {
		fmt.Println("fherr:",fherr.Error())
		return
	}
	defer fh.Close()
	buf := make([]byte,fi.Size())
	_, nerr  :=  fh.Read(buf)
	if nerr != nil {
		fmt.Println("nerr:",nerr.Error())
		return
	}
	doc := string(buf)
	fmt.Println("\nXML doc:\n",doc)

	root, berr := DocToTree(doc)
	if berr != nil {
		fmt.Println("berr:",berr.Error())
		return
	}
	fmt.Println("\nDocToTree():\n",root.WriteTree())

	m := make(map[string]interface{})
	m[root.key] = root.treeToMap(false)
	fmt.Println("\ntreeToMap, recast==false:\n",WriteMap(m))

	j,jerr := json.MarshalIndent(m,"","  ")
	if jerr != nil {
		fmt.Println("jerr:",jerr.Error())
	}
	fmt.Println("\njson.MarshalIndent, recast==false:\n",string(j))

	// test DocToMap() with recast
	mm, mmerr := DocToMap(doc,true)
	if mmerr != nil {
		println("mmerr:",mmerr.Error())
		return
	}
	println("\nDocToMap(), recast==true:\n",WriteMap(mm))

	// test DocToJsonIndent() with recast
	s,serr := DocToJsonIndent(doc,true)
	if serr != nil {
		fmt.Println("serr:",serr.Error())
	}
	fmt.Println("\nDocToJsonIndent, recast==true:\n",s)

	// test ValueInMap()
	doc = `<entry><vars><foo>bar</foo><foo2><hello>world</hello></foo2></vars></entry>`
	fmt.Println("\nRead doc:",doc)
	fmt.Println("Looking for value: entry.vars")
	mm,mmerr = DocToMap(doc)
	if mmerr != nil {
		fmt.Println("merr:",mmerr.Error())
	}
	v,verr := ValueInMap(mm,"entry.vars")
	if verr != nil {
		fmt.Println("verr:",verr.Error())
	} else {
		j, jerr := json.MarshalIndent(v,"","  ")
		if jerr != nil {
			fmt.Println("jerr:",jerr.Error())
		} else {
			fmt.Println(string(j))
		}
	}
	fmt.Println("Looking for value: entry.vars.foo2.hello")
	v,verr = ValueInMap(mm,"entry.vars.foo2.hello")
	if verr != nil {
		fmt.Println("verr:",verr.Error())
	} else {
		fmt.Println(v.(string))
	}
	fmt.Println("Looking with error in path: entry.var")
	v,verr = ValueInMap(mm,"entry.var")
	fmt.Println("verr:",verr.Error())

	// test DocToValue()
	fmt.Println("DocToValue() for tag path entry.vars")
	v,verr = DocToValue(doc,"entry.vars")
	if verr != nil {
		fmt.Println("verr:",verr.Error())
	}
	j,_ = json.MarshalIndent(v,"","  ")
	fmt.Println(string(j))
}

