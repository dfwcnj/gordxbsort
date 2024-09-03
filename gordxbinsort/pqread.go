package gordxbinsort

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"io"
	"log"
	"os"
)

// kln.key serves as the priority
type ritem struct {
	kln   kvalline
	br    *bufio.Reader
	rlen  int
	index int
}

type SPQ []*ritem

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
	ritem := x.(*ritem)
	ritem.index = n
	*pq = append(*pq, ritem)
}

func (pq *SPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	ritem := old[n-1]
	ritem.index = -1 // for safety
	*pq = old[0 : n-1]
	return ritem
}

func (pq *SPQ) update(ritem *ritem, value []byte, priority []byte) {
	ritem.kln.line = value
	ritem.kln.key = priority
	heap.Fix(pq, ritem.index)
}

func nextitem(itm ritem, kg func([]byte) [][]byte) (kvalline, error) {

	var kln kvalline
	var bln []byte

	if itm.rlen == 0 {
		l, err := itm.br.ReadString('\n')
		if err != nil {
			// log.Println("nextitem readstring ", err)
			return kln, err
		}
		// log.Print("nextitem readstring ", l)
		bln = []byte(l)
	} else {
		bln = make([]byte, itm.rlen)
		_, err := io.ReadFull(itm.br, bln)
		if err != nil {
			// log.Println("nextitem readfull ", err)
			return kln, err
		}
	}

	// default key is the whole line
	kln.line = bln
	kln.key = kln.line
	// key generator
	if kg != nil {
		bls := kg(bln)
		if len(bls) != 2 {
			log.Fatal("nextitem ", string(bln), "wanted 2  got ", len(bls))
		}
		kln.key = bls[0]
		kln.line = bls[1]
	}

	return kln, nil
}

func pqreademit(ofp *os.File, reclen int, kg func([]byte) [][]byte, fns []string) {

	pq := make(SPQ, len(fns))

	for i, fn := range fns {
		var itm ritem

		fp, err := os.Open(fn)
		if err != nil {
			log.Fatal("pqreademit open ", fn, " ", err)
		}

		itm.rlen = reclen
		rdr := io.Reader(fp)
		itm.br = bufio.NewReader(rdr)

		itm.kln, err = nextitem(itm, kg)
		if err != nil {
			log.Fatal("pqreademit setup nextitem ", fn, " ", err)
		}
		itm.index = i

		pq[i] = &itm
	}

	heap.Init(&pq)

	//nw := bufio.NewWriter(ofp)
	nw := bufio.NewWriterSize(ofp, 1<<30)

	for pq.Len() > 0 {
		ritem := heap.Pop(&pq).(*ritem)
		//s := fmt.Sprintf("%s\n", string(ritem.kln.line))
		s := fmt.Sprintf("%s", string(ritem.kln.line))
		_, err := nw.WriteString(s)
		if err != nil {
			log.Fatal("pqreademit write ", err)
		}

		ritem.kln, err = nextitem(*ritem, kg)
		if err != nil {
			continue
		}

		heap.Push(&pq, ritem)
		pq.update(ritem, ritem.kln.line, ritem.kln.key)
	}
	err := nw.Flush()
	if err != nil {
		log.Fatal("pqreademit flush", err)
	}
}
