// io.Reader --> map[string]interface{} or JSON string
// nothing magic - just implements generic Go case

package x2j

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

// ToTree() - parse a XML io.Reader to a tree of Nodes
func ToTree(rdr io.Reader) (*Node, error) {
	// We need to put an *os.File reader in a ByteReader or the xml.NewDecoder
	// will wrap it in a bufio.Reader and seek on the file beyond where the
	// xml.Decoder parses!
	if _, ok := rdr.(io.ByteReader); !ok {
		rdr = myByteReader(rdr) // see code at EOF
	}

	p := xml.NewDecoder(rdr)
	p.CharsetReader = X2jCharsetReader
	n, perr := xmlToTree("", nil, p)
	if perr != nil {
		return nil, perr
	}

	return n, nil
}

// ToMap() - parse a XML io.Reader to a map[string]interface{}
func ToMap(rdr io.Reader, recast ...bool) (map[string]interface{}, error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	n, err := ToTree(rdr)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m[n.key] = n.treeToMap(r)

	return m, nil
}

// ToJson() - parse a XML io.Reader to a JSON string
func ToJson(rdr io.Reader, recast ...bool) (string, error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m, merr := ToMap(rdr, r)
	if m == nil || merr != nil {
		return "", merr
	}

	b, berr := json.Marshal(m)
	if berr != nil {
		return "", berr
	}

	return string(b), nil
}

// ToJsonIndent - the pretty form of ReaderToJson
func ToJsonIndent(rdr io.Reader, recast ...bool) (string, error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m, merr := ToMap(rdr, r)
	if m == nil || merr != nil {
		return "", merr
	}

	b, berr := json.MarshalIndent(m, "", "  ")
	if berr != nil {
		return "", berr
	}

	// NOTE: don't have to worry about safe JSON marshaling with json.Marshal, since '<' and '>" are reservedin XML.
	return string(b), nil
}

// ReaderValuesFromTagPath - io.Reader version of ValuesFromTagPath()
func ReaderValuesFromTagPath(rdr io.Reader, path string, getAttrs ...bool) ([]interface{}, error) {
	var a bool
	if len(getAttrs) == 1 {
		a = getAttrs[0]
	}
	m, err := ToMap(rdr)
	if err != nil {
		return nil, err
	}

	return ValuesFromKeyPath(m, path, a), nil
}

// ReaderValuesForTag - io.Reader version of ValuesForTag()
func ReaderValuesForTag(rdr io.Reader, tag string) ([]interface{}, error) {
	m, err := ToMap(rdr)
	if err != nil {
		return nil, err
	}

	return ValuesForKey(m, tag), nil
}

//============================ from github.com/clbanning/mxj/mxl.go ==========================

type byteReader struct {
	r io.Reader
	b []byte
}

func myByteReader(r io.Reader) io.Reader {
	b := make([]byte, 1)
	return &byteReader{r, b}
}

// need for io.Reader - but we don't use it ...
func (b *byteReader) Read(p []byte) (int, error) {
	return 0, nil
}

func (b *byteReader) ReadByte() (byte, error) {
	_, err := b.r.Read(b.b)
	if len(b.b) > 0 {
		return b.b[0], nil
	}
	var c byte
	return c, err
}
