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
		//log.Printf("flreadn %s  seeking to %d\n", fp.Name(), offset)
		o, err := fp.Seek(offset, 0)
		if err != nil {
			log.Fatal(err)
		}
		if o != offset {
			log.Fatal("flreadn seek wanted", offset, " got ", o)
		}
	}
	for {

		if memused >= iomem {
			offset, err = fp.Seek(0, 1)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			return klns, offset, err
		}

		buf := make([]byte, reclen)
		if bl, err = io.ReadFull(fp, buf); err != nil {
			if err == io.EOF && bl == 0 {
				return klns, offset, err
			}
		}

		var kln kvalline
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

		memused += int64(reclen)

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
		return vlreadall(fp, offset, keyoff, keylen, finf.Size())
	}

	if offset != 0 {
		if fp.Name() == "/dev/stdin" {
			log.Fatal("vlreadn(stdin) offset ", offset)
		}
		//log.Printf("vlreadn %s  seeking to %d\n", fp.Name(), offset)
		fp.Seek(offset, 0)
	}

	r := io.Reader(fp)
	nw := bufio.NewReader(r)

	for {
		if memused >= iomem {
			//log.Println("vlreadn memused >= iomem")
			return klns, offset, err
		}

		l, err := nw.ReadString('\n')
		// Seek seens to giv the buffer seek
		offset += int64(len(l))
		if err != nil {
			if err == io.EOF && len(l) == 0 {
				//log.Println("vlreadn readstring EOF ", offset)
				return klns, offset, err
			}
			log.Fatal(err)
		}

		var kln kvalline

		bln := []byte(l)
		bls := klnullsplit(bln)
		if len(bls) == 2 {
			kln.key = bls[0]
			if bls[0][len(bls[0])-1] == '\n' {
				kln.key = bls[0][:len(bls)-1]
			}
			kln.line = bls[1]
		} else {
			kln.key = bln
			if bln[len(bln)-1] == '\n' {
				kln.key = bln[:len(bln)-1]
			}
			kln.line = bln
		}
		if keyoff != 0 {
			kln.key = kln.line[keyoff:]
			if keylen != 0 {
				kln.key = kln.line[keyoff : keyoff+keylen]
			}
		}

		klns = append(klns, kln)

		memused += int64(len(l))
	}
}
