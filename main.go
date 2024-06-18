package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

type line []byte
type lines []line

//type kline struct {
//	key  []byte
//	line []byte
//}
//type klines []kline

func main() {
	var fns []string
	var keylen, keyoff int
	var lps, lpm int
	flag.IntVar(&keylen, "keylen", 0, "length of the key if not whole line")
	flag.IntVar(&keyoff, "keyoff", 0, "offset of the key")
	flag.IntVar(&lps, "lps", 1<<20, "lines per sort ")
	flag.IntVar(&lpm, "lpm", 1<<20, "lines per merge ")
	flag.Parse()
	fns = flag.Args()

	var klns klines

	var err error

	for _, fn := range fns {
		fp := os.Stdin
		if fn != "" {
			fp, err = os.Open(fn)
			if err != nil {
				log.Fatal(err)
			}
			defer fp.Close()
		}

		scanner := bufio.NewScanner(fp)
		// option, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			l := scanner.Text()

			bln := []byte(l)
			var kl kline
			kl.line = bln
			kl.key = kl.line
			klns = append(klns, kl)
		}
		klrsort2a(klns, 0)
		for _, l := range klns {
			fmt.Print(string(l.line))
		}
	}

}
