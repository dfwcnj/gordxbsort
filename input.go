package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

//type kvalline struct {
//        key  []byte
//        line []byte
//}
//type kvallines []kvalline

func flreadall(fp *os.File, offset int64, reclen int, keyoff int, keylen int, iomem int64) (kvallines, int64, error) {

	var klns kvallines

	buf, err := io.ReadAll(fp)
	if err != nil {
		log.Fatal(err)
	}
	var r io.Reader = bytes.NewReader(buf)

	recbuf := make([]byte, reclen)
	for {
		var kln kvalline
		_, err := io.ReadFull(r, recbuf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			return klns, 0, nil
		}
		kln.line = recbuf
		kln.key = kln.line
		klns = append(klns, kln)
	}
}

func flreadn(fp *os.File, offset int64, reclen int, keyoff int, keylen int, iomem int64) (kvallines, int64, error) {

	var klns kvallines
	var nr int // number records read
	var bl int
	var err error
	var memused int64

	finf, err := fp.Stat()
	if err != nil {
		log.Fatal()
	}
	if finf.Size() < iomem {
		return flreadall(fp, offset, reclen, keyoff, keylen, finf.Size())
	}

	if keyoff+keylen > reclen {
		log.Fatal("key dimension extends beyond reclen")
	}

	if offset != 0 {
		if fp.Name() == "/dev/stdin" {
			log.Fatal("flreadn(stdin) more than iomem")
		}
		if offset == finf.Size() {
			log.Fatal("flreadn ", fp.Name(), " end of file")
		}
		log.Printf("flreadn %s  seeking to %d\n", fp.Name(), offset)
		fp.Seek(offset, 0)
	}
	for {
		buf := make([]byte, reclen)
		if bl, err = io.ReadFull(fp, buf); err != nil {
			if err == io.EOF {
				log.Println("flreadn readfull eof")
				return klns, offset, err
			}
			log.Fatal("flreadn: ", err)
		}

		memused += int64(reclen)
		if memused >= iomem || bl == 0 {
			offset, err = fp.Seek(0, 1)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			log.Println("flreadn memused ", memused, " iomem ", iomem)
			return klns, offset, err
		}

		var kln kvalline
		// to avoid having to make buf in the loop
		// mistake??
		bls := klnullsplit(buf)
		if len(bls) == 2 {
			kln.key = bls[0]
			kln.line = bls[1]
		} else {
			kln.line = buf
			kln.key = kln.line
		}
		if keyoff != 0 {
			kln.key = kln.line[keyoff:]
			if keylen != 0 {
				kln.key = kln.line[keyoff : keyoff+keylen]
			}
		}
		klns = append(klns, kln)
		nr++
	}
}

func vlreadall(fp *os.File, offset int64, keyoff int, keylen int, iomem int64) (kvallines, int64, error) {
	var klns kvallines
	buf, err := io.ReadAll(fp)
	if err != nil {
		return klns, offset, err
	}
	lines := strings.Split(string(buf), "\n")
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		var kln kvalline
		bln := []byte(l)
		bls := klnullsplit(bln)
		if len(bls) == 2 {
			kln.key = bls[0]
			kln.line = bls[1]
		} else {
			kln.line = bln
			kln.key = bln
		}
		if keyoff != 0 {
			kln.key = kln.line[keyoff:]
			if keylen != 0 {
				kln.key = kln.line[keyoff : keyoff+keylen]
			}
		}
		klns = append(klns, kln)
	}
	return klns, offset, nil
}

func vlreadn(fp *os.File, offset int64, keyoff int, keylen int, iomem int64) (kvallines, int64, error) {

	var klns kvallines
	var memused int64

	finf, err := fp.Stat()
	if err != nil {
		log.Fatal()
	}
	if finf.Size() < iomem {
		log.Println("vlreadn vlreadall")
		return vlreadall(fp, offset, keyoff, keylen, finf.Size())
	}

	if offset != 0 {
		if fp.Name() == "/dev/stdin" {
			log.Fatal("vlreadn(stdin) offset ", offset)
		}
		if offset == finf.Size() {
			log.Fatal("vlreadn ", fp.Name(), " end of file")
		}
		log.Printf("vlreadn %s  seeking to %d\n", fp.Name(), offset)
		fp.Seek(offset, 0)
	}

	r := io.Reader(fp)
	nw := bufio.NewReader(r)

	for {
		l, err := nw.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return klns, 0, err
			}
		}
		memused += int64(len(l))
		if memused >= iomem {
			offset, err = fp.Seek(0, 1)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			log.Println("vlreadn memused ", memused, " iomem ", iomem)
			log.Println("vlreadn ", len(klns), " ", offset, " ", err)
			return klns, offset, err
		}
		bln := []byte(l)
		bls := klnullsplit(bln)
		var kln kvalline
		if len(bls) == 2 {
			kln.key = bls[0]
			kln.line = bls[1]
		} else {
			kln.line = bln
			kln.key = bln
		}
		if keyoff != 0 {
			kln.key = kln.line[keyoff:]
			if keylen != 0 {
				kln.key = kln.line[keyoff : keyoff+keylen]
			}
		}
		klns = append(klns, kln)
	}
}
