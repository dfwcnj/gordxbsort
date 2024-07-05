package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

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

	nw := bufio.NewWriter(ofp)
	for len(items) > 0 {
		items = iteminsertionsort(items)

		s := fmt.Sprintf("%s\n", string(items[0].kln.line))
		_, err := nw.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Fprintf(ofp, "%s\n", string(items[0].kln.line))

		kln, ok := <-items[0].inch
		if !ok {
			items = items[1:]
			continue
		}
		items[0].kln.key = kln.key
		items[0].kln.line = kln.line
	}
	err := nw.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
