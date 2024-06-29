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

// sortflrecfile(fn, dn, reclen, keyoff, keylen, lpo)
func sortflrecfile(fn string, dn string, reclen int, keyoff int, keylen int, lpo int) (kvallines, string, error) {
	var klns kvallines
	var offset int64
	var err error
	var i int

	fp := os.Stdin
	if fn != "" {
		fp, err = os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {

		klns, offset, err = flreadn(fp, offset, reclen, keyoff, keylen, lpo)
		sklns := klrsort2a(klns, 0)

		if err == io.EOF {
			if offset != 0 {
				break
			}
			return sklns, dn, err
		}

		// probably sorting a singleton
		if dn == "" {
			dn, err = initmergedir("rdxsort")
			if err != nil {
				log.Fatal(err)
			}
		}
		mfn := filepath.Join(dn, fmt.Sprintf("%s%d", fn, i))

		if savemergefile(sklns, mfn) == "" {
			log.Fatal("savemergefile failed: ", fn, " ", dn)
		}
		i++

	}
	//clear(klns)
	//mergefiles("", dn, lpo)
	return klns, dn, err
}

// sort variable lengh records file
func sortvlrecfile(fn string, dn string, reclen int, keyoff int, keylen int, lpo int) (kvallines, string, error) {
	var offset int64
	var klns kvallines
	var err error
	var i int

	fp := os.Stdin
	if fn != "" {
		fp, err = os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		klns, offset, err = vlscann(fp, offset, keyoff, keylen, lpo)
		sklns := klrsort2a(klns, 0)

		if err == io.EOF {
			if offset != 0 {
				break
			}
			return sklns, dn, err
		}

		// probably sorting a singleton
		if dn == "" {
			dn, err = initmergedir("rdxsort")
			if err != nil {
				log.Fatal(err)
			}
		}
		mfn := filepath.Join(dn, fmt.Sprintf("%s%d", fn, i))

		if savemergefile(sklns, mfn) == "" {
			log.Fatal("savemergefile failed: ", fn, " ", dn)
		}
		i++

	}
	//clear(klns)
	//mergefiles("", dn, lpo)
	return klns, dn, err
}

func sortfiles(fns []string, ofn string, reclen int, keyoff int, keylen int, lpo int) {

	var klns kvallines
	var dn string
	var err error
	log.Printf("sortfiles reclen %d\n", reclen)

	if len(fns) == 0 {
		if reclen != 0 {
			klns, dn, err = sortflrecfile("", "", reclen, keyoff, keylen, lpo)
			if dn != "" {
				mergefiles(ofn, dn, lpo)
				return
			}
		} else {
			klns, dn, err = sortvlrecfile("", "", reclen, keyoff, keylen, lpo)
			if dn != "" {
				if err != io.EOF {
					log.Fatal(err)
				}
				mergefiles(ofn, dn, lpo)
				return
			}
		}

		// not enough recs to call mergefiles
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
		log.Printf("sortfiles fn %s\n", fns[0])
		if reclen != 0 {
			klns, _, err = sortflrecfile(fns[0], "", reclen, keyoff, keylen, lpo)
			if dn != "" {
				mergefiles(ofn, dn, lpo)
				return
			}
		} else {
			klns, _, err = sortvlrecfile(fns[0], "", reclen, keyoff, keylen, lpo)
			if dn != "" {
				mergefiles(ofn, dn, lpo)
				return
			}
		}

		// not enough recs to call mergefiles
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

	dn, err = initmergedir("rdxsort")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dn)

	for n, fn := range fns {
		log.Printf("sortfiles fn %s\n", fn)
		var klns kvallines
		var d string
		var err error
		if reclen != 0 {
			klns, d, err = sortflrecfile(fn, dn, reclen, keyoff, keylen, lpo)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			// merge files already saved
			if d != "" {
				continue
			}
		} else {
			klns, d, err = sortvlrecfile(fn, dn, reclen, keyoff, keylen, lpo)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			// merge files already saved
			if d != "" {
				continue
			}
		}

		mfn := fmt.Sprintf("%s%d", filepath.Base(fn), n)
		mpath := filepath.Join(dn, mfn)
		savemergefile(klns, mpath)
	}
	mergefiles(ofn, dn, lpo)
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
