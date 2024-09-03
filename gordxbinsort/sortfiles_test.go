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

func Test_sortfiles(t *testing.T) {
	var l int = 32
	var r bool = true
	var e bool = false
	var nrs int64 = 1 << 20
	//var iomem int64 = 1 << 29
	var iomem int64 = 1<<24 + 1<<20
	var nmf = 10
	var dlim string
	dlim = "\n"

	log.Print("sortfiles test")

	dn, err := initmergedir("/tmp", "rdxsort")

	var fns []string
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
			log.Fatal("sortfiles test before sort wanted len ", l, " got ", len(klns))
		}

		//log.Println("sorting file", i)
		slns := KLrsortsl(klns, 0)
		var fn = filepath.Join(dn, fmt.Sprint("sortfilestest", i))
		//log.Println("saving file", i)
		savemergefile(slns, fn, dlim)
		fns = append(fns, fn)
	}

	mfn := "mergeout.txt"
	mpath := filepath.Join(dn, mfn)

	//log.Println("sortfiles test sorting files to ", mpath)
	Sortfiles(fns, mpath, dn, 0, 0, 0, iomem)

	mfp, err := os.Open(mpath)
	if err != nil {
		log.Fatal("sortfiles test ", err)
	}
	defer mfp.Close()

	scanner := bufio.NewScanner(mfp)
	var mlns []string
	for scanner.Scan() {
		l := scanner.Text()
		mlns = append(mlns, l)
	}
	if len(mlns) != int(nrs)*nmf {
		log.Fatal("sortfiles test n wanted ", int(nrs)*nmf, " got ", len(mlns))
	}
	if !slices.IsSorted(mlns) {
		log.Fatal("sortfiles test lines in ", mfn, " not in sort order")
	}
	log.Print("sortfiles test passed")

}
