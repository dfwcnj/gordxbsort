package main

import (
	"log"
	"os"
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
	var l int = 32
	var r bool = true
	//ls := []int{1, 2, 1 << 4, 1 << 8, 1 << 16, 1 << 20, 1 << 24}
	ls := []int{1 << 4, 1 << 10, 1 << 20, 1 << 24}

	for _, nl := range ls {

		log.Print("klrsort2a test ", nl)
		var klns kvallines

		//log.Print("testing sort of ", nl)
		rsl := randomstrings(nl, l, r)
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
			log.Fatal("klrsort2a test ssl: wanted len ", nl, " got ", len(ssl))
		}
		if !sort.StringsAreSorted(ssl) {
			fp, err := os.OpenFile("/tmp/klrsort2atest", os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
			for _, l := range ssl {
				l = l + "\n"
				fp.Write([]byte(l))
			}
			fp.Close()
			log.Fatal("klrrsort2a test not in sort order")
		} else {
			log.Print("klrsort2a test passed")
		}
	}
}
