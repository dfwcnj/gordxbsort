package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// kln.key serves as the priority
type sitem struct {
	kln   kvalline
	scn   *bufio.Reader
	index int
}

type SPQ []*sitem

func (pq SPQ) Len() int { return len(pq) }

func (pq SPQ) Less(i, j int) bool {
	return bytes.Compare(pq[i].kln.key, pq[j].kln.key) < 0
}

func (pq SPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *SPQ) Push(x interface{}) {
	n := len(*pq)
	sitem := x.(*sitem)
	sitem.index = n
	*pq = append(*pq, sitem)
}

func (pq *SPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	sitem := old[n-1]
	sitem.index = -1 // for safety
	*pq = old[0 : n-1]
	return sitem
}

func (pq *SPQ) update(sitem *sitem, value []byte, priority []byte) {
	sitem.kln.line = value
	sitem.kln.key = priority
	heap.Fix(pq, sitem.index)
}

func nextitem(scn *bufio.Reader, kg func([]byte) [][]byte) (kvalline, error) {

	var kln kvalline

	l, err := scn.ReadString('\n')
	if err != nil {
		return kln, err
	}

	bln := []byte(l)
	// default key is the whole line
	kln.line = bln
	kln.key = kln.line
	// only if there is a key generator
	if kg != nil {
		bls := kg(bln)
		if len(bls) != 2 {
			log.Fatal("nextitem len ", bls, "wanted 2  got ", len(bls))
		}
		kln.key = bls[0]
		kln.line = bls[1]
	}

	return kln, nil
}

func pqreademit(ofp *os.File, dn string, kg func([]byte) [][]byte, finfs []fs.DirEntry) {
	pq := make(SPQ, len(finfs))

	for i, finf := range finfs {
		fn := filepath.Join(dn, finf.Name())
		var itm sitem

		fp, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}

		//itm.scn = bufio.NewScanner(fp)
		r := io.Reader(fp)
		itm.scn = bufio.NewReader(r)
		itm.kln, err = nextitem(itm.scn, kg)
		if err != nil {
			log.Fatal(err)
		}
		itm.index = i

		pq[i] = &itm
	}

	heap.Init(&pq)

	nw := bufio.NewWriter(ofp)

	for pq.Len() > 0 {
		sitem := heap.Pop(&pq).(*sitem)
		//s := fmt.Sprintf("%s\n", string(sitem.kln.line))
		s := fmt.Sprintf("%s", string(sitem.kln.line))
		_, err := nw.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}

		sitem.kln, err = nextitem(sitem.scn, kg)
		if err != nil {
			continue
		}

		heap.Push(&pq, sitem)
		pq.update(sitem, sitem.kln.line, sitem.kln.key)
	}
	err := nw.Flush()
	if err != nil {
		log.Fatal(err)
	}
	err = nw.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
