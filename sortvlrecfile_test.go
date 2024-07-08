package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"testing"
)

//type kvalline struct {
//	key  []byte
//	line []byte
//}

func Test_sortvlrecfile(t *testing.T) {
	var l uint = 32
	var nrs uint = 1 << 20
	var nss uint
	var iomem int64 = 1<<24 + 1<<20

	//var klns kvallines
	var tklns kvallines
	var err error
	var nr int

	dn, err := initmergedir("/tmp", "rdxsort")

	log.Println("sortlrecfile test")

	rsl := randomstrings(nrs, l)

	fn := path.Join(dn, "sortvlrecfiletest")
	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	for i, _ := range rsl {
		fmt.Fprintln(fp, rsl[i])
		nr++
	}
	fp.Close()

	//log.Printf("sortvlrecfile test sortvlrecfile %s, %d\n", fn, iomem)
	_, fns, err := sortvlrecfile(fn, dn, int(l)+1, 0, 0, iomem)

	//log.Println("sortvlrecfile test after  klns ", len(klns))
	//log.Println("sortvlrecfile test after fns ", fns)

	for _, f := range fns {
		//log.Println("sortvlrecfile chacking ", f)
		mfp, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		tklns, _, err = vlreadn(mfp, 0, 0, 0, iomem*2)

		var lns = make([]string, 0)
		for _, t := range tklns {
			lns = append(lns, string(t.line))
		}
		if slices.IsSorted(lns) == false {
			log.Fatal(f, " is not sorted")
		}
		nss += uint(len(tklns))
	}
	if nrs != nss {
		log.Fatal("sortflrecfile test wanted ", nrs, " got ", nss)
	}
	log.Println("sortvlrecfile passed")
}
