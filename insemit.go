package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type sitem struct {
	kln   kvalline
	inch  chan kvalline
	index int
}

var sitems []sitem

func iteminsertionsort(sitems []sitem) []sitem {
	n := len(sitems)
	if n == 1 {
		return sitems
	}
	for i := 0; i < n; i++ {
		for j := i; j > 0 && string(sitems[j-1].kln.key) > string(sitems[j].kln.key); j-- {
			sitems[j], sitems[j-1] = sitems[j-1], sitems[j]
		}
	}
	return sitems
}

func insemit(ofp *os.File, fns []string) {
	var sitems = make([]sitem, 0)

	// populate the priority queue
	for _, fn := range fns {

		var itm sitem

		inch := make(chan kvalline)
		go klchan(fn, klnullsplit, inch)

		itm.kln = <-inch
		itm.inch = inch
		sitems = append(sitems, itm)
	}

	nw := bufio.NewWriter(ofp)
	for len(sitems) > 0 {
		sitems = iteminsertionsort(sitems)

		s := fmt.Sprintf("%s\n", string(sitems[0].kln.line))
		_, err := nw.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Fprintf(ofp, "%s\n", string(sitems[0].kln.line))

		kln, ok := <-sitems[0].inch
		if !ok {
			sitems = sitems[1:]
			continue
		}
		sitems[0].kln.key = kln.key
		sitems[0].kln.line = kln.line
	}
	err := nw.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
