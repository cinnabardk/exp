package main

import (
	e "github.com/cinnabardk/allancorfix2/internal/errors"
	"github.com/Lazin/go-ngram"
	"github.com/argusdusty/Ferret"
	"github.com/hashicorp/go-immutable-radix"
	"github.com/cosn/collections/tst"
	"github.com/derekparker/trie"
	"github.com/tchap/go-patricia/patricia"
	"github.com/sauerbraten/radix"
	"strings"
	"github.com/Mitranim/codex"
	// 	"github.com/m4rw3r/uuid"
)

var ExampleWords = []string{
	"coca cola",
	"coce med is",
	"koala",
	"koalant afregning",
	"regning med rykker",
	"det regner",
	"regn og blest",
	"borte med blesten",
	"blost og uvejr",
	"blost og storm",
}
var ExampleUint32 = []uint32{66, 77, 88, 11, 22, 33, 44, 55, 99, 111}


// No delete
func (_ Test) Ferret() {
	e.InfoLog.Println("\nFerret:")
	var Correction = func(b []byte) [][]byte { return ferret.ErrorCorrect(b, ferret.LowercaseLetters) }
	var Converter = func(s string) []byte { return []byte(s) }
	var Data []interface{}

	for i := range ExampleUint32 {
		Data = append(Data, ExampleUint32[i])
	}
	f := ferret.New(ExampleWords, ExampleWords, Data, Converter)
	f.Insert("blest i most", "blest i most", 7)

	e.InfoLog.Println(f.ErrorCorrectingQuery("blest", 4, Correction))
	e.InfoLog.Println(f.Query("blest", 10))
}


// http://hackthology.com/ternary-search-tries-for-fast-flexible-string-search-part-1.html
func (_ Test) RTrie() {
	e.InfoLog.Println("\nRTrie:")
	t := trie.New()

	for i := range ExampleWords {
		t.Add(ExampleWords[i], ExampleUint32[i])
	}
	t.Remove("coca cola")

	node, _ := t.Find("blost og storm")
	e.InfoLog.Println(node.Meta().(uint32))
	e.InfoLog.Println(t.PrefixSearch("blost"))
	e.InfoLog.Println(t.FuzzySearch("blest"))

}

// No delete
func (_ Test) Ngram() {
	e.InfoLog.Println("\nN-Gram:")
	index, _ := ngram.NewNGramIndex(ngram.SetN(3))

	var token ngram.TokenID
	for i := range ExampleWords {
		token, _ = index.Add(ExampleWords[i])
	}

	str, _ := index.GetString(token) // str == "hello"
	e.InfoLog.Println(str)
	resultsList, _ := index.Search("blest")
	for _, v := range resultsList {
		e.InfoLog.Println(v.TokenID)
	}

}
func (_ Test) TST() {
	e.InfoLog.Println("\nTST:")
	n := tst.T{}

	for i := range ExampleWords {
		n.Insert(ExampleWords[i], nil)
	}
	n.Delete("coca cola")

	result := n.StartsWith("blost") //buggy, try searching for blest
	for _, v := range result {
		e.InfoLog.Println(v)
	}
}

func (_ Test) radix() {
	e.InfoLog.Println("\nRadix:")
	r := radix.New()

	for i := range ExampleWords {
		r.Set(ExampleWords[i], ExampleUint32[i])
	}
	x := r.GetAllWithPrefix("blost")
	for _, v := range x {
		i := v.(uint32)
		e.InfoLog.Println(i)
	}
}

func (_ Test) Patricia() {
	e.InfoLog.Println("\nPatricia:")
	printItem := func(prefix patricia.Prefix, item patricia.Item) error {
		e.InfoLog.Println(string(prefix), item.(uint32))
		return nil
	}
	trie := patricia.NewTrie()

	for i := range ExampleWords {
		trie.Insert(patricia.Prefix(ExampleWords[i]), ExampleUint32[i])
	}
	trie.Set(patricia.Prefix("coca cola"), 188)

	e.InfoLog.Println("SubTree:")
	trie.VisitSubtree(patricia.Prefix("blost"), printItem)
	e.InfoLog.Println("Prefixes:")
	trie.VisitPrefixes(patricia.Prefix("borte med blesten mega"), printItem)

	trie.Delete(patricia.Prefix("coca cola"))
	trie.DeleteSubtree(patricia.Prefix("blost"))

	e.InfoLog.Println("What is left:")
	trie.Visit(printItem)
}

func (_ Test) IRadix() {
	e.InfoLog.Println("\nIradix:")
	out := []string{}

	fn := func(k []byte, v interface{}) bool {
		out = append(out, string(k))
		return false
	}
	r := iradix.New()
	for i := range ExampleWords {
		r, _, _ = r.Insert([]byte(ExampleWords[i]), ExampleUint32[i])
	}
	x := r.Root()
	m, _, _ := x.LongestPrefix([]byte("borte med blesten mega"))
	e.InfoLog.Println(string(m))

	x.WalkPrefix([]byte("blost"), fn)

	for i := range out {
		e.InfoLog.Println(out[i])
	}
}


func (_ Test) Deaccent4() {
	e.InfoLog.Println("\nDeaccent:")
	str := "æblekage med ål og øl"

	e.InfoLog.Println(string(StripÆØÅ(str)))
}

func StripÆØÅ(str string) []byte {
	str = strings.ToLower(str)
	b := (make([]byte, len(str)))

	pos := 0
	for _, char := range str {

		if char > 127 {
			switch char {
			case 'å':
				b[pos] = 'a'
			case 'æ':
				b[pos] = 'a'
			case 'ø':
				b[pos] = 'o'
			}
		} else {
			b[pos] = byte(char)
		}
		pos += 1
	}
	return b[:pos]
}

func (_ Test) Codex(){
	source := []string{ "ral", "mag", "dyn", "kunda", "sim"}

	traits, err := codex.NewTraits(source)
	if err != nil {
		panic(err)
	}
	gen := traits.Generator()

	var str [200]string
	// Print twelve random words.
	for i := 0; i < 200; i++ {
		str[i] = gen()
	}
	e.InfoLog.Println("Codex, Words generated: ", str)

	// Find out how many words can be generated from this sample.
	gen = traits.Generator()
	i := 0
	for gen() != "" {
		i++
	}
	e.InfoLog.Println("total:", i)
}