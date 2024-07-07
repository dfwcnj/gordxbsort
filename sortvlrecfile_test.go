package main

import (
	"fmt"
	"log"
	"os"
	"path"
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

	var klns kvallines
	var tklns kvallines
	var err error
	var nr int
	dn := "/tmp"

	log.Println("sortlrecfile test")

	rsl := randomstrings(nrs, l)

	fn := path.Join(dn, "sortvlrecfiletest")
	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	for i, _ := range rsl {
		fmt.Fprintln(fp, rsl[i])
		nr++
	}

	log.Printf("sortvlrecfile %s, %s, %d\n", fn, dn, iomem)
	klns, fns, err := sortvlrecfile(fn, dn, int(l)+1, 0, 0, iomem)

	log.Println("sortvlrecfile after klns ", len(klns))
	log.Println("sortvlrecfile after fns ", len(fns))

	for _, f := range fns {
		mfp, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		tklns, _, err = vlreadn(mfp, 0, 0, 0, iomem)
		nss += uint(len(tklns))
	}
	if nrs != nss {
		log.Fatal("sortflrecfile test wanted ", nrs, " got ", nss)
	}
}
