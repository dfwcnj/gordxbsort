package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
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
	return string(pq[i].kln.key) < string(pq[j].kln.key)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// lowest priority item
func (pq *PriorityQueue) Bottom() any {
	old := *pq
	item := old[0]
	return item
}

// highest priority item
func (pq *PriorityQueue) Top() any {
	old := *pq
	n := len(*pq)
	item := old[n-1]
	return item
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an item in the queue.
func (pq *PriorityQueue) update(item *item, value string, priority string) {
	item.kln.line = []byte(value)
	item.kln.key = []byte(priority)
	heap.Fix(pq, item.index)
}

func initmergedir(dn string) (string, error) {
	mdn, err := makemergedir(dn)
	if err != nil {
		if os.IsExist(err) {
			os.RemoveAll(mdn)
			return makemergedir(dn)
		}
		log.Fatal(err)
	}
	return mdn, err

}

func makemergedir(dn string) (string, error) {
	if dn == "" {
		dn = "rdxsort"
	}
	mdn, err := os.MkdirTemp("", dn)
	return mdn, err
}

func mergefiles(ofn string, dn string, lpo int) {
	log.Print("multi step merge not implemented")
	ofp := os.Stdout
	if ofn != "" {
		var mpath = filepath.Join(dn, ofn)
		log.Print("mpath ", mpath)
		log.Print("opening ", mpath)
		ofp, err := os.OpenFile(mpath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer ofp.Close()
	}
	finfs, err := os.ReadDir(dn)
	if err != nil {
		log.Fatal("ReadDir ", dn, ": ", err)
	}
	pq := make(PriorityQueue, len(finfs))
	i := 0

	// populate the priority queue
	for _, finf := range finfs {
		fn := finf.Name()
		inch := make(chan kvalline)
		go klchan(fn, klnullsplit, inch)
		var nit item
		nit.kln = <-inch
		nit.inch = inch
		nit.index = i
		pq[i] = &nit
		i++
	}

	for pq.Len() > 0 {
		item := pq.Top().(*item)
		fmt.Fprint(ofp, item.kln.line)
		kln, ok := <-item.inch
		if !ok {
			_ = heap.Pop(&pq)
		}
		pq.update(item, string(kln.line), string(kln.key))
	}
}

// save merge file
// save key and line separated by null bute
func savemergefile(klns kvallines, fn string) string {

	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	var n = byte(0)

	for _, kln := range klns {

		knl := string(kln.key) + string(n) + string(kln.line) + "\n"

		_, err := fp.Write([]byte(knl))
		if err != nil {
			log.Fatal(err)
		}
	}
	return fn
}

// bufSplit(buf, reclen)
//
// split the buffer into a slice containing reclen records
func bufSplit(buf []byte, reclen int) lines {
	buflen := len(buf)
	var lns lines
	for o := 0; o < buflen; o += reclen {
		rec := buf[o : o+reclen-1]
		lns = append(lns, rec)
	}
	return lns
}

// klnullsplit(bl)
// example function for generating a key from a byte array
// this example assumes that the line contains a key and value
// separated by a null byte
func klnullsplit(bln []byte) [][]byte {
	var bls [][]byte
	var sep = make([]byte, 1)
	// split on null byte
	bls = bytes.Split(bln, sep)
	// there can be only one
	if len(bls) != 2 {
		log.Println("klnullsplit wanted ", 2, " got ", len(bls), " parts")
	}
	return bls
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
