package main

import (
	"github.com/AlasdairF/Conv"
	"github.com/AlasdairF/Sort/Uint32Uint32"
	"github.com/GeertJohan/coop"
	e "github.com/cinnabardk/allancorfix2/internal/errors"
	"github.com/coocood/freecache"
	"github.com/gwwfps/onetime"
	mailgun "github.com/mailgun/mailgun-go"
	garbler "github.com/michaelbironneau/garbler/lib"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sfreiberg/gotwilio"
	"github.com/xxtea/xxtea-go/xxtea"
	// "github.com/ryanskidmore/GoWork"
	// "github.com/inconshreveable/go-update"
	// "https://github.com/GeertJohan/coop" ***
	// "https://github.com/haxpax/gosms"
	// "https://github.com/dgrijalva/jwt-go"
	// "https://github.com/hashicorp/memberlist" ***
	// "https://github.com/smartystreets/goconvey" //*** testing with web ui
	// "github.com/AlasdairF/BinSearch"
	// "github.com/eapache/queue"
	//"github.com/corsc/go-geohash" // Int
	// "github.com/patrick-higgins/rtreego"
	// "github.com/imdario/mergo"
	// "https://github.com/fatih/structs"
	// "https://github.com/fatih/set"
	// "https://github.com/emirpasic/gods#redblacktree" ***
	// "https://github.com/arnauddri/algorithms" ***
	// "https://github.com/h2non/imaginary"
	// "https://github.com/h2non/bimg" ***
	// "https://github.com/bluele/gcache"
	// https://github.com/golang/go/wiki/SliceTricks
	// "http://godoc.org/github.com/robfig/pathtree"
	// "github.com/AlasdairF/Custom"*/
	"github.com/apcera/termtables"
	// "github.com/h2non/bimg"
	// "github.com/hashicorp/golang-lru"
	// "github.com/flosch/trindex" Not working
	// "github.com/jamra/LevenshteinTrie"
	// "github.com/antzucaro/matchr" Algorithms
	// "http://godoc.org/github.com/zond/god/radix"
	// "github.com/Lazin/go-ngram"
	// "github.com/EricR/phrase_search"
	// "github.com/cosn/collections/tst"
	// "github.com/gnanderson/trie" **
	// "github.com/aybabtme/trie"
	// "github.com/armon/go-radix"
	// "github.com/miekg/radix"
	// "github.com/manveru/trie"
	// "github.com/IanLewis/codekata" // exercise
	// Single word search:
	// "github.com/AlasdairF/Tokenize"
	// "github.com/typeflow/typeflow-go"
	// "github.com/blevesearch/bleve/search"
	// "github.com/typerandom/cleo"
	// "github.com/sajari/fuzzy"
	// "github.com/smartystreets/mafsa"
	"runtime"
	"os"
)

var gun = mailgun.NewMailgun("aabenhave.dk", "key-2dcm94f8dttfm4k4xid3y8x2c8et-op6", "pubkey-3c1amkiwdizjhoa-cqnkrlvylm0pmiq4")

func main(){
	c := runtime.NumCPU()
	runtime.GOMAXPROCS(c)
	os.Mkdir("data", 0777)
	e.LogInit()
	testRandom()
	test := Test{}
	test.Codex()
}

func initExp() {
	e.InfoLog.Println("Init various")
	e.InfoLog.Println(conv.String(74534537787))
	test := Test{}
	//test.strings()
	test.Garbler()
	test.GeoHash()
	test.RTree()
	test.GeoHashInt()
	test.BlueMonday()
	test.FreeCache()
	test.TermTables()
	test.GZip()
	test.Coop()
}

func (_ Test) strings() {

	test := Test{}
	test.Ferret()
	test.Sort()
	test.Deaccent4()
	test.RTrie()
	test.Ngram()
	test.TST()
	test.radix()
	test.Patricia()
	test.IRadix()
}

type Test struct{}

var testData = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")


func (_ Test) sms() {
	accountSid := "ABC123..........ABC123"
	authToken := "ABC123..........ABC123"
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := "+15555555555"
	to := "+15555555555"
	message := "Welcome to gotwilio!"
	twilio.SendSMS(from, to, message, "", "")
}

func (_ Test) OneTimeCode() {
	var secret = []byte("SOME_SECRET")
	var otp, _ = onetime.Simple(8)
	var code = otp.TOTP(secret)
	e.InfoLog.Println(code)
}

func (_ Test) EncryptXXTea() {

	str := "Hello World! 你好，中国！"
	key := "1234567890"
	encrypt_data := xxtea.Encrypt([]byte(str), []byte(key))
	decrypt_data := string(xxtea.Decrypt(encrypt_data, []byte(key)))
	if str == decrypt_data {
		e.InfoLog.Println("success!")
	} else {
		e.InfoLog.Println("fail!")
	}
}

func (_ Test) FreeCache2() {
	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)

	key := []byte("abc")
	val := []byte("def")
	expire := 60 // expire in 60 seconds
	cache.Set(key, val, expire)
	got, err := cache.Get(key)
	if err != nil {
		e.InfoLog.Println(err)
	} else {
		e.InfoLog.Println(string(got))
	}
	affected := cache.Del(key)
	e.InfoLog.Println("deleted key ", affected)
	e.InfoLog.Println("entry count ", cache.EntryCount())
}


func (_ Test) Coop() {
	printFn := func() {
		e.InfoLog.Println("Hello world")
	}
	<-coop.All(printFn, printFn, printFn)
}

func (_ Test) TermTables() {
	table := termtables.CreateTable()
	table.AddHeaders("Name", "Age")
	table.AddRow("John", "30")
	table.AddRow("Sam", 18)
	table.AddRow("Julie", 20.14)

	e.InfoLog.Println("\n", table.Render())
}

func (_ Test) FreeCache() {

	freeCache := freecache.NewCache(512 * 1024 * 1024)
	val := make([]byte, 10)
	freeCache.Set([]byte("Key"), val, 0)
	_, _ = freeCache.Get([]byte("Key"))
}

func (_ Test) BlueMonday() {
	p := bluemonday.UGCPolicy()
	html := p.Sanitize(
		`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	)

	// Output:
	// <a href="http://www.google.com" rel="nofollow">Google</a>
	e.InfoLog.Println(html)
}

func (_ Test) Garbler() {
	g := &garbler.PasswordStrengthRequirements{MinimumTotalLength: 6, MaximumTotalLength: 9,
		Uppercase: 0, Digits: 4, Punctuation: 0}
	s, _ := garbler.NewPassword(g)
	e.InfoLog.Println(s)
}





func (_ Test) Sort() {
	e.InfoLog.Println("\nSort:")
	s := []sortUint32Uint32.KeyVal{
		{77, 212},
		{300, 500},
		{80, 53},
		{50, 52},
		{99, 2},
	}
	sortUint32Uint32.Asc(s)
	e.InfoLog.Println(s)
}
