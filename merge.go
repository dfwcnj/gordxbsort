package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// kln.key serves as the priority
type item struct {
	kln   kvalline
	inch  chan kvalline
	index int
}

func initmergedir(dn string) (string, error) {
	mdn, err := makemergedir(dn)
	if err != nil {
		if os.IsExist(err) {
			os.RemoveAll(mdn)
			return makemergedir(dn)
		}
		log.Fatal(err)
	}
	return mdn, err

}

func makemergedir(dn string) (string, error) {
	if dn == "" {
		dn = "rdxsort"
	}
	mdn, err := os.MkdirTemp("", dn)
	return mdn, err
}

// save merge file
// save key and line separated by null bute
func savemergefile(klns kvallines, fn string) string {

	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	var n = byte(0)

	for _, kln := range klns {

		knl := string(kln.key) + string(n) + string(kln.line) + "\n"

		_, err := fp.Write([]byte(knl))
		if err != nil {
			log.Fatal(err)
		}
	}
	return fn
}

// bufSplit(buf, reclen)
//
// split the buffer into a slice containing reclen records
func bufSplit(buf []byte, reclen int) lines {
	buflen := len(buf)
	var lns lines
	for o := 0; o < buflen; o += reclen {
		rec := buf[o : o+reclen-1]
		lns = append(lns, rec)
	}
	return lns
}

// klnullsplit(bl)
// example function for generating a key from a byte array
// this example assumes that the line contains a key and value
// separated by a null byte
func klnullsplit(bln []byte) [][]byte {
	var bls [][]byte
	var sep = make([]byte, 1)
	// split on null byte
	bls = bytes.Split(bln, sep)
	return bls
}

// klchan(fn, kg, out)
// klchan reads lines from file fn, creates a kvalline structure,
// populates the structure with the output of kg
func klchan(fn string, kg func([]byte) [][]byte, out chan kvalline) {
	fp, e := os.Open(fn)
	if e != nil {
		log.Fatal(e)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		var kln kvalline

		l := scanner.Text()
		if len(l) == 0 {
			continue
		}
		bln := []byte(l)

		// default key is the whole line
		kln.line = bln
		kln.key = kln.line
		// only if there is a key generator
		if kg != nil {
			bls := kg(bln)
			if len(bls) != 2 {
				log.Fatal("klchan ", fn, " ", l, " ", len(bls))
			}
			kln.key = bls[0]
			kln.line = bls[1]
		}
		out <- kln
	}
	close(out)
}

func iteminsertionsort(items []item) []item {
	n := len(items)
	if n == 1 {
		return items
	}
	for i := 0; i < n; i++ {
		for j := i; j > 0 && string(items[j-1].kln.key) > string(items[j].kln.key); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
	return items
}

func insemit(ofp *os.File, dn string, finfs []fs.DirEntry) {
	var items = make([]item, 0)

	// populate the priority queue
	for _, finf := range finfs {

		fn := filepath.Join(dn, finf.Name())
		var itm item

		inch := make(chan kvalline)
		go klchan(fn, klnullsplit, inch)

		itm.kln = <-inch
		itm.inch = inch
		items = append(items, itm)
	}

	for len(items) > 0 {
		items = iteminsertionsort(items)

		fmt.Fprintf(ofp, "%s\n", string(items[0].kln.line))

		kln, ok := <-items[0].inch
		if !ok {
			items = items[1:]
			continue
		}
		items[0].kln.key = kln.key
		items[0].kln.line = kln.line
	}
}

func mergefiles(ofn string, dn string, lpo int) {
	log.Print("multi step merge not implemented")

	finfs, err := os.ReadDir(dn)
	if err != nil {
		log.Fatal("ReadDir ", dn, ": ", err)
	}
	//log.Println("mergefiles dn ", dn)

	ofp := os.Stdout
	if ofn != "" {
		ofp, err = os.OpenFile(ofn, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer ofp.Close()
	}

	insemit(ofp, dn, finfs)
}
