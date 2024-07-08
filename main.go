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

func sortflrecfile(fn string, dn string, reclen int, keyoff int, keylen int, iomem int64) (kvallines, []string, int, error) {
	var klns kvallines
	var offset int64
	var mrlen int
	var err error
	var dlim string
	dlim = ""
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
		dn, err = initmergedir("", "rdxsort")
		if err != nil {
			log.Fatal(err)
		}
	}

	for {

		klns, offset, err = flreadn(fp, offset, reclen, keyoff, keylen, iomem)

		if err == io.EOF && len(mfiles) == 0 {
			return klns, mfiles, mrlen, err
		}
		if len(klns) == 0 {
			return klns, mfiles, mrlen, err
		}

		sklns := klrsort2a(klns, 0)

		if offset > 0 && len(sklns) > 0 {
			mfn := filepath.Join(dn, filepath.Base(fmt.Sprintf("%s%d", fn, i)))
			fn, mrlen = savemergefile(sklns, mfn, dlim)
			if fn == "" {
				log.Fatal("savemergefile failed: ", fn, " ", dn)
			}
			mfiles = append(mfiles, mfn)
		}
		if err == io.EOF {
			return sklns, mfiles, mrlen, err
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
	var dlim string
	dlim = ""
	var mfiles []string

	fp := os.Stdin
	if fn != "" {
		fp, err = os.Open(fn)
		if err != nil {
			log.Fatal("sortvlrecfile ", err)
		}
	}
	if dn == "" {
		dn, err = initmergedir("", "rdxsort")
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("sortvlrecfile dn ", dn)
	}

	for {
		klns, offset, err = vlreadn(fp, offset, keyoff, keylen, iomem)

		if err == io.EOF && len(mfiles) == 0 {
			return klns, mfiles, err
		}
		//log.Println("sortvlrecfile vlreadn klns ", len(klns))
		if len(klns) == 0 {
			return klns, mfiles, err
		}

		sklns := klrsort2a(klns, 0)
		//log.Println("sortvlrecfile sklns ", len(sklns))

		if offset > 0 && len(sklns) > 0 {
			mfn := filepath.Join(dn, filepath.Base(fmt.Sprintf("%s%d", fn, i)))
			f, _ := savemergefile(sklns, mfn, dlim)
			if f == "" {
				log.Fatal("savemergefile failed: ", mfn, " ", dn)
			}
			mfiles = append(mfiles, mfn)
			//log.Println("sortvlrecfile savemergefile ", mfn)
		}
		if err == io.EOF {
			return klns, mfiles, err
		}
		i++

	}
}

func sortfiles(fns []string, ofn string, dn string, reclen int, keyoff int, keylen int, iomem int64) {

	var klns kvallines
	var err error
	var mfiles []string
	var mrlen int = reclen
	var dlim string = ""
	if reclen == 0 {
		dlim = "\n"
	}
	//log.Printf("sortfiles ofn %s\n", ofn)
	if len(dn) == 0 {
		dn, err = initmergedir("", "rdxsort")
		if err != nil {
			log.Fatal(err)
		}
	}

	fp := os.Stdout
	if ofn != "" {
		fp, err := os.OpenFile(ofn, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()
	}

	if len(fns) == 0 {
		log.Println("sortfiles stdin ", reclen)
		if reclen != 0 {
			klns, mfiles, mrlen, err = sortflrecfile("", "", reclen, keyoff, keylen, iomem)
		} else {
			klns, mfiles, err = sortvlrecfile("", "", reclen, keyoff, keylen, iomem)
		}
		if err != nil && err != io.EOF {
			log.Fatal("sortfiles after sort ", err)
		}
		if len(mfiles) > 0 {
			mergefiles(ofn, mrlen, mfiles)
			return
		}

		for _, kln := range klns {

			_, err := fp.Write(kln.line)
			if err != nil {
				log.Fatal("sortfiles writing ", err)
			}
		}

		return
	}

	if len(fns) == 1 {
		//log.Printf("sortfiles fn %s\n", fns[0])
		if reclen != 0 {
			klns, mfiles, mrlen, err = sortflrecfile(fns[0], "", reclen, keyoff, keylen, iomem)
		} else {
			klns, mfiles, err = sortvlrecfile(fns[0], "", reclen, keyoff, keylen, iomem)
		}
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if len(mfiles) > 0 {
			mergefiles(ofn, mrlen, mfiles)
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

	for _, fn := range fns {
		var klns kvallines
		var mfns []string

		//log.Println("sortfiles sort ", fn, "", reclen)
		if reclen != 0 {
			klns, mfns, mrlen, err = sortflrecfile(fn, dn, reclen, keyoff, keylen, iomem)
		} else {
			klns, mfns, err = sortvlrecfile(fn, dn, reclen, keyoff, keylen, iomem)
		}
		if err != nil && err != io.EOF {
			log.Fatal("sortfiles after sort ", err)
		}
		if len(mfns) > 0 {
			mfiles = append(mfiles, mfns...)
			continue
		}

		mfn := fmt.Sprintf("%s", filepath.Base(fn))
		mpath := filepath.Join(dn, mfn)
		//log.Println("sortfiles saving merge file ", mpath)
		var mf string
		mf, mrlen = savemergefile(klns, mpath, dlim)
		if mf == "" {
			log.Fatal("sortfiles savemergefile failes ", mpath)
		}
		mfiles = append(mfiles, mpath)
	}
	if reclen > 0 {
		//log.Println("sortfiles merging", ofn, " ", mrlen)
		mergefiles(ofn, mrlen, mfiles)
	} else {
		//log.Println("sortfiles merging", ofn, " ", reclen)
		mergefiles(ofn, 0, mfiles)
	}
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
	var ofn, iomem, md string
	var reclen, keylen, keyoff int
	flag.StringVar(&ofn, "ofn", "", "output file name")
	flag.StringVar(&iomem, "iomem", "500mb", "max read memory size in kb, mb or gb")
	flag.StringVar(&md, "md", "", "merge sirectory")
	flag.IntVar(&reclen, "reclen", 0, "length of the fixed length record")
	flag.IntVar(&keyoff, "keyoff", 0, "offset of the key")
	flag.IntVar(&keylen, "keylen", 0, "length of the key if not whole line")
	flag.Parse()
	fns = flag.Args()

	var iom int64
	if iomem != "" {
		iom = parseiomem(iomem)
	}
	sortfiles(fns, ofn, md, reclen, keyoff, keylen, iom)

}
