package gordxbinsort

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Sortfiles(fns []string, ofn string, dn string, reclen int, keyoff int, keylen int, iomem int64) {

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
			Mergefiles(ofn, mrlen, mfiles)
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
		Mergefiles(ofn, mrlen, mfiles)
	} else {
		//log.Println("sortfiles merging", ofn, " ", reclen)
		Mergefiles(ofn, 0, mfiles)
	}
}
