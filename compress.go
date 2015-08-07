package main

import (
	e "github.com/cinnabardk/allancorfix2/internal/errors"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"bytes"
	"github.com/pierrec/lz4"
)


func (_ Test) GZip() {

	var in bytes.Buffer

	w, _ := gzip.NewWriterLevel(&in, gzip.BestSpeed)
	w.Write(testData)
	w.Close()
	out := in.Bytes()
	e.InfoLog.Println("Result of Gzip compression", len(out))

	var in2 bytes.Buffer

	f, _ := flate.NewWriter(&in2, flate.BestSpeed)
	f.Write(testData)
	f.Close()
	out = in2.Bytes()
	e.InfoLog.Println("Result of Flate compression", len(out))
}

func (_ Test) LZ4Compress() {
	var b = []byte("dette er en test af LZ4 komprimering")
	e.InfoLog.Println(string(b))
	b = testData

	var compressed = make([]byte, len(b))
	size, err := lz4.CompressBlockHC(b, compressed, 0)
	if err != nil || size == 0 {
		e.InfoLog.Println(err, size)
	}
	compressed = compressed[:size]

	var out = make([]byte, len(b))
	_, err = lz4.UncompressBlock(compressed, out, 0)
	if err != nil {
		e.InfoLog.Println(err, size)
	}

	e.InfoLog.Println("Result of LZ4 compression: ", len(b), size)
}

func (_ Test) CompressLZ4(data []byte) []byte {
	in := bytes.Buffer{}

	w := lz4.NewWriter(&in)
	w.Write(data)
	w.Close()
	return in.Bytes()
}