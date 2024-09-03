package gordxbinsort

import (
	"fmt"
	"github.com/dfwcnj/randomdata"
	"log"
	"os"
	"path"
	"testing"
)

func Test_vlreadn(t *testing.T) {
	var l int = 32
	var r bool = true
	var e bool = false
	var nrs int64 = 1 << 20
	var iomem int64 = 1 << 30
	var nr int

	var klns kvallines

	log.Print("vlreadn test")

	dn, err := initmergedir("/tmp", "rdxsort")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dn)

	rsl := randomdata.Randomstrings(nrs, l, r, e)

	fn := path.Join(dn, "vlreadn")
	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	for _, l := range rsl {
		fmt.Fprintln(fp, l+"\n")
		nr++
	}

	_, err = fp.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	klns, _, err = Vlreadn(fp, 0, 0, 0, iomem)
	for _, kln := range klns {
		if len(kln.line) == 0 {
			log.Fatal("vlreadn len(kln.line) == 0")
		}
		if len(kln.key) != len(kln.line) {
			log.Fatal("vlreadn len(kln.key) != len(kln.line)")
		}
		//log.Print(string(kln.line))
	}
	if len(klns) != int(nrs) {
		log.Fatal("vlreadn: expected ", nrs, " got ", len(klns))
	}
	log.Print("vlreadn test passed")
}
