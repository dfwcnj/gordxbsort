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
	var nrs uint = 1 << 20
	var dlim string
	dlim = "\n"
	// var mrlen int

	log.Print("klchan test")

	dn, err := initmergedir("/tmp", "rdxsort")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dn)

	for i := range 1 {
		var klns kvallines
		var kln kvalline

		rsl := randomstrings(nrs, l)
		for _, s := range rsl {
			bln := []byte(s)
			kln.line = bln
			kln.key = kln.line
			klns = append(klns, kln)
		}
		if len(klns) != int(nrs) {
			//log.Print(klns)
			log.Fatal("klns: before sort wanted len ", l, " got ", len(klns))
		}

		slns := klrsort2a(klns, 0)
		var fn = filepath.Join(dn, fmt.Sprint("file", i))
		//mrlen = len(slns[0])
		fn, _ = savemergefile(slns, fn, dlim)

		inch := make(chan kvalline)

		go klchan(fn, klnullsplit, inch)

		var cklns kvallines

		for ckln := range inch {
			cklns = append(cklns, ckln)
		}
		if len(cklns) != int(nrs) {
			log.Fatal("klchan len(klns) wanted ", nrs, " got ", len(cklns))
		}
		log.Print("klchan test passed")

	}
}
