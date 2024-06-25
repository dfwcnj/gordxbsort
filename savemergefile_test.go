package main

import (
	"fmt"
	"log"
	"testing"
)

//type kvalline struct {
//	key  []byte
//	line []byte
//}

func Test_savemergefile(t *testing.T) {
	var l uint = 32
	var lpo uint = 1 << 8
	var dn = "/private/tmp"

	for i := range 10 {
		var klns kvallines
		var kln kvalline

		rsl := randomstrings(lpo, l)
		for _, s := range rsl {
			bln := []byte(s)
			kln.line = bln
			kln.key = kln.line
			klns = append(klns, kln)
		}
		if len(klns) != int(lpo) {
			//log.Print(klns)
			log.Fatal("klns: before sort wanted len ", l, " got ", len(klns))
		}

		slns := klrsort2a(klns, 0)
		var fn = fmt.Sprint("file", i)
		savemergefile(slns, fn, dn)
	}
}
