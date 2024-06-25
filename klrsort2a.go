package main

import (
	"log"
)

const THRESHOLD int = 1 << 5

type kvalline struct {
	key  []byte
	line []byte
}
type kvallines []kvalline

func binsertionsort(klns kvallines) kvallines {
	n := len(klns)
	if n == 1 {
		return klns
	}
	for i := 0; i < n; i++ {
		for j := i; j > 0 && string(klns[j-1].key) > string(klns[j].key); j-- {
			klns[j], klns[j-1] = klns[j-1], klns[j]
		}
	}
	return klns
}

// bostic
func klrsort2a(klns kvallines, recix int) kvallines {
	var piles = make([]kvallines, 256)
	var nc int
	var li int
	nl := len(klns)
	pilelen := make([]int, 256)

	if nl == 0 {
		log.Fatal("klrsort2a: 0 len lines: ", recix)
	}
	if nl < THRESHOLD {
		return binsertionsort(klns)
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
		return binsertionsort(piles[li])
	}

	for i, _ := range piles {
		if len(piles[i]) == 0 {
			continue
		}
		pilelen[i] = len(piles[i])
		// sort pile
		if len(piles[i]) < THRESHOLD {
			piles[i] = binsertionsort(piles[i])
			if len(piles[i]) != pilelen[i] {
				log.Fatal("pilelen[", i, "] ", pilelen[i], "len(piles[i]) ", len(piles[i]))
			}
		} else {
			piles[i] = klrsort2a(piles[i], recix+1)
			if len(piles[i]) != pilelen[i] {
				log.Fatal("pilelen[", i, "] ", pilelen[i], "len(piles[i]) ", len(piles[i]))
			}
		}
		nc--
		if nc == 0 {
			break
		}
	}

	var slns kvallines
	for i, _ := range piles {
		if len(piles[i]) != pilelen[i] {
			log.Fatal("pilelen[", i, "] ", pilelen[i], "len(piles[i]) ", len(piles[i]))
		}
		for j, _ := range piles[i] {
			slns = append(slns, piles[i][j])
		}
	}
	if len(slns) != nl {
		log.Fatal("slns: ", len(slns), " nl ", nl)
	}
	return slns
}
