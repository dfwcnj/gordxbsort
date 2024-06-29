package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func Test_klchan(t *testing.T) {
	var l uint = 32
	var lpo uint = 1 << 16

	log.Print("klchan test")

	dn, err := initmergedir("rdxsort")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dn)

	for i := range 1 {
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
		var fn = filepath.Join(dn, fmt.Sprint("file", i))
		savemergefile(slns, fn)

		inch := make(chan kvalline)

		go klchan(fn, klnullsplit, inch)

		var cklns kvallines

		for ckln := range inch {
			cklns = append(cklns, ckln)
		}
		if len(cklns) != int(lpo) {
			log.Fatal("klchan len(klns) wanted ", lpo, " got ", len(cklns))
		}
		log.Print("klchan test passed")

	}
}
