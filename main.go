package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type line []byte
type lines []line

//	type kvalline struct {
//		key  []byte
//		line []byte
//	}
//
// type kvallines []kvalline
func mergefiles(dn string, keyoff int, keylen int, lpm int) {
}

func savemergefile(klns kvallines, dn string) {

}

func sortflrecfile(fn string, dn string, reclen int, keyoff int, keylen int, lps int, lpm int) {
	var nr int
	var klns kvallines

	fp, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, reclen)
	for {
		if _, err := io.ReadFull(fp, buf); err != nil {
			log.Fatal(err)
		}
		// radix sort used bytes
		bln := buf
		var kl kvalline
		kl.line = bln
		kl.key = kl.line
		klns = append(klns, kl)

		if lps > 0 && nr >= lps {
			if dn == "" {
				log.Fatal("sortfile, no temporary directory")
			}
			klrsort2a(klns, 0)

			// call savemergefile()
			savemergefile(klns, dn)

			// remove after testing
			// need to merge
			for _, l := range klns {
				fmt.Print(string(l.line))
			}
		}
	}
}

func sortfile(fn string, dn string, reclen int, keyoff int, keylen int, lps int, lpm int) {
	if reclen > 0 {
		sortflrecfile(fn, dn, reclen, keyoff, keylen, lps, lpm)
	}
	var nr int
	var offset int64
	var klns kvallines

	fp, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		l := scanner.Text()

		// radix sort used bytes
		bln := []byte(l)
		var kl kvalline
		kl.line = bln
		kl.key = kl.line
		klns = append(klns, kl)

		if lps > 0 && nr >= lps {
			if dn == "" {
				log.Fatal("sortfile, no temporary directory")
			}
			offset, _ = fp.Seek(0, 1)
			klrsort2a(klns, 0)

			// call savemergefile()

			// remove after testing
			// need to merge
			for _, l := range klns {
				fmt.Print(string(l.line))
			}

		}
		nr++
	}
	// have to integrate with multifile merge
	if offset > 0 {
		mergefiles(dn, keyoff, keylen, lpm)
	}
	klrsort2a(klns, 0)
	for _, l := range klns {
		fmt.Print(string(l.line))
	}

}

func sortfiles(fns []string, reclen int, keyoff int, keylen int, lps int, lpm int) {

	if len(fns) == 0 {
		sortfile("", "", reclen, keyoff, keylen, lps, lpm)
		return
	}
	if len(fns) == 1 {
		sortfile(fns[0], "", reclen, keyoff, keylen, lps, 0)
		return
	}
	dn, err := os.MkdirTemp("", "sort")
	if err != nil {
		log.Fatal(err)
	}
	for _, fn := range fns {
		sortfile(fn, dn, reclen, keyoff, keylen, lps, lpm)
	}

}

func main() {
	var fns []string
	var reclen, keylen, keyoff int
	var lps, lpm int
	flag.IntVar(&reclen, "reclen", 0, "length of the fixed length record")
	flag.IntVar(&keyoff, "keyoff", 0, "offset of the key")
	flag.IntVar(&keylen, "keylen", 0, "length of the key if not whole line")
	flag.IntVar(&lps, "lps", 1<<20, "lines per sort ")
	flag.IntVar(&lpm, "lpm", 1<<20, "lines per merge ")
	flag.Parse()
	fns = flag.Args()

	sortfiles(fns, reclen, keyoff, keylen, lps, lpm)

}
