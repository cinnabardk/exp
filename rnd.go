package main

import (
	e "github.com/cinnabardk/allancorfix2/internal/errors"
	"github.com/losalamos/rdrand"
	"github.com/AndreasBriese/breeze"
	"github.com/db47h/rand64"
	//"github.com/tildeleb/cuckoo"
	"github.com/db47h/rand64/xorshift"
	"github.com/zephyrtronium/xer"
	"math/rand"
	"time"
)

func testRandom(){
	test := Test{}
	test.Breeze()
	test.Rand64()
	test.Xer()
}


func (_ Test) Random() {
	r := rdrand.Uint64()
	println(r)
}


func (_ Test)Breeze() {
	//
	// drain 1000 random bytes from chaos
	//
	var bmap128 breeze.BreezeCS128
	err := bmap128.Init()
	if err != nil {
		e.InfoLog.Println(err)
		panic(1)
	}
	resultI := make([]uint8, 70)
	for i, _ := range resultI {
		resultI[i] = uint8(bmap128.RandIntn()) >> 4
	}
	e.InfoLog.Println("Breeze Random: ", resultI)
}
func (_ Test)Rand64(){
	const 		SEED1 = 1387366483214
	var s rand64.Source64 = xorshift.New64star()
	s.Seed64(SEED1)

		r64 := rand64.New(s)

	resultI := make([]uint8, 70)
	for i, _ := range resultI {
		resultI[i] = uint8(r64.Uint32()) >> 4
	}
	e.InfoLog.Println("Rand64 Random: ", resultI)
}

func (_ Test) Xer(){
	var rng *rand.Rand
	rng = rand.New(xer.New(time.Now().UnixNano(), 256))

	resultI := make([]uint8, 70)
	for i, _ := range resultI {
		resultI[i] = uint8(rng.Uint32()) >> 4
	}
	e.InfoLog.Println("Xer Random: ", resultI)
}