package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type line []byte
type lines []line

// sortflrecfile(fn, dn, reclen, keyoff, keylen, lpo)
func sortflrecfile(fn string, dn string, reclen int, keyoff int, keylen int, lpo int, iomem int64) (kvallines, string, error) {
	var klns kvallines
	var offset int64
	var err error
	var i int

	fp := os.Stdin
	if fn != "" {
		fp, err = os.Open(fn)
		if err != nil {
			log.Fatal("sortflrecfile ", err)
		}
	}

	log.Println("sortflrecfile")

	for {

		klns, offset, err = flreadn(fp, offset, reclen, keyoff, keylen, lpo, iomem)
		sklns := klrsort2a(klns, 0)
		//inch := make(chan kvallines, 0)
		//go cklrsort2a(klns, 0, inch)
		//sklns := <-inch

		if offset == 0 {
			return sklns, "", err
		}

		// probably sorting a singleton
		if dn == "" {
			dn, err = initmergedir("rdxsort")
			if err != nil {
				log.Fatal(err)
			}
		}
		mfn := filepath.Join(dn, filepath.Base(fmt.Sprintf("%s%d", fn, i)))

		if savemergefile(sklns, mfn) == "" {
			log.Fatal("savemergefile failed: ", fn, " ", dn)
		}
		i++

	}
	//clear(klns)
	//mergefiles("", dn, lpo)
	// return klns, dn, err
}

// sort variable lengh records file
func sortvlrecfile(fn string, dn string, reclen int, keyoff int, keylen int, lpo int, iomem int64) (kvallines, string, error) {
	var offset int64
	var klns kvallines
	var err error
	var i int

	fp := os.Stdin
	if fn != "" {
		fp, err = os.Open(fn)
		if err != nil {
			log.Fatal("sortvlrecfile ", err)
		}
	}

	for {
		klns, offset, err = vlscann(fp, offset, keyoff, keylen, lpo, iomem)

		if err != nil {
			log.Fatal("sortvlrecfile after vlscann ", fn, " ", err)
		}

		sklns := klrsort2a(klns, 0)
		// inch := make(chan kvallines, 0)
		// go cklrsort2a(klns, 0, inch)
		// sklns := <-inch

		if offset == 0 {
			return sklns, "", err
		}

		// probably sorting a singleton
		if dn == "" {
			dn, err = initmergedir("rdxsort")
			if err != nil {
				log.Fatal(err)
			}
		}
		mfn := filepath.Join(dn, filepath.Base(fmt.Sprintf("%s%d", fn, i)))
		if savemergefile(sklns, mfn) == "" {
			log.Fatal("savemergefile failed: ", fn, " ", dn)
		}
		i++

	}
}

func sortfiles(fns []string, ofn string, reclen int, keyoff int, keylen int, lpo int, iomem int64) {

	var klns kvallines
	var dn string
	var err error
	// log.Printf("sortfiles ofn %s\n", ofn)

	if len(fns) == 0 {
		if reclen != 0 {
			klns, dn, err = sortflrecfile("", "", reclen, keyoff, keylen, lpo, iomem)
			if dn != "" {
				mergefiles(ofn, dn, lpo)
				return
			}
		} else {
			klns, dn, err = sortvlrecfile("", "", reclen, keyoff, keylen, lpo, iomem)
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
			klns, _, err = sortflrecfile(fns[0], "", reclen, keyoff, keylen, lpo, iomem)
			if dn != "" {
				mergefiles(ofn, dn, lpo)
				return
			}
		} else {
			klns, _, err = sortvlrecfile(fns[0], "", reclen, keyoff, keylen, lpo, iomem)
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
	// log.Print("sortfiles dn ", dn)

	for _, fn := range fns {
		var klns kvallines
		var d string
		if reclen != 0 {
			log.Println("sortfiles fn reclen ", fn, " ", reclen)
			klns, d, err = sortflrecfile(fn, dn, reclen, keyoff, keylen, lpo, iomem)
			if d != "" {
				continue
			}
		} else {
			klns, d, err = sortvlrecfile(fn, dn, reclen, keyoff, keylen, lpo, iomem)
			if d != "" {
				continue
			}
		}

		mfn := fmt.Sprintf("%s", filepath.Base(fn))
		mpath := filepath.Join(dn, mfn)
		savemergefile(klns, mpath)
	}
	mergefiles(ofn, dn, lpo)
}

func parseiomem(iomem string) int64 {

	ns := iomem[0 : len(iomem)-2]
	n, err := strconv.ParseInt(ns, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	ms := iomem[len(iomem)-2:]
	switch ms {
	case "kb":
		return n * 1 << 10
	case "mb":
		return n * 1 << 20
	case "gb":
		return n * 1 << 30
	default:
		log.Fatal("bad iomem argument: ", iomem)
	}
	return 0
}

func main() {
	var fns []string
	var ofn, iomem string
	var reclen, keylen, keyoff int
	var lpo int
	flag.StringVar(&ofn, "ofn", "", "output file name")
	flag.StringVar(&iomem, "iomem", "500mb", "max read memory size in kb, mb or gb")
	flag.IntVar(&reclen, "reclen", 0, "length of the fixed length record")
	flag.IntVar(&keyoff, "keyoff", 0, "offset of the key")
	flag.IntVar(&keylen, "keylen", 0, "length of the key if not whole line")
	flag.IntVar(&lpo, "lpo", 0, "lines per sort/merge")
	flag.Parse()
	fns = flag.Args()

	var iom int64
	if iomem != "" {
		iom = parseiomem(iomem)
	}
	sortfiles(fns, ofn, reclen, keyoff, keylen, lpo, iom)

}
