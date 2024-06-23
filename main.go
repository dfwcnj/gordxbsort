package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type line []byte
type lines []line

//	type kvalline struct {
//		key  []byte
//		line []byte
//	}
//
// type kvallines []kvalline

// sortflrecfile(fn, dn, reclen, keyoff, keylen, lpo)
func sortflrecfile(fn string, dn string, reclen int, keyoff int, keylen int, lpo int) (kvallines, string, error) {
	var klns kvallines
	var offset int64
	var err error

	fp, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	for {
		klns, offset, err = flreadn(fp, 0, reclen, keyoff, keylen, lpo)
		klrsort2a(klns, 0)
		// call savemergefile()
		if offset == 0 || err == io.EOF {
			return klns, dn, err
		}
		if dn == "" {
			dn, err = os.MkdirTemp("", "sort")
			if err != nil {
				log.Fatal(err)
			}
		}
		mfn := savemergefile(klns, fn, dn)
		if mfn == "" {
			log.Fatal("savemergefile failed: ", fn, " ", dn)
		}

		// for debugging
		for _, l := range klns {
			fmt.Print(string(l.line))
		}

	}
	//return klns, dn, nil
}

func sortvlrecfile(fn string, dn string, reclen int, keyoff int, keylen int, lpo int) (kvallines, string, error) {
	if reclen > 0 {
		sortflrecfile(fn, dn, reclen, keyoff, keylen, lpo)
	}
	var offset int64
	var klns kvallines
	var err error

	fp, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	for {
		klns, offset, err = vlscann(fp, offset, keyoff, keylen, lpo)
		klrsort2a(klns, 0)
		if offset == 0 || err == io.EOF {
			return klns, dn, err
		}
		if dn == "" {
			dn, err = os.MkdirTemp("", "sort")
			if err != nil {
				log.Fatal(err)
			}
		}
		mfn := savemergefile(klns, fn, dn)
		if mfn == "" {
			log.Fatal("savemergefile failed: ", fn, " ", dn)
		}

		// for debugging
		for _, l := range klns {
			fmt.Print(string(l.line))
		}

	}
	//return klns, dn, nil

}

func sortfiles(fns []string, ofn string, reclen int, keyoff int, keylen int, lpo int) {

	var klns kvallines
	var err error
	if len(fns) == 0 {
		if reclen != 0 {
			klns, _, err = sortflrecfile("", "", reclen, keyoff, keylen, lpo)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			klns, _, err = sortvlrecfile("", "", reclen, keyoff, keylen, lpo)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err != nil {
			log.Fatal(err)
		}
		fp := os.Stdout
		if ofn != "" {
			fp, err := os.OpenFile(ofn, os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
			defer fp.Close()
		}
		for _, kln := range klns {

			_, err := fp.Write(kln.line)
			if err != nil {
				log.Fatal(err)
			}
		}

		return
	}
	if len(fns) == 1 {
		if reclen != 0 {
			klns, _, err = sortflrecfile(fns[0], "", reclen, keyoff, keylen, lpo)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			klns, _, err = sortvlrecfile(fns[0], "", reclen, keyoff, keylen, lpo)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}
		fp := os.Stdout
		if ofn != "" {
			fp, err := os.OpenFile(ofn, os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
			defer fp.Close()
		}
		for _, kln := range klns {

			_, err := fp.Write(kln.line)
			if err != nil {
				log.Fatal(err)
			}
		}

		return
	}
	dn, err := os.MkdirTemp("", "sort")
	if err != nil {
		log.Fatal(err)
	}
	for _, fn := range fns {
		if reclen != 0 {
			_, dn, err = sortflrecfile(fn, dn, reclen, keyoff, keylen, lpo)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, dn, err = sortvlrecfile(fn, dn, reclen, keyoff, keylen, lpo)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Fatal(err)
		}
	}
	mergefiles(dn, lpo)

}

func main() {
	var fns []string
	var ofn string
	var reclen, keylen, keyoff int
	var lpo int
	flag.StringVar(&ofn, "ofn", "", "output file name")
	flag.IntVar(&reclen, "reclen", 0, "length of the fixed length record")
	flag.IntVar(&keyoff, "keyoff", 0, "offset of the key")
	flag.IntVar(&keylen, "keylen", 0, "length of the key if not whole line")
	flag.IntVar(&lpo, "lpo", 1<<20, "lines per sort/merge")
	flag.Parse()
	fns = flag.Args()

	sortfiles(fns, ofn, reclen, keyoff, keylen, lpo)

}
