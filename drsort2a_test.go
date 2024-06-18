package main

import (
	"log"
	"sort"
	"testing"
)

func Test_rsort2a(t *testing.T) {

	var klns klines
	var l uint = 32
	ls := []uint{1, 2, 1<<4, 1<<8, 1<<16, 1<<20, 1<<24}

	for _, i := range ls {

		var kln kline

		log.Print("testing sort of ", i)
		rsl := randomstrings(i, l)
		if len(rsl) != int(i) {
			log.Fatal("rsl: wanted len ", i, " got ", len(rsl))
		}
		klns = klns[:0]
		for _, s := range rsl {
			bln := []byte(s)
			kline.line = bln
			kline.key = kline.line
			klns = append(klns, kline)
		}
		if len(kline) != int(i) {
			log.Print(kline)
			log.Fatal("kline: before sort wanted len ", i, " got ", len(lns))
		}
		drsort2a(klns, 0)
		if len(klns) != int(i) {
			log.Print(klns)
			log.Fatal("klns: after sort wanted len ", i, " got ", len(lns))
		}

		var ssl []string
		for _, s := range klns {
			ssl = append(ssl, string(s.line))
		}
		if len(ssl) != int(i) {
			log.Print(ssl)
			log.Fatal("ssl: wanted len ", i, " got ", len(ssl))
		}
		if !sort.StringsAreSorted(ssl) {
			t.Error("rsort2a failed for size ", i)
		} else {
			log.Print("sort test passed for ", i)
		}
	}
}
