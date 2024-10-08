package gordxbinsort

import (
	"fmt"
	"github.com/dfwcnj/randomdata"
	"log"
	"os"
	"path"
	"testing"
)

func Test_flreadn(t *testing.T) {
	var l int = 32
	var r bool = false
	var e bool = false
	var lrs int64 = 1 << 20
	var iomem int64 = 1 << 30

	var klns kvallines
	var tklns kvallines
	var offset int64
	var err error
	var nr int

	log.Println("flreadn test")

	rsl := randomdata.Randomstrings(lrs, l, r, e)
	log.Println("flreadn test rsl ", len(rsl))

	dn, err := initmergedir("/tmp", "rdxsort")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dn)

	fn := path.Join(dn, "flreadntest")
	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	for i, _ := range rsl {
		fmt.Fprint(fp, rsl[i])
		nr++
	}
	fp.Sync()

	// file length
	offset, err = fp.Seek(0, 1)
	if err != nil {
		log.Fatal(err)
	}

	// rewind file
	offset, err = fp.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	for {
		//log.Println("flreadn test flreadn ", fn, " ", l)
		klns, offset, err = Flreadn(fp, offset, int(l), 0, 0, iomem)
		if len(klns) == 0 {
			break
		}
		for _, kln := range klns {
			if len(kln.line) != int(l) {
				log.Fatal("kln.line ", kln.line, " len ", len(kln.line))
			}
			if len(kln.key) != len(kln.line) {
				log.Fatal("kln.key ", kln.line, " len ", len(kln.line))
			}
			//log.Print(string(kln.line))
		}
		tklns = append(tklns, klns...)
	}
	if len(tklns) != int(lrs) {
		log.Fatal("flreadn: expected ", lrs, " got ", len(klns))
	}
	log.Print("flreadn test passed")
}
