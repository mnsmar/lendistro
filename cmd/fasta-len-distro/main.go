package main

import (
	"log"
	"os"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"github.com/mnsmar/lendistro"
	"github.com/mnsmar/lendistro/cmd"
)

const descr = "Measure the length distribution of the records in FASTA file/s. " +
	cmd.OutDescription

// Opts is the struct with the options that the program accepts.
type Opts struct{ cmd.Opts }

// Version returns the program version.
func (Opts) Version() string { return "fasta-len-distro 0.1" }

// Description returns an extended description of the program.
func (Opts) Description() string { return descr }

// reader encapsulates a fasta.Reader and satisfies lendistro.Reader.
type reader struct{ r *fasta.Reader }

// Read reads one record from r.
func (r *reader) Read() (lendistro.Record, error) { return r.r.Read() }

func main() {
	var opts Opts
	arg.MustParse(&opts)
	if opts.Delim == "" {
		opts.Delim = "\t"
	}

	// create a length distro.
	distro := lendistro.NewDistro()

	var wg sync.WaitGroup
	for _, input := range opts.Input {
		wg.Add(1)
		go func(input string) {
			// open input file.
			f, err := os.Open(input)
			if err != nil {
				log.Fatal(err)
			}

			// create BED reader.
			r, err := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNA)), nil
			if err != nil {
				log.Fatal(err)
			}

			// encapsulate BED reader in reader.
			br := &reader{r}

			// update distro.
			if err = distro.Update(br); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(input)
	}

	// wait for updates to finish.
	wg.Wait()

	// print length distribution to stdout
	err := distro.Write(
		os.Stdout, opts.SkipZeros, !opts.NoHeader, []rune(opts.Delim)[0])
	if err != nil {
		log.Fatal(err)
	}
}
