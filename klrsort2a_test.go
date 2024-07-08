package main

import (
	"log"
	"sort"
	"testing"
)

//type kvalline struct {
//	key  []byte
//	line []byte
//}
//type kvallines []kvalline

func Test_klrsort2a(t *testing.T) {

	log.Print("klrsort2a test")
	var l uint = 32
	//ls := []uint{1, 2, 1 << 4, 1 << 8, 1 << 16, 1 << 20, 1 << 24}
	ls := []uint{1 << 24}

	for _, nl := range ls {

		var klns kvallines

		//log.Print("testing sort of ", nl)
		rsl := randomstrings(nl, l)
		if len(rsl) != int(nl) {
			log.Fatal("rsl: wanted len ", nl, " got ", len(rsl))
		}
		for _, s := range rsl {
			var kln kvalline
			bln := []byte(s)
			kln.line = bln
			kln.key = kln.line
			klns = append(klns, kln)
		}
		if len(klns) != int(nl) {
			log.Fatal("klns: before sort wanted len ", nl, " got ", len(klns))
		}
		slns := klrsort2a(klns, 0)
		if len(slns) != int(nl) {
			log.Fatal("slns: after sort wanted len ", nl, " got ", len(slns))
		}

		var ssl []string
		for _, s := range slns {
			ssl = append(ssl, string(s.line))
		}
		if len(ssl) != 1 && ssl[0] == ssl[len(ssl)-1] {
			log.Fatal("strings are all equal")
		}
		if len(ssl) != int(nl) {
			log.Fatal("ssl: wanted len ", nl, " got ", len(ssl))
		}
		for i, _ := range ssl {
			if len(ssl[i]) != int(l) {
				log.Fatal("ssl[i]: wanted len ", l, " got ", len(ssl[i]))
			}
		}
		if !sort.StringsAreSorted(ssl) {
			t.Error("rsort2a failed")
		} else {
			log.Print("klrsort2a test passed")
		}
	}
}
