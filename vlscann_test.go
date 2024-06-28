package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
)

//type kvalline struct {
//	key  []byte
//	line []byte
//}

func Test_vlscann(t *testing.T) {
	var l uint = 32
	var lpo uint = 1 << 16
	var nr int

	var klns kvallines

	rsl := randomstrings(lpo, l)
	td := os.TempDir()
	fn := path.Join(td, "rdxsort", "rtxt.txt")

	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	log.Println(fn, " created")

	for _, l := range rsl {
		fmt.Fprintln(fp, l+"\n")
		nr++
	}
	log.Println(fn, " filled ", nr)

	_, err = fp.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fn, " file pointer  rewound")
	klns, _, err = vlscann(fp, 0, 0, 0, 0)
	for _, kln := range klns {
		if len(kln.line) == 0 {
			log.Fatal("vlscann len(kln.line) == 0")
		}
		if len(kln.key) != len(kln.line) {
			log.Fatal("vlscann len(kln.key) != len(kln.line)")
		}
		//log.Print(string(kln.line))
	}
	if len(klns) != int(lpo) {
		log.Fatal("vlscann: expected ", lpo, " got ", len(klns))
	}
	log.Print("vlscann test passed")
}
