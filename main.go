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

func sortflrecfile(fn string, dn string, reclen int, keyoff int, keylen int, iomem int64) (kvallines, []string, error) {
	var klns kvallines
	var offset int64
	var err error
	var i int
	var mfiles []string

	fp := os.Stdin
	if fn != "" {
		fp, err = os.Open(fn)
		if err != nil {
			log.Fatal("sortflrecfile ", err)
		}
	}
	if dn == "" {
		dn, err = initmergedir("rdxsort")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("sortflrecfile")

	for {

		log.Println("sortflrecfile flreadn ", offset)
		klns, offset, err = flreadn(fp, offset, reclen, keyoff, keylen, iomem)
		if len(klns) == 0 {
			log.Println("sortflrecfile klns ", len(klns))
			return klns, mfiles, err
		}

		sklns := klrsort2a(klns, 0)

		if offset > 0 && len(sklns) > 0 {
			mfn := filepath.Join(dn, filepath.Base(fmt.Sprintf("%s%d", fn, i)))
			log.Println("sortflrecfile mfn ", mfn)
			if savemergefile(sklns, mfn) == "" {
				log.Fatal("savemergefile failed: ", fn, " ", dn)
			}
			mfiles = append(mfiles, mfn)
		}
		if err == io.EOF {
			log.Println("sortflrecfile eof")
			return sklns, mfiles, nil
		}

		i++

	}
}

// sort variable lengh records file
func sortvlrecfile(fn string, dn string, reclen int, keyoff int, keylen int, iomem int64) (kvallines, []string, error) {
	var offset int64
	var klns kvallines
	var err error
	var i int
	var mfiles []string

	fp := os.Stdin
	if fn != "" {
		fp, err = os.Open(fn)
		if err != nil {
			log.Fatal("sortvlrecfile ", err)
		}
	}
	if dn == "" {
		dn, err = initmergedir("rdxsort")
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		//log.Println("sortvlrecfile vlreadn ", fn, " ", offset, " ", iomem)
		klns, offset, err = vlreadn(fp, offset, keyoff, keylen, iomem)

		if err != nil {
			log.Fatal("sortvlrecfile after vlreadn ", fn, " ", err)
		}
		if len(klns) == 0 {
			return klns, mfiles, err
		}

		sklns := klrsort2a(klns, 0)

		if offset == 0 {
			return sklns, mfiles, err
		}
		if offset > 0 && len(sklns) > 0 {
			mfn := filepath.Join(dn, filepath.Base(fmt.Sprintf("%s%d", fn, i)))
			log.Println("sortvlrecfile mfn ", mfn)
			if savemergefile(sklns, mfn) == "" {
				log.Fatal("savemergefile failed: ", fn, " ", dn)
			}
			mfiles = append(mfiles, mfn)
		}
		i++

	}
}

func sortfiles(fns []string, ofn string, reclen int, keyoff int, keylen int, iomem int64) {

	var klns kvallines
	var dn string
	var err error
	var mfiles []string
	// log.Printf("sortfiles ofn %s\n", ofn)

	fp := os.Stdout
	if ofn != "" {
		fp, err := os.OpenFile(ofn, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()
	}

	if len(fns) == 0 {
		if reclen != 0 {
			klns, mfiles, err = sortflrecfile("", "", reclen, keyoff, keylen, iomem)
		} else {
			klns, mfiles, err = sortvlrecfile("", "", reclen, keyoff, keylen, iomem)
		}
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if len(mfiles) > 0 {
			mergefiles(ofn, mfiles)
			return
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
			klns, mfiles, err = sortflrecfile(fns[0], "", reclen, keyoff, keylen, iomem)
		} else {
			klns, mfiles, err = sortvlrecfile(fns[0], "", reclen, keyoff, keylen, iomem)
		}
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if len(mfiles) > 0 {
			mergefiles(ofn, mfiles)
			return
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
		var mfns []string
		if reclen != 0 {
			klns, mfns, err = sortflrecfile(fn, dn, reclen, keyoff, keylen, iomem)
		} else {
			klns, mfns, err = sortvlrecfile(fn, dn, reclen, keyoff, keylen, iomem)
		}
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if len(mfns) > 0 {
			mfiles = append(mfiles, mfns...)
			continue
		}

		mfn := fmt.Sprintf("%s", filepath.Base(fn))
		mpath := filepath.Join(dn, mfn)
		savemergefile(klns, mpath)
		mfiles = append(mfiles, mpath)
	}
	mergefiles(ofn, mfiles)
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
	flag.StringVar(&ofn, "ofn", "", "output file name")
	flag.StringVar(&iomem, "iomem", "500mb", "max read memory size in kb, mb or gb")
	flag.IntVar(&reclen, "reclen", 0, "length of the fixed length record")
	flag.IntVar(&keyoff, "keyoff", 0, "offset of the key")
	flag.IntVar(&keylen, "keylen", 0, "length of the key if not whole line")
	flag.Parse()
	fns = flag.Args()

	var iom int64
	if iomem != "" {
		iom = parseiomem(iomem)
	}
	sortfiles(fns, ofn, reclen, keyoff, keylen, iom)

}
