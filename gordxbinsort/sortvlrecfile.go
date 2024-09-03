package gordxbinsort

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

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
		klns, offset, err = Vlreadn(fp, offset, keyoff, keylen, iomem)

		if err == io.EOF && len(mfiles) == 0 {
			return klns, mfiles, err
		}
		//log.Println("sortvlrecfile vlreadn klns ", len(klns))
		if len(klns) == 0 {
			return klns, mfiles, err
		}

		sklns := KLrsortsl(klns, 0)
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
