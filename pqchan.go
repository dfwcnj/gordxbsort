package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// kln.key serves as the priority
type item struct {
	kln   kvalline
	inch  chan kvalline
	index int
}

type PriorityQueue []*item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return bytes.Compare(pq[i].kln.key, pq[j].kln.key) < 0
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *item, value []byte, priority []byte) {
	item.kln.line = value
	item.kln.key = priority
	heap.Fix(pq, item.index)
}

// klchan(fn, kg, out)
// klchan reads lines from file fn, creates a kvalline structure,
// populates the structure with the output of kg
func klchan(fn string, kg func([]byte) [][]byte, out chan kvalline) {
	fp, e := os.Open(fn)
	if e != nil {
		log.Fatal(e)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		var kln kvalline

		l := scanner.Text()
		if len(l) == 0 {
			continue
		}
		bln := []byte(l)

		// default key is the whole line
		kln.line = bln
		kln.key = kln.line
		// only if there is a key generator
		if kg != nil {
			bls := kg(bln)
			if len(bls) != 2 {
				log.Fatal("klchan ", fn, " ", l, " ", len(bls))
			}
			kln.key = bls[0]
			kln.line = bls[1]
		}
		out <- kln
	}
	close(out)
}

func pqchanemit(ofp *os.File, dn string, finfs []fs.DirEntry) {
	pq := make(PriorityQueue, len(finfs))

	for i, finf := range finfs {
		fn := filepath.Join(dn, finf.Name())
		var itm item

		inch := make(chan kvalline)
		go klchan(fn, klnullsplit, inch)

		itm.kln = <-inch
		itm.inch = inch
		itm.index = i
		pq[i] = &itm
	}

	heap.Init(&pq)

	nw := bufio.NewWriter(ofp)

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*item)
		s := fmt.Sprintf("%s\n", string(item.kln.line))
		_, err := nw.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}

		kln, ok := <-item.inch
		if !ok {
			continue
		}
		item.kln = kln
		heap.Push(&pq, item)
		pq.update(item, item.kln.line, item.kln.key)
	}
	err := nw.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
