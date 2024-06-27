package main

import (
	"bufio"
	"bytes"

	"log"
	"os"
)

type item struct {
	kln   kvalline
	in    chan kvalline
	index int
}

type PriorityQuele []*item

func mergefiles(dn string, lpo int) {
	fns, err := os.ReadDir(dn)
	if err != nil {
		log.Fatal(err)
	}

	for fn := range fns {
		log.Println(fn)
	}
	log.Fatal("still need to merge files")
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
	// there can be only one
	if len(bls) != 2 {
		log.Println("klnullsplit wanted ", 2, " got ", len(bls), " parts")
	}
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
