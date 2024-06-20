package main

import (
	"bytes"
	"bufio"
	"log"
	"os"
)

// klnulldelim(bl)
// example function for generating a key from a byte array
// this example assumes that the line contains a key and value
// separated by a null byte
func klnulldelim(bl []byte) [][]byte {
    var bls [][]byte
    var sep = make([]byte, 1)
    // split on null byte
    bls = bytes.Split(bl, sep)
    // there can be only one
    if len(bls) != 2 {
        log.Println("", string(bl), " ", len(bls), " parts")
    }
    return bls
}

// klchan(fn, kg, out)
// klchan reads lines from file fn, creates a kvalline structure,
// populates the structure with the output of kg
func klchan(fn string, kg func(bln []byte) lines, out chan kvalline) {
	fp, e := os.Open(fn)
	if e != nil {
		log.Fatal(e)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		var kln kvalline

		l := scanner.Text()
		bln := []byte(l)
		// default key is the whole line
		kln.line = bln
		kln.key = kln.line
		// only if there is a key generator
		if kg != nil {
			bls := kg(bln)
			if len(bls) != 2 {
				log.Fatal("klchan ", fn, " ", l, " ", len(bls) )
			}
			kln.key = bls[0]
			kln.line = bls[1]
		}
		out <- kln
	}
}
