package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type line []byte
type lines []line

//	type kvalline struct {
//		key  []byte
//		line []byte
//	}
//
// type kvallines []kvalline

func mergefiles(dn string, lpo int) {
}

// save merge file
// save key and line separated by null bute
func savemergefile(klns kvallines, fn string, dn string) string {
	bn := filepath.Base(fn)
	pfn := filepath.Join(dn, bn)
	fp, err := os.OpenFile(pfn, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	for _, kln := range klns {

		var n = byte(0)
		knl := string(kln.key) + string(n) + string(kln.line)

		_, err := fp.Write([]byte(knl))
		if err != nil {
			log.Fatal(err)
		}
	}

	return fn
}

// bufSplit(buf, reclen)
//
// split the buffer into a slice containing reclen records
func bufSplit(buf []byte, reclen int) lines {
	buflen := len(buf)
	var lns lines
	for o := 0; o < buflen; o += reclen {
		rec := buf[o : o+reclen-1]
		lns = append(lns, rec)
	}
	return lns
}

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
	return klns, dn, nil
}

func sortfile(fn string, dn string, reclen int, keyoff int, keylen int, lpo int) (kvallines, string, error) {
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
	return klns, dn, nil

}

func sortfiles(fns []string, ofn string, reclen int, keyoff int, keylen int, lpo int) {

	if len(fns) == 0 {
		klns, _, err := sortfile("", "", reclen, keyoff, keylen, lpo)
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
		klns, _, err := sortfile(fns[0], "", reclen, keyoff, keylen, lpo)
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
		_, dn, err = sortfile(fn, dn, reclen, keyoff, keylen, lpo)
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
