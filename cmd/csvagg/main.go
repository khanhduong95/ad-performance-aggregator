package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ad-performance-aggregator/internal/aggregator"
)

func main() {
	input := flag.String("input", "", "path to input CSV file (required)")
	output := flag.String("output", "", "path to output directory (required)")
	flag.Parse()

	if *input == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "usage: csvagg --input <csv_path> --output <output_dir>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := run(*input, *output); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(input, output string) error {
	// TODO: validate input file exists and is readable before processing

	start := time.Now()
	fmt.Fprintf(os.Stderr, "processing %s ...\n", input)

	metrics, err := aggregator.Process(input)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "aggregated %d campaigns in %s\n", len(metrics), time.Since(start))

	if err := aggregator.WriteReports(metrics, output); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "reports written to %s/\n", output)
	return nil
}
