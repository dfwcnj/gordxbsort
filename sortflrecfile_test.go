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

func Test_sortflrecfile(t *testing.T) {
	var l int = 32
	var r bool = true
	var nrs int = 1 << 20
	var iomem int64 = 1<<24 + 1<<20
	var mrlen int

	var tklns kvallines
	var err error
	var nr int

	dn, err := initmergedir("/tmp", "rdxsort")

	log.Println("sortflrecfile test")

	rsl := randomstrings(nrs, l, r)

	fn := path.Join(dn, "sortflrecfiletest")
	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for i, _ := range rsl {
		fmt.Fprint(fp, rsl[i])
		nr++
	}
	fp.Close()

	_, fns, mrlen, err := sortflrecfile(fn, dn, int(l), 0, 0, iomem)

	var nss int
	for _, f := range fns {
		mfp, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		finf, err := mfp.Stat()
		if err != nil {
			log.Fatal("sortflrecfiletest ", err)
		}
		tklns, _, err = flreadn(mfp, 0, mrlen, 0, 0, finf.Size())
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
		log.Fatal("sortflrecfile test wanted ", nrs, " got ", nss)
	}
	log.Println("sortflrecfile passed")

}
