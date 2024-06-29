package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
)

//type kvalline struct {
//        key  []byte
//        line []byte
//}
//type kvallines []kvalline

// flread(fp, reclen, keyoff, keylen, lpo)
// reads up to lpo fixed length records of reclen length from file ptr fp
// at offset
// returns a slice of kvalline, the current offset, and err
func flreadn(fp *os.File, offset int64, reclen int, keyoff int, keylen int, lpo int) (kvallines, int64, error) {

	var klns kvallines
	var nr int // number records read
	var bl int
	var err error

	if keyoff+keylen > reclen {
		log.Fatal("key dimension extends beyond reclen")
	}

	if offset != 0 {
		if fp.Name() == "/dev/stdin" {
			log.Fatal("flreadn(stdin) lpo less than input lines")
		}
		log.Printf("sfpread seeking to %d\n", offset)
		_, err := fp.Seek(offset, 0)
		if err != nil {
			if err == io.EOF {
				return klns, 0, err
			}
			log.Fatal("flreadn: ", err)
		}
	}
	if lpo != 0 {
		buf := make([]byte, reclen*lpo)
		if bl, err = io.ReadFull(fp, buf); err != nil {
			if err != io.EOF {
				log.Fatal("flreadn: ", err)
			}
		}
		for i := 0; i < bl/reclen; i++ {
			var kln kvalline
			kln.line = buf[i*reclen : (i+1)*reclen]
			if keyoff != 0 {
				kln.key = kln.line[keyoff:]
				if keylen != 0 {
					kln.key = kln.line[keyoff : keyoff+keylen]
				}
			}
			klns = append(klns, kln)
			nr++
		}
		if nr < lpo {
			return klns, 0, nil
		} else {
			return klns, int64(bl / reclen), errors.New("flreadn lpo")
		}
	} else {
		for {
			buf := make([]byte, reclen)
			if bl, err = io.ReadFull(fp, buf); err != nil {
				if err == io.EOF {
					return klns, 0, err
				}
				log.Fatal("flreadn: ", err)
			}
			if bl == 0 {
				return klns, 0, err
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
}

func vlscann(fp *os.File, offset int64, keyoff int, keylen int, lpo int) (kvallines, int64, error) {

	var klns kvallines
	var nr int // number records read

	if offset != 0 {
		if fp.Name() == "/dev/stdin" {
			log.Fatal("vlscann(stdin) offset ", offset)
		}
		_, err := fp.Seek(offset, 0)
		if err != nil {
			if err == io.EOF {
				return klns, 0, err
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
		nr++
		if lpo != 0 && nr >= lpo {
			offset, err := fp.Seek(0, 1)
			if err != nil {
				log.Fatal("vlscann Seek ", err)
			}
			return klns, offset, errors.New("vlscann lpo")
		}
	}
	e := scanner.Err()
	if e != nil {
		log.Fatal("vlscann: ", e)
	}
	return klns, 0, nil

}
