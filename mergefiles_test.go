package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func Test_mergefiles(t *testing.T) {
	var l uint = 32
	var nrs uint = 1 << 20
	var nmf = 10
	var fns []string

	log.Print("mergefiles test")

	dn, err := initmergedir("rdxsort")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dn)

	for i := range nmf {
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
			log.Fatal("klns: before sort wanted len ", l, " got ", len(klns))
		}

		slns := klrsort2a(klns, 0)
		var fn = filepath.Join(dn, fmt.Sprint("file", i))
		savemergefile(slns, fn)
		fns = append(fns, fn)
	}

	mfn := "mergeout.txt"
	mpath := filepath.Join(dn, mfn)
	mergefiles(mpath, fns)

	mfp, err := os.Open(mpath)
	if err != nil {
		log.Fatal(err)
	}
	defer mfp.Close()

	scanner := bufio.NewScanner(mfp)
	var mlns []string
	for scanner.Scan() {
		l := scanner.Text()
		mlns = append(mlns, l)
	}
	if len(mlns) != int(nrs)*nmf {
		log.Fatal(mpath, " wanted ", int(nrs)*nmf, " got ", len(mlns))
	}
	if !slices.IsSorted(mlns) {
		log.Fatal("lines in ", mfn, " not in sort order")
	}
	log.Print("mergefiles test passed")

}
