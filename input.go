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
		log.Printf("sfpread seeking to %d\n", offset)
		_, err := fp.Seek(offset, 0)
		if err != nil {
			if err == io.EOF {
				return klns, offset, err
			}
			log.Fatal("flreadn: ", err)
		}
	}
	for {
		buf := make([]byte, reclen)
		if bl, err = io.ReadFull(fp, buf); err != nil {
			if err == io.EOF {
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
			return klns, memused, err
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

func vlscann(fp *os.File, offset int64, keyoff int, keylen int, iomem int64) (kvallines, int64, error) {

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
			log.Fatal("vlscann(stdin) offset ", offset)
		}
		_, err := fp.Seek(offset, 0)
		if err != nil {
			if err == io.EOF {
				return klns, offset, err
			}
			log.Fatal("vlscann", err)
		}
	}

	n := 1 << 30
	sbuf := make([]byte, n)

	scanner := bufio.NewScanner(fp)
	scanner.Buffer(sbuf, n*2)

	for scanner.Scan() {
		var kln kvalline
		l := scanner.Text()
		if len(l) == 0 {
			continue
		}
		memused += int64(len(l))
		if memused >= iomem {
			offset, err = fp.Seek(0, 1)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			return klns, memused, err
		}

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
	if scanner.Err() != nil {
		log.Fatal("vlscann: ", scanner.Err())
	}
	return klns, offset, nil

}
