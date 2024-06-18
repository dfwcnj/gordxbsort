package main

const THRESHOLD int = 1 << 5

//type line []byte
//type lines []line

//type kline struct {
//	key  []byte
//	line []byte
//}
//type klines []kline

func binsertionsort(klns klines) {
	n := len(klns)
	if n == 1 {
		return
	}
	for i := 0; i < len(klns); i++ {
		for j := i; j > 0 && string(klns[j-1].key) > string(klns[j].key); j-- {
			klns[j], klns[j-1] = klns[j-1], klns[j]
		}
	}
}

// bostic
func drsort2a(klns klines, recix int) {
	var piles = make([]klines, 256)
	var nc int

	if len(klns) == 0 {
		return
	}
	if len(klns) < THRESHOLD {
		binsertionsort(klns)
		return
	}

	for _, l := range klns {

		if recix >= len(l.key) {
			continue
		}

		// aooend kline to the pile indexed by c
		c := int(l.key[recix])
		piles[int(c)] = append(piles[int(c)], l)
		if len(piles[c]) == 1 {
			nc++ // number of piles so far
		}
	}

	for _, p := range piles {
		if len(p) == 0 {
			continue
		}
		// sort pile
		drsort2a(p, recix+1)
		nc--
		if nc == 0 {
			break
		}
	}
	clear(klns)
	for _, p := range piles {
		if len(p) == 0 {
			continue
		}
		for _, l := range p {
			klns = append(klns, l)
		}
	}
}
