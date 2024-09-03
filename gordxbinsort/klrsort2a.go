package gordxbinsort

// shæmelessly plagiarized from
// https://github.com/rsc/tmp/blob/master/rsort/rsort.go sortWithTmp
// without the goto, using byte slices instead of strings
// commented so that I can understand it better
// the rsc github directory is a great place to see an analysis of radix
// sort compared to slices.sort.
// if you are interested in radix sort, his three functions in rsort.go
// are highly instructive.
// same license as rsc code - BSD

func KLrsort2a(klns kvallines) {
	klrsort2array(klns, make(kvallines, len(klns)), 0)
}

func klrsort2array(klns, klns2 kvallines, ix int) {
	if len(klns) < 16 {
		binsertionsort(klns) // insertion sort
		return
	}

	// lines are sorted into bins based on byte at offset ix
	// compute bin sizes
	var sizes, ends [257]int // bin sizes and end indices
	fbin := 256              // first assignable bin
	lbin := 1                // last assignable bin
	for _, s := range klns {
		c := 0 // for shorties
		if ix < len(s.key) {
			c = int(s.key[ix]) + 1
		}
		sizes[c]++

		if sizes[c] == 1 && c > 0 { // bin is newly assigned
			fbin = min(fbin, c)
			lbin = max(lbin, c)
		}
	}

	// compute bin ends for move into bins below
	off := sizes[0] // skip unassignable bine
	ends[0] = off
	for i := fbin; i <= lbin; i++ {
		n := sizes[i]
		if n == 0 {
			continue
		}
		off += sizes[i]
		ends[i] = off // bin end
	}

	// klns simplifies moving lines into their bins
	copy(klns2, klns)

	// move lines into bins
	for i := len(klns) - 1; i >= 0; i-- {
		s := klns2[i]
		c := 0 // bin for for shorties
		if ix < len(s.key) {
			c = int(s.key[ix]) + 1 // destination bin
		}
		ends[c]--         // ends contained size, not last offset
		klns[ends[c]] = s // move to home
	}

	// recurse binning the lines with the next byte offset
	off = sizes[0]
	for c := fbin; c <= lbin; c++ {
		n := sizes[c] // size of bin
		if c > 0 && n > 1 {
			klrsort2array(klns[off:off+n], klns2, ix+1)
		}
		off += n // offset to next bin
	}
}
