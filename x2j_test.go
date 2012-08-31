package x2j

import (
	"os"
	"testing"
)

func TestX2j(t *testing.T) {
	fi, fierr := os.Stat("x2j_test.xml")
	if fierr != nil {
		println("fierr:",fierr.Error())
		return
	}
	fh, fherr := os.Open("x2j_test.xml")
	if fherr != nil {
		println("fherr:",fherr.Error())
		return
	}
	defer fh.Close()
	buf := make([]byte,fi.Size())
	_, nerr  :=  fh.Read(buf)
	if nerr != nil {
		println("nerr:",nerr.Error())
		return
	}
	doc := string(buf)
	println("\nXML DOC:\n",doc)

	root, berr := DocToTree(doc)
	if berr != nil {
		println("berr:",berr.Error())
		return
	}
	println("\nTREE:\n",WriteTree(0,root))

	m := make(map[string]interface{})
	m[root.key] = treeToMap(root,false)
	println("\nMAP:\n",WriteMap(m))

	s,serr := DocToJsonIndent(doc,true)
	if serr != nil {
		println("serr:",serr.Error())
	}
	println("\nDocToJsonIndent(doc,true):\n",s)
}

