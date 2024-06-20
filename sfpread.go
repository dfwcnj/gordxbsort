package main

import (
    "bufio"
    "io"
    "log"
    "os"
)

//type kvalline struct {
//        key  []byte
//        line []byte
//}
//type kvallines []kvalline

// flread(fp, reclen, keyoff, keylen, lps)
// reads up to lps fixed length records of reclen length from file ptr fp
// returns a slice of kvalline, the current offset, and err
func flreadn(fp *os.File, reclen int, keyoff int, keylen int, lps int) (k kvallines, offset int64, err error) {

    var klns kvallines
    var nr int              // number records read
    buf := make([]byte, reclen)

    for {
        if _, err := io.ReadFull(fp, buf); err != nil {
            if err == io.EOF {
                return klns, 0, err
            }
            log.Fatal(err)
        }
        var kln kvalline
        // to avoid having to make buf in the loop
        // mistake??
        copy(kln.line, buf)
        kln.key = kln.line
        if keyoff != 0 {
            kln.key = kln.line[keyoff:]
            if keylen != 0 {
                kln.key = kln.line[keyoff:keyoff+keylen]
            }
        }
        klns = append(klns, kln)
        nr++
        if nr >= lps {
            offset, err := fp.Seek(0, 1)
            if err != nil {
                log.Fatal(err)
            }
            return klns, offset, nil
        }
    }

}

func fscann(fp *os.File, keyoff int, keylen int, lps int) (k kvallines, offset int64, err error) {

    var klns kvallines
    var nr int              // number records read

    scanner := bufio.NewScanner(fp)
    for {
        l := scanner.Text()
        bln := []byte(l)
        var kln kvalline
        kln.line = bln
        kln.key = bln
        if keyoff != 0 {
            kln.key = kln.line[keyoff:]
            if keylen != 0 {
                kln.key = kln.line[keyoff:keyoff+keylen]
            }
        }
        klns = append(klns, kln)
        nr++
        if nr >= lps {
            offset, err := fp.Seek(0, 1)
            if err != nil {
                log.Fatal(err)
            }
            return klns, offset, nil
        }
    }

}

