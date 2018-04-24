package lendistro

import (
	"encoding/csv"
	"io"
	"strconv"
	"sync"
)

const maxInt = int(^uint(0) >> 1) // The maximum int value.

// A Record is a generic element that has length.
type Record interface {
	Len() int
}

// A Reader is the common Record reader interface.
type Reader interface {
	Read() (Record, error)
}

// Distro is the length distribution. It stores the number of entries
// corresponding to every length. It is safe for concurrent use.
type Distro struct {
	counts         map[int]int
	minLen, maxLen int
	mux            sync.Mutex
}

// NewDistro returns a new, initialized Distro.
func NewDistro() *Distro {
	return &Distro{counts: make(map[int]int), minLen: maxInt}
}

// MaxLen returns the maximum length stored in d.
func (d *Distro) MaxLen() int {
	return d.maxLen
}

// MinLen returns the minimum length stored in d.
func (d *Distro) MinLen() int {
	return d.minLen
}

// Sum returns the total number of entries stored in d.
func (d *Distro) Sum() int {
	d.mux.Lock()
	s := d.sum()
	d.mux.Unlock()
	return s
}

// sum returns the total number of entries stored in d (not safe for
// concurrent use).
func (d *Distro) sum() int {
	sum := 0
	for _, v := range d.counts {
		sum += v
	}
	return sum
}

// Counts returns a map with the raw counts stored in d. Key is the length and
// value is the corresponding count.
func (d *Distro) Counts() map[int]int {
	d.mux.Lock()
	counts := make(map[int]int)
	for k, v := range d.counts {
		counts[k] = v
	}
	d.mux.Unlock()
	return counts
}

// Add adds n records of length len in Distro.
func (d *Distro) Add(len, n int) {
	d.mux.Lock()
	d.add(len, n)
	d.mux.Unlock()
}

// Add adds n elements at length len in Distro (not safe for concurrent use).
func (d *Distro) add(len, n int) {
	d.counts[len] += n
	if len > d.maxLen {
		d.maxLen = len
	}
	if len < d.minLen {
		d.minLen = len
	}
}

// Update updates d with the length distribution of the records in r.
func (d *Distro) Update(r Reader) error {
	// store in a local unsafe distro first, for performance.
	unsafe := NewDistro()
	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		len := rec.Len()
		unsafe.add(len, 1)
	}
	// safely add counts to d.
	for l, v := range unsafe.counts {
		d.Add(l, v)
	}
	return nil
}

// Write writes d in w and returns the number of bytes written and any error.
// It only writes lengths from min to max. Skips lengths with 0 intervals if
// skipZeros is set.
func (d *Distro) Write(
	w io.Writer, skipZeros, header bool, comma rune) error {

	csvW := csv.NewWriter(w)
	csvW.Comma = comma

	if header {
		if err := csvW.Write([]string{"len", "count", "density"}); err != nil {
			return err
		}
	}

	d.mux.Lock()
	sum := d.sum()
	for l := 0; l <= d.MaxLen(); l++ {
		v, ok := d.counts[l]
		if skipZeros && !ok {
			continue
		}
		lenStr := strconv.Itoa(l)
		valStr := strconv.Itoa(v)
		densStr := strconv.FormatFloat(float64(v)/float64(sum), 'f', -1, 32)
		if err := csvW.Write([]string{lenStr, valStr, densStr}); err != nil {
			return err
		}
	}
	d.mux.Unlock()

	csvW.Flush()
	return csvW.Error()
}
