package probably

import (
	"fmt"
	"hash/crc32"
	"math"
)

// A count-min sketcher.
type Sketch [][]uint64

// Create a new count-min sketch with the given width and depth.
func NewSketch(w, d int) *Sketch {
	if d < 1 || w < 1 {
		panic("Dimensions must be positive")
	}

	rv := make(Sketch, w)
	for i := 0; i < w; i++ {
		rv[i] = make([]uint64, d)
	}

	return &rv
}

func (s Sketch) String() string {
	return fmt.Sprintf("{Sketch %dx%d}", len(s), len(s[0]))
}

func hash(s string, d, lim int) int {
	h1 := crc32.Update(0, crc32.IEEETable, []byte(s))
	h2 := int(crc32.Update(h1, crc32.IEEETable, []byte{byte(d)}))

	rv := h2 % lim
	if rv < 0 {
		rv = 0 - rv
	}
	return rv
}

// Increment the count for the given input.
func (s *Sketch) Increment(h string) {
	d := len((*s)[0])
	w := len(*s)
	for i := 0; i < d; i++ {
		(*s)[hash(h, i, w)][i]++
	}
}

// Get the estimate count for the given input.
func (s Sketch) Count(h string) uint64 {
	var min uint64 = math.MaxUint64
	d := len(s[0])
	w := len(s)
	for i := 0; i < d; i++ {
		v := s[hash(h, i, w)][i]
		if v < min {
			min = v
		}
	}
	return min
}
