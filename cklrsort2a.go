package main

import (
	"log"
)

//const THRESHOLD int = 1 << 5
//
//type kvalline struct {
//	key  []byte
//	line []byte
//}
//type kvallines []kvalline

func cbinsertionsort(klns kvallines, out chan kvallines) {
	n := len(klns)
	if n == 1 {
		out <- klns
	}
	for i := 0; i < n; i++ {
		for j := i; j > 0 && string(klns[j-1].key) > string(klns[j].key); j-- {
			klns[j], klns[j-1] = klns[j-1], klns[j]
		}
	}
	out <- klns
}

// bostic
func cklrsort2a(klns kvallines, recix int, out chan kvallines) {
	var piles = make([]kvallines, 256)
	var nc int
	var li int
	nl := len(klns)
	pilelen := make([]int, 256)

	if nl == 0 {
		log.Fatal("cklrsort2a: 0 len lines: ", recix)
	}
	if nl < THRESHOLD {
		cbinsertionsort(klns, out)
		return
	}

	for i, _ := range klns {

		if recix >= len(klns[i].key) {
			continue
		}

		// append kvalline to the pile indexed by c
		c := int(klns[i].key[recix])
		piles[int(c)] = append(piles[c], klns[i])
		if len(piles[c]) == 1 {
			nc++ // number of piles so far
		}
		li = c
	}
	if nc == 1 {
		cbinsertionsort(piles[li], out)
	}

	chans := make([]chan kvallines, 0)
	for i, _ := range piles {
		if len(piles[i]) == 0 {
			continue
		}
		pilelen[i] = len(piles[i])
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
