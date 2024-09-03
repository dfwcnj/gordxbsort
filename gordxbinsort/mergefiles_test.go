package gordxbinsort

import (
	"bufio"
	"fmt"
	"github.com/dfwcnj/randomdata"
	"log"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func Test_mergefiles(t *testing.T) {
	var l int = 32
	var r bool = true
	var e bool = false
	var nrs int64 = 1 << 20
	var dlim string
	dlim = "\n"
	var nmf = 10
	var fns []string

	log.Print("mergefiles test")

	dn, err := initmergedir("/tmp", "rdxsort")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dn)

	for i := range nmf {
		var klns kvallines
		var kln kvalline

		rsl := randomdata.Randomstrings(nrs, l, r, e)
		for _, s := range rsl {
			bln := []byte(s)
			kln.line = bln
			kln.key = kln.line
			klns = append(klns, kln)
		}
		if len(klns) != int(nrs) {
			log.Fatal("klns: before sort wanted len ", l, " got ", len(klns))
		}

		slns := KLrsortsl(klns, 0)
		var fn = filepath.Join(dn, fmt.Sprint("file", i))
		savemergefile(slns, fn, dlim)
		fns = append(fns, fn)
	}

	mfn := "mergeout.txt"
	mpath := filepath.Join(dn, mfn)
	Mergefiles(mpath, 0, fns)

	mfp, err := os.Open(mpath)
	if err != nil {
		log.Fatal(err)
	}
	defer mfp.Close()

	scanner := bufio.NewScanner(mfp)
	var mlns []string
	for scanner.Scan() {
		l := scanner.Text()
		mlns = append(mlns, l)
	}
	if len(mlns) != int(nrs)*nmf {
		log.Fatal(mpath, " wanted ", int(nrs)*nmf, " got ", len(mlns))
	}
	if !slices.IsSorted(mlns) {
		log.Fatal("lines in ", mfn, " not in sort order")
	}
	log.Print("mergefiles test passed")

}
