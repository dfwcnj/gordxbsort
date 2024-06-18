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

type kline struct {
	key  []byte
	line []byte
}

type klines []kline

func main() {
	var fn string
	flag.StringVar(&fn, "file", "", "name of file to sort")
	flag.Parse()
	var klns klines

	var err error

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
	drsort2a(klns, 0)
	for _, l := range klns {
		fmt.Print(string(l.line))
	}

}
