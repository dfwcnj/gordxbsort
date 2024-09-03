package gordxbinsort

import (
	"github.com/dfwcnj/randomdata"
	"log"
	"os"
	"sort"
	"testing"
)

func Test_klrsortsl(t *testing.T) {

	log.Print("klrsortsl test")
	var l int = 32
	var r bool = true
	var e bool = false
	//ls := []int64{1, 2, 1 << 4, 1 << 8, 1 << 16, 1 << 20, 1 << 24}
	ls := []int64{1 << 4, 1 << 10, 1 << 20, 1 << 24}

	for _, nl := range ls {

		log.Print("klrsortsl test ", nl)
		var klns kvallines

		//log.Print("testing sort of ", nl)
		rsl := randomdata.Randomstrings(nl, l, r, e)
		if len(rsl) != int(nl) {
			log.Fatal("rsl: wanted len ", nl, " got ", len(rsl))
		}
		for _, s := range rsl {
			var kln kvalline
			bln := []byte(s)
			kln.line = bln
			kln.key = kln.line
			klns = append(klns, kln)
		}
		if len(klns) != int(nl) {
			log.Fatal("klns: before sort wanted len ", nl, " got ", len(klns))
		}
		slns := KLrsortsl(klns, 0)
		if len(slns) != int(nl) {
			log.Fatal("slns: after sort wanted len ", nl, " got ", len(slns))
		}

		var ssl []string
		for _, s := range slns {
			ssl = append(ssl, string(s.line))
		}
		if len(ssl) != 1 && ssl[0] == ssl[len(ssl)-1] {
			log.Fatal("strings are all equal")
		}
		if len(ssl) != int(nl) {
			log.Fatal("klrsortsl test ssl: wanted len ", nl, " got ", len(ssl))
		}
		if !sort.StringsAreSorted(ssl) {
			fp, err := os.OpenFile("/tmp/klrsortsltest", os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
			for _, l := range ssl {
				l = l + "\n"
				fp.Write([]byte(l))
			}
			fp.Close()
			log.Fatal("klrrsort2a test not in sort order")
		} else {
			log.Print("klrsortsl test passed")
		}
	}
}
