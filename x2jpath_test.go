package x2j

/*
import (
	"fmt"
	"testing"
)

// the basic demo/test case - a small bibliography with mixed element types
func TestValuesFromTagPath(t *testing.T) {
var doc = `
<doc>
	<books>
		<book>
			<author>William H. Gaddis</author>
			<title>The Recognitions</title>
			<review>One of the great seminal American novels of the 20th century.</review>
		</book>
		<book>
			<author>Austin Tappan Wright</author>
			<title>Islandia</title>
			<review>An example of earlier 20th century American utopian fiction.</review>
		</book>
		<book>
			<author>John Hawkes</author>
			<title>The Beetle Leg</title>
			<review>A lyrical novel about the construction of Ft. Peck Dam in Montana.</review>
		</book>
		<book>
			<author>
				<first_name>T.E.</first_name>
				<last_name>Porter</last_name>
			</author>
			<title>King's Day</title>
			<review>A magical novella.</review>
		</book>
	</books>
</doc>
`
	fmt.Println("\nTestValuesFromTagPath()\n",doc)

	m,_ := DocToMap(doc)
	fmt.Println("map:",WriteMap(m))

	v,_ := ValuesFromTagPath(doc,"doc.books")
	fmt.Println("\npath == doc.books: len(v):",len(v))
	for key,val := range v {
		fmt.Println(key,":",val)
	}

	v,_ = ValuesFromTagPath(doc,"doc.books.*")
	fmt.Println("\npath == doc.books.*: len(v):",len(v))
	for key,val := range v {
		fmt.Println(key,":",val)
	}

	v,_ = ValuesFromTagPath(doc,"doc.books.*.author")
	fmt.Println("\npath == doc.books.*.author: len(v):",len(v))
	for key,val := range v {
		fmt.Println(key,":",val)
	}

	v,_ = ValuesFromTagPath(doc,"doc.*.*.author")
	fmt.Println("\npath == doc.books.*.*.author: len(v):",len(v))
	for key,val := range v {
		fmt.Println(key,":",val)
	}

	v,_ = ValuesFromTagPath(doc,"doc.*.*.title")
	fmt.Println("\npath == doc.books.*.*.title: len(v):",len(v))
	for key,val := range v {
		fmt.Println(key,":",val)
	}

	v,_ = ValuesFromTagPath(doc,"doc.*.*.*")
	fmt.Println("\npath == doc.books.*.*.*: len(v):",len(v))
	for key,val := range v {
		fmt.Println(key,":",val)
	}
}

// demo how to compensate for irregular tag labels in data
// "netid" vs. "idnet"
func TestValuesFromTagPath2(t *testing.T) {
var doc1 = `
<?xml version="1.0" encoding="UTF-8"?>
<data>
    <netid>
        <disable>no</disable>
        <text1>default:text</text1>
        <word1>default:word</word1>
    </netid>
</data>
`
var doc2 = `
<?xml version="1.0" encoding="UTF-8"?>
<data>
    <idnet>
        <disable>yes</disable>
        <text1>default:text</text1>
        <word1>default:word</word1>
    </idnet>
</data>
`
	var docs = []string{doc1,doc2}

	for n,doc := range docs {
		fmt.Println("\nTestValuesFromTagPath2(), iteration:",n,"\n",doc)

		m,_ := DocToMap(doc)
		fmt.Println("map:",WriteMap(m))

		v,_ := ValuesFromTagPath(doc,"data.*")
		fmt.Println("\npath == data.*: len(v):",len(v))
		for key,val := range v {
			fmt.Println(key,":",val)
		}

		v,_ = ValuesFromTagPath(doc,"data.*.*")
		fmt.Println("\npath == data.*.*: len(v):",len(v))
		for key,val := range v {
			fmt.Println(key,":",val)
		}
	}
}
*/
