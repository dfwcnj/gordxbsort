package gordxbinsort

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"log"
	"os"
)

// kln.key serves as the priority
type chitem struct {
	kln   kvalline
	inch  chan kvalline
	index int
}

type PriorityQueue []*chitem

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
	chitem := x.(*chitem)
	chitem.index = n
	*pq = append(*pq, chitem)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	chitem := old[n-1]
	chitem.index = -1 // for safety
	*pq = old[0 : n-1]
	return chitem
}

func (pq *PriorityQueue) update(chitem *chitem, value []byte, priority []byte) {
	chitem.kln.line = value
	chitem.kln.key = priority
	heap.Fix(pq, chitem.index)
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

func pqchanemit(ofp *os.File, fns []string) {
	pq := make(PriorityQueue, len(fns))

	for i, fn := range fns {
		var itm chitem

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
		chitem := heap.Pop(&pq).(*chitem)
		s := fmt.Sprintf("%s\n", string(chitem.kln.line))
		_, err := nw.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}

		kln, ok := <-chitem.inch
		if !ok {
			continue
		}
		chitem.kln = kln
		heap.Push(&pq, chitem)
		pq.update(chitem, chitem.kln.line, chitem.kln.key)
	}
	err := nw.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
