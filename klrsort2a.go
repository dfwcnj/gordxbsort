package main

import (
	"bytes"
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
		for j := i; j > 0 && bytes.Compare(klns[j-1].key, klns[j].key) < 0; j-- {
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

	if nl == 0 {
		log.Fatal("klrsort2a: 0 len lines: ", recix)
	}
	if nl < THRESHOLD {
		return binsertionsort(klns)
	}

	for i, _ := range klns {

		if recix >= len(klns[i].key) {
			piles[0] = append(piles[0], klns[i])
			if nc == 0 {
				nc = 1
			}
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
	if nc == 0 {
		return piles[0]
	}
	if nc == 1 {
		return binsertionsort(piles[li])
	}

	for i, _ := range piles {
		if len(piles[i]) == 0 {
			continue
		}
		// sort pile
		if len(piles[i]) < THRESHOLD {
			piles[i] = binsertionsort(piles[i])
		} else {
			piles[i] = klrsort2a(piles[i], recix+1)
		}
		nc--
		if nc == 0 {
			break
		}
	}

	var slns kvallines
	for i, _ := range piles {
		for j, _ := range piles[i] {
			slns = append(slns, piles[i][j])
		}
	}
	if len(slns) != nl {
		log.Fatal("slns: ", len(slns), " nl ", nl)
	}
	return slns
}
