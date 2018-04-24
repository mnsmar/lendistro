package cmd

// Opts encapsulates common command line options.
type Opts struct {
	Input     []string `arg:"positional,required" help:"input file/s (STDIN if -)"`
	SkipZeros bool     `arg:"--skip-zeros" help:"skip lengths with no records"`
	NoHeader  bool     `arg:"" help:"do not print header line"`
	Delim     string   `arg:"" help:"delimiter for output [default: \\t]"`
}

// OutDescription returns the output description.
const OutDescription = "Output is a delimited file that reports the number of records (count) for each length. The output density column is calculated by dividing the count by the sum of all counts."
