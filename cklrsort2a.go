package main

import (
	"bytes"
	"log"
)

func cbinsertionsort(klns kvallines, out chan kvallines) {
	n := len(klns)
	if n == 1 {
		out <- klns
	}
	for i := 0; i < n; i++ {
		for j := i; j > 0 && bytes.Compare(klns[j-1].key, klns[j].key) > 0; j-- {
			klns[j], klns[j-1] = klns[j-1], klns[j]
		}
	}
	out <- klns
}

// bostic
func cklrsort2a(klns kvallines, recix int, out chan kvallines) {
	var piles = make([]kvallines, 256)
	var nc int
	nl := len(klns)

	if nl == 0 {
		log.Fatal("cklrsort2a: 0 len lines: ", recix)
	}
	if nl < THRESHOLD {
		cbinsertionsort(klns, out)
		return
	}

	for i, _ := range klns {

		var c int
		if recix >= len(klns[i].key) {
			c = 0
		} else {
			c = int(klns[i].key[recix])
		}
		piles[int(c)] = append(piles[c], klns[i])
		if len(piles[c]) == 1 {
			nc++ // number of piles so far
		}
	}
	if len(piles[0]) > 1 {
		piles[0] = binsertionsort(piles[0])
	}
	if nc == 1 {
		cbinsertionsort(klns, out)
	}

	chans := make([]chan kvallines, 0)
	for i, _ := range piles {
		if len(piles[i]) == 0 {
			continue
		}
		// sort pile
		cini := make(chan kvallines, 0)
		chans = append(chans, cini)
		if len(piles[i]) < THRESHOLD {
			go cbinsertionsort(piles[i], cini)
		} else {
			go cklrsort2a(piles[i], recix+1, cini)
		}
		nc--
		if nc == 0 {
			break
		}
	}

	var slns kvallines
	for i, _ := range chans {
		klns := <-chans[i]
		slns = append(slns, klns...)
	}
	out <- slns
}
