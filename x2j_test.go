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
	fmt.Println("\nDocToTree():\n",WriteTree(root))

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

}

