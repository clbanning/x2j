package x2j

import (
	"bytes"
	"fmt"
	"testing"
)

var doc = `<entry><vars><foo>bar</foo><foo2><hello>world</hello></foo2></vars></entry>`

func TestToTree(t *testing.T) {
	fmt.Println("\nToTree - Read doc:",doc)
	rdr := bytes.NewBufferString(doc)
	n,err := ToTree(rdr)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println(n.WriteTree())
}

func TestToMap(t *testing.T) {
	fmt.Println("\nToMap - Read doc:",doc)
	rdr := bytes.NewBufferString(doc)
	m,err := ToMap(rdr)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println(WriteMap(m))
}

func TestToJson(t *testing.T) {
	fmt.Println("\nToJson - Read doc:",doc)
	rdr := bytes.NewBufferString(doc)
	s,err := ToJson(rdr)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("json:",s)
}

func TestToJsonIndent(t *testing.T) {
	fmt.Println("\nToJsonIndent - Read doc:",doc)
	rdr := bytes.NewBufferString(doc)
	s,err := ToJsonIndent(rdr)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("json:",s)
}

