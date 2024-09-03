package gordxbinsort

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

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

		klns, offset, err = Flreadn(fp, offset, reclen, keyoff, keylen, iomem)

		if err == io.EOF && len(mfiles) == 0 {
			return klns, mfiles, mrlen, err
		}
		if len(klns) == 0 {
			return klns, mfiles, mrlen, err
		}

		sklns := KLrsortsl(klns, 0)

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
