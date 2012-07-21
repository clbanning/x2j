// Copyright 2012 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file
/*	
	Unmarshal an arbitrary XML document to a map or a JSON string. 

	DocToMap() returns an intermediate result with the XML doc unmarshal'd to a map
	of type map[string]interface{}. It is analogous to unmarshal'ng a JSON string to
	a map using json.Unmarshal(). (This was the original purpose of this library.)

	DocToTree()/WriteTree() let you examine the parsed XML doc.

	XML values are all type 'string'. The optional argument 'recast' for DocToJson()
	and DocToMap() will convert element values to JSON data types - 'float64' and 'bool' -
	if possible.  This, however, should be done with caution as it will recast ALL numeric
	and boolean string values, even those that are meant to be of type 'string'.
 */
package x2j

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type node struct {
	dup bool
	key string
	val string
	nodes []*node
}

// DocToJson - return an XML doc as a JSON string.
//	If the optional argument 'recast' is 'true', then values will be converted to boolean or float64 if possible.
func DocToJson(doc string,recast ...bool) (string,error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m,merr := DocToMap(doc,r)
	if m == nil || merr != nil {
		return "",merr
	}

	b, berr := json.Marshal(m)
	if berr != nil {
		return "",berr
	}

	return string(b),nil
}

// DocToJsonIndent - return an XML doc as a prettified JSON string.
//	If the optional argument 'recast' is 'true', then values will be converted to boolean or float64 if possible.
func DocToJsonIndent(doc string,recast ...bool) (string,error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m,merr := DocToMap(doc,r)
	if m == nil || merr != nil {
		return "",merr
	}

	b, berr := json.MarshalIndent(m,"","  ")
	if berr != nil {
		return "",berr
	}

	return string(b),nil
}

// DocToMap - convert an XML doc into a map[string]interface{}.
// (This is analogous to unmarshalling a JSON string to map[string]interface{} using json.Unmarshal().)
//	If the optional argument 'recast' is 'true', then values will be converted to boolean or float64 if possible.
func DocToMap(doc string,recast ...bool) (map[string]interface{},error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	n,err := DocToTree(doc)
	if err != nil {
		return nil,err
	}

	m := make(map[string]interface{})
	m[n.key] = treeToMap(n,r)

	return m,nil
}

// DocToTree - convert an XML doc into a tree of nodes.
func DocToTree(doc string) (*node, error) {
	reg,_ := regexp.Compile("[ \t]*<")
	doc = reg.ReplaceAllString(doc,"<")
	r := strings.NewReplacer("\n","","\r","","\t","")
	doc = r.Replace(doc)
	b := bytes.NewBufferString(doc)
	p := xml.NewDecoder(b)

	n, berr := xmlToTree("",nil,p)
	if berr != nil {
		return nil, berr
	}

	return n,nil
}

// WriteTree - convert a tree of nodes into a printable string.
//	'indent' is the starting indentation count; typically: WriteTree(0,n).
func WriteTree(indent int,n *node) string {
	var s string
	if n.val != "" {
		for i := 0 ; i < indent ; i++ {
			s += "  "
		}
		s += n.key+" : "+n.val+"\n"
	} else {
		for i := 0 ; i < indent ; i++ {
			s += "  "
		}
		s += n.key+" :"+"\n"
		for _,v := range n.nodes {
			s += WriteTree(indent+1,v)
		}
	}
	return s
}

