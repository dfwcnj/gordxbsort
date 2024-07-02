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

func Test_sortfiles(t *testing.T) {
	var l uint = 32
	var lpo uint = 1 << 20
	var iomem int64 = 1 << 30
	var nmf = 10
	var tmpdir = "/tmp"

	log.Print("sortfiles test")

	var fns []string
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
			log.Fatal("sortfiles test before sort wanted len ", l, " got ", len(klns))
		}

		//log.Println("sorting file", i)
		slns := klrsort2a(klns, 0)
		var fn = filepath.Join(tmpdir, fmt.Sprint("file", i))
		//log.Println("saving file", i)
		savemergefile(slns, fn)
		fns = append(fns, fn)
	}

	mfn := "mergeout.txt"
	mpath := filepath.Join(tmpdir, mfn)

	log.Println("sorting files to ", mpath)
	sortfiles(fns, mpath, 0, 0, 0, 0, iomem)

	mfp, err := os.Open(mpath)
	if err != nil {
		log.Fatal("sortfiles test ", err)
	}
	defer mfp.Close()

	scanner := bufio.NewScanner(mfp)
	var mlns []string
	for scanner.Scan() {
		l := scanner.Text()
		mlns = append(mlns, l)
	}
	if len(mlns) != int(lpo)*nmf {
		log.Fatal("sortfiles test n wanted ", int(lpo)*nmf, " got ", len(mlns))
	}
	if !slices.IsSorted(mlns) {
		log.Fatal("sortfiles test lines in ", mfn, " not in sort order")
	}
	log.Print("sortfiles test passed")

}
