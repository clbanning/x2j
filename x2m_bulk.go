// Copyright 2012-2013 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file
/*	
	Process files with multiple XML messages.

	Architecture:
		Read file into buffer 
			-> read message from buffer
				-> unmarshal to map[string]interface{}
				-> dispatch to process or error handler
			-> continue until io.EOF
*/
package x2j

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"regexp"
)

type msg struct {
	val map[string]interface{}
	err error
}

// ParseXmlMsgsFromFile()
//	'fname' is name of file
//	'phandler' is the map processing handler. If return of 'false' stops further processing.
//	'ehandler' is the parsing error handler. If return of 'false' stops further processing.
//	Note: phandler() and ehandler() calls are blocking, so reading and processing of messages is serialized.
//	      This means that you can stop reading the file on error or on processing a particular message.
//	      To have reading and handling run concurrently, pass arguments to a go routine in handler and return true.
func ParseXmlMsgsFromFile(fname string, phandler func(map[string]interface{})(bool), ehandler func(error)(bool), recast ...bool) error {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	fi, fierr := os.Stat(fname)
	if fierr != nil {
		return fierr
	}
	fh, fherr := os.Open(fname)
	if fherr != nil {
		return fherr
	}
	defer fh.Close()
	buf := make([]byte,fi.Size())
	_, rerr  :=  fh.Read(buf)
	if rerr != nil {
		return rerr
	}
	doc := string(buf)

	// xml.Decoder doesn't properly handle whitespace in some doc
	// see songTextString.xml test case ... 
	reg,_ := regexp.Compile("[ \t\n\r]*<")
	doc = reg.ReplaceAllString(doc,"<")
	b := bytes.NewBufferString(doc)
	// indirectly stop go routine if it's still running
	defer b.Reset()

	// now start reading the buffer and process messages
	mout := make(chan msg,1)
	go ReadDocsFromBuffer(b,mout,r)
	for {
		m :=<-mout
		if m.err != nil && m.err != io.EOF {
			if ok := ehandler(m.err); !ok {
				break
			 }
		}
		if m.val != nil {
			if ok := phandler(m.val); !ok {
				break
			}
		}
		if m.err == io.EOF {
			break
		}
	}
	return nil
}

// start as a go routine to feed the file results
func ReadDocsFromBuffer(b *bytes.Buffer, mout chan msg,recast ...bool) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	for {
		m := new(msg)
		m.val, m.err = XmlBufferToMap(b,r)
		if m.err != nil {
			if m.err == io.EOF {
				mout<- *m
				return
			}
			mout<- *m
		}
		mout<- *m
	}
}

// XmlBufferToMap - derived from DocToMap()
func XmlBufferToMap(b *bytes.Buffer,recast ...bool) (map[string]interface{},error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}

	n,err := BufferToTree(b)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m[n.key] = n.treeToMap(r)

	return m,nil
}

// BufferToTree - derived from DocToTree()
func BufferToTree(b *bytes.Buffer) (*Node, error) {
	p := xml.NewDecoder(b)
	n, berr := xmlToTree("",nil,p)
	if berr != nil {
		return nil, berr
	}

	return n,nil
}


