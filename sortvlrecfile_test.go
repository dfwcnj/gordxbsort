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
	var l int = 32
	var r bool = true
	var nrs int = 1 << 20
	var nss int
	var iomem int64 = 1<<24 + 1<<20

	//var klns kvallines
	var tklns kvallines
	var err error
	var nr int

	dn, err := initmergedir("/tmp", "rdxsort")

	log.Println("sortlrecfile test")

	rsl := randomstrings(nrs, l, r)

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
		//log.Println("sortvlrecfile test tklns ", len(tklns))

		var lns = make([]string, 0)
		for _, t := range tklns {
			lns = append(lns, string(t.line))
		}
		if slices.IsSorted(lns) == false {
			log.Fatal(f, " is not sorted")
		}
		nss += int(len(tklns))
	}
	if nrs != nss {
		log.Fatal("sortvlrecfile test wanted ", nrs, " got ", nss)
	}
	log.Println("sortvlrecfile passed")
}