// xmlToTree - load a 'clean' XML doc into a tree of *node.
func xmlToTree(skey string,a []xml.Attr,p *xml.Decoder) (*node, error) {
	n := new(node)
	n.nodes = make([]*node,0)

	if skey != "" {
		n.key = skey
		if len(a) > 0 {
			for _,v := range a {
				na := new(node)
				na.key = `-`+v.Name.Local
				na.val = v.Value
				n.nodes = append(n.nodes,na)
			}
		}
	}
	for {
		t,err := p.Token()
		if err != nil {
			return nil, err
		} else {
			switch t.(type) {
				case xml.StartElement:
					tt := t.(xml.StartElement)
					// handle root
					if n.key == "" {
						n.key = tt.Name.Local
						if len(tt.Attr) > 0 {
							for _,v := range tt.Attr {
								na := new(node)
								na.key = `-`+v.Name.Local
								na.val = v.Value
								n.nodes = append(n.nodes,na)
							}
						}
					} else {
						nn, nnerr := xmlToTree(tt.Name.Local,tt.Attr,p)
						if nnerr != nil {
							return nil, nnerr
						}
						n.nodes = append(n.nodes,nn)
					}
				case xml.EndElement:
					// scan v.nodes for duplicate v.key values
					markDuplicateKeys(n)
					return n, nil
				case xml.CharData:
					tt := string(t.(xml.CharData))
					if len(n.nodes) > 0 {
						nn := new(node)
						nn.key = "#text"
						nn.val = tt
						n.nodes = append(n.nodes,nn)
					} else {
						n.val = tt
					}
				default:
					// noop
			}
		}
	}
	return nil, errors.New("EndElement not found for: "+n.key)
}

// markDuplicateKeys - set node.dup flag for loading map[string]interface{}.
func markDuplicateKeys(n *node) {
	l := len(n.nodes)
	for i := 0 ; i < l ; i++ {
		if n.nodes[i].dup {
			continue
		}
		for j := i+1 ; j < l ; j++ {
			if n.nodes[i].key == n.nodes[j].key {
				n.nodes[i].dup = true
				n.nodes[j].dup = true
			}
		}
	}
}

// treeToMap - convert a tree of nodes into a map[string]interface{}.
//	(Parses to map that is structurally the same as from json.Unmarshal().)
// Note: root is not instantiated; call with: "m[n.key] = treeToMap()".
func treeToMap(n *node,r bool) interface{} {
	if len(n.nodes) == 0 {
		return recast(n.val,r)
	}

	m := make(map[string]interface{},0)
	for _,v := range n.nodes {
		// just a value
		if !v.dup && len(v.nodes) == 0 {
			m[v.key] = recast(v.val,r)
			continue
		}

		// a list of values
		if v.dup {
			var a []interface{}
			if vv,ok := m[v.key]; ok {
				a = vv.([]interface{})
			} else {
				a = make([]interface{},0)
			}
			a = append(a,treeToMap(v,r))
			m[v.key] = interface{}(a)
			continue
		}

		// it's a unique key
		m[v.key] = treeToMap(v,r)
	}

	return interface{}(m)
}

// recast - try to cast string values to bool or float64
func recast(s string,r bool) interface{} {
	if r {
		// handle numeric strings ahead of boolean
		if f, err := strconv.ParseFloat(s,64); err == nil {
			return interface{}(f)
		}
		// ParseBool treats "1"==true & "0"==false
		if b, err := strconv.ParseBool(s); err == nil {
			return interface{}(b)
		}
	}
	return interface{}(s)
}

// WriteMap - dumps the map[string]interface{} for examination.
//	'indent' is initial indentation count; typically: WriteMap(0,m).
//	NOTE: with XML all element types are 'string'.
//	But code written as generic for use with maps[string]interface{} values from json.Unmarshal().
func WriteMap(indent int,m interface{}) string {
	var s string
	switch m.(type) {
		case nil:
			return "[nil] nil"
		case string:
			return "[string] "+m.(string)
		case float64:
			return "[float64] "+strconv.FormatFloat(m.(float64),'e',2,64)
		case bool:
			return "[bool] "+strconv.FormatBool(m.(bool))
		case []interface{}:
			s += "[[]interface{}]"
			for i,v := range m.([]interface{}) {
				s += "\n"
				for i := 0 ; i < indent ; i++ {
					s += "  "
				}
				s += "item: "+strconv.FormatInt(int64(i),10)
				switch v.(type) {
					case string,float64,bool:
						s += "\n"
					default:
						// noop
				}
				for i := 0 ; i < indent ; i++ {
					s += "  "
				}
				s += WriteMap(indent+1,v)
			}
		case map[string]interface{}:
			for k,v := range m.(map[string]interface{}) {
				s += "\n"
				for i := 0 ; i < indent ; i++ {
					s += "  "
				}
				s += "[map[string]interface{}] "+k+" :"+WriteMap(indent+1,v)
		}
		default:
			// shouldn't ever be here ...
			s += fmt.Sprintf("unknown type for: %v",m)
	}
	return s
}
