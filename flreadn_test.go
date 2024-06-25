package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
)

//type kvalline struct {
//	key  []byte
//	line []byte
//}

func Test_flreadn(t *testing.T) {
	var l uint = 32
	var lpo uint = 1 << 16

	var klns kvallines
	var offset int64
	var err error
	var nr int

	rsl := randomstrings(lpo, l)
	log.Println("rsl ", len(rsl))

	td := os.TempDir()
	fn := path.Join(td, "rtxt.txt")

	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for i, _ := range rsl {
		fmt.Fprintln(fp, rsl[i])
		nr++
	}
	log.Println(fn, " ", nr)

	offset, err = fp.Seek(0, 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fn, " offset ", offset)

	offset, err = fp.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fn, " offset ", offset)

	//klns, offset, err = flreadn(fp, offset, int(l), 0, 0, int(lpo))
	klns, offset, err = flreadn(fp, offset, int(l)+1, 0, 0, 0)
	for _, kln := range klns {
		if len(kln.line) != int(l)+1 {
			log.Fatal("kln.line ", kln.line, " len ", len(kln.line))
		}
		if len(kln.key) != len(kln.line) {
			log.Fatal("kln.key ", kln.line, " len ", len(kln.line))
		}
		//log.Print(string(kln.line))
	}
	if len(klns) != int(lpo) {
		log.Fatal("flreadn: expected ", lpo, " got ", len(klns))
	}
	log.Print("flreadn test passed")
}
