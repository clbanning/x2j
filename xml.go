// Copyright 2012-2016 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package x2j

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"strings"
)

// xmlToMap - convert a XML doc into map[string]interface{} value
func xmlToMap(doc []byte, r bool) (map[string]interface{}, error) {
	b := bytes.NewReader(doc)
	p := xml.NewDecoder(b)
	p.CharsetReader = X2jCharsetReader
	return xmlToMapParser("", nil, p, r)
}

// ===================================== where the work happens =============================

// xmlToMapParser (2015.11.12) - load a 'clean' XML doc into a map[string]interface{} directly.
// A refactoring of xmlToTreeParser(), markDuplicate() and treeToMap() - here, all-in-one.
// We've removed the intermediate *node tree with the allocation and subsequent rescanning.
func xmlToMapParser(skey string, a []xml.Attr, p *xml.Decoder, r bool) (map[string]interface{}, error) {
	// NOTE: all attributes and sub-elements parsed into 'na', 'na' is returned as value for 'skey'
	// Unless 'skey' is a simple element w/o attributes, in which case the xml.CharData value is the value.
	var n, na map[string]interface{}

	// Allocate maps and load attributes, if any.
	if skey != "" {
		n = make(map[string]interface{})  // old n
		na = make(map[string]interface{}) // old n.nodes
		if len(a) > 0 {
			for _, v := range a {
				na[`-` + v.Name.Local] = cast(v.Value, r)
			}
		}
	}
	for {
		t, err := p.Token()
		if err != nil {
			if err != io.EOF {
				return nil, errors.New("xml.Decoder.Token() - " + err.Error())
			}
			return nil, err
		}
		switch t.(type) {
		case xml.StartElement:
			tt := t.(xml.StartElement)

			// First call to xmlToMapParser() doesn't pass xml.StartElement - the map key.
			// So when the loop is first entered, the first token is the root tag along
			// with any attributes, which we process here.
			//
			// Subsequent calls to xmlToMapParser() will pass in tag+attributes for
			// processing before getting the next token which is the element value,
			// which is done above.
			if skey == "" {
				return xmlToMapParser(tt.Name.Local, tt.Attr, p, r)
			}

			// If not initializing the map, parse the element.
			// len(nn) == 1, necessarily - it is just an 'n'.
			nn, err := xmlToMapParser(tt.Name.Local, tt.Attr, p, r)
			if err != nil {
				return nil, err
			}

			// The nn map[string]interface{} value is a na[nn_key] value.
			// We need to see if nn_key already exists - means we're parsing a list.
			// This may require converting na[nn_key] value into []interface{} type.
			// First, extract the key:val for the map - it's a singleton.
			var key string
			var val interface{}
			for key, val = range nn {
				break
			}

			// 'na' holding sub-elements of n.
			// See if 'key' already exists.
			// If 'key' exists, then this is a list, if not just add key:val to na.
			if v, ok := na[key]; ok {
				var a []interface{}
				switch v.(type) {
				case []interface{}:
					a = v.([]interface{})
				default: // anything else - note: v.(type) != nil
					a = []interface{}{v}
				}
				a = append(a, val)
				na[key] = a
			} else {
				na[key] = val // save it as a singleton
			}
		case xml.EndElement:
			// len(n) > 0 if this is a simple element w/o xml.Attrs - see xml.CharData case.
			if len(n) == 0 {
				// If len(na)==0 we have an empty element == "";
				// it has no xml.Attr nor xml.CharData.
				// Note: in original node-tree parser, val defaulted to "";
				// so we always had the default if len(node.nodes) == 0.
				if len(na) > 0 {
					n[skey] = na
				} else {
					n[skey] = "" // empty element
				}
			}
			return n, nil
		case xml.CharData:
			// clean up possible noise
			tt := strings.Trim(string(t.(xml.CharData)), "\t\r\b\n ")
			if len(tt) > 0 {
				if len(na) > 0 {
					na["#text"] = cast(tt, r)
				} else if skey != "" {
					n[skey] = cast(tt, r)
				} else {
					// per Adrian (http://www.adrianlungu.com/) catch stray text
					// in decoder stream -
					// https://github.com/clbanning/mxj/pull/14#issuecomment-182816374
					// NOTE: CharSetReader must be set to non-UTF-8 CharSet or you'll get
					// a p.Token() decoding error when the BOM is UTF-16 or UTF-32.
					continue
				}
			}
		default:
			// noop
		}
	}
}

var castNanInf bool

// Cast "Nan", "Inf", "-Inf" XML values to 'float64'.
// By default, these values will be decoded as 'string'.
func CastNanInf(b bool) {
	castNanInf = b
}

// cast - try to cast string values to bool or float64
func cast(s string, r bool) interface{} {
	if r {
		// handle nan and inf
		if !castNanInf {
			switch strings.ToLower(s) {
			case "nan", "inf", "-inf":
				return interface{}(s)
			}
		}

		// handle numeric strings ahead of boolean
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return interface{}(f)
		}
		// ParseBool treats "1"==true & "0"==false
		// but be more strick - only allow TRUE, True, true, FALSE, False, false
		if s != "t" && s != "T" && s != "f" && s != "F" {
			if b, err := strconv.ParseBool(s); err == nil {
				return interface{}(b)
			}
		}
	}
	return interface{}(s)
}

