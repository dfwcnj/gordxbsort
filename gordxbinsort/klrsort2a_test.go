package gordxbinsort

import (
	"github.com/dfwcnj/randomdata"
	"log"
	"slices"
	"testing"
)

func Test_klrsort2a(t *testing.T) {

	ls := []int{1 << 3, 1 << 4, 1 << 5, 1 << 6}
	ns := []int64{1 << 3, 1 << 16, 1 << 20, 1 << 24}

	for _, ll := range ls {
		for _, nl := range ns {

			var klns kvallines
			var l int = ll
			var r bool = true
			var e bool = false
			log.Print("testing klrsort2a of ", nl, " random strings length ", l)
			rsl := randomdata.Randomstrings(nl, l, r, e)
			if len(rsl) != int(nl) {
				log.Fatal("klrsort2a test rsl: wanted len ", nl, " got ", len(rsl))
			}
			for _, s := range rsl {
				var kln kvalline
				bln := []byte(s)
				kln.line = bln
				kln.key = kln.line
				klns = append(klns, kln)
			}
			if len(klns) != int(nl) {
				log.Print(klns)
				log.Fatal("klrsort2a test lns: before klrsort2a wanted len ", nl, " got ", len(klns))
			}
			KLrsort2a(klns)
			if len(klns) != int(nl) {
				//log.Print(ulns)
				log.Fatal("klrsort2a test ulns: after klrsort2a wanted len ", nl, " got ", len(klns))
			}
			var ssl []string
			for _, s := range klns {
				ssl = append(ssl, string(s.key))
			}

			if !slices.IsSorted(ssl) {
				log.Fatal("klrsort2a failed for size ", nl)
			} else {
				log.Print("klrsort2a test passed for ", nl)
			}

		}
	}
}
