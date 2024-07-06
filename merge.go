package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

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
	nw := bufio.NewWriter(fp)

	var n = byte(0)

	for _, kln := range klns {

		knl := string(kln.key) + string(n) + string(kln.line) + "\n"

		//_, err := fp.Write([]byte(knl))
		_, err := nw.WriteString(knl)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = nw.Flush()
	if err != nil {
		log.Fatal(err)
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

func mergefiles(ofn string, fns []string) {
	log.Print("multi step merge not implemented")

	var err error

	ofp := os.Stdout
	if ofn != "" {
		ofp, err = os.OpenFile(ofn, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer ofp.Close()
	}

	pqreademit(ofp, klnullsplit, fns)
	//pqchanemit(ofp, knullsplit, fns)
	//insemit(ofp, knullsplit, fns)
}
