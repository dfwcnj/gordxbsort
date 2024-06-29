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
	var lpo uint = 1 << 20
	var nmf = 10

	log.Print("mergefiles test")

	dn, err := initmergedir("rdxsort")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dn)

	for i := range nmf {
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
			log.Fatal("klns: before sort wanted len ", l, " got ", len(klns))
		}

		slns := klrsort2a(klns, 0)
		var fn = filepath.Join(dn, fmt.Sprint("file", i))
		savemergefile(slns, fn)
	}

	mfn := "mergeout.txt"
	mpath := filepath.Join(dn, mfn)
	mergefiles(mpath, dn, int(lpo))

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
	if len(mlns) != int(lpo)*nmf {
		log.Fatal("mergefiles n wanted ", int(lpo)*nmf, " got ", len(mlns))
	}
	if !slices.IsSorted(mlns) {
		log.Fatal("lines in ", mfn, " not in sort order")
	}
	log.Print("mergefiles test passed")

}
