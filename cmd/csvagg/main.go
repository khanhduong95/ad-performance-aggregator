package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/khanhduong95/ad-performance-aggregator/internal/aggregator"
)

func main() {
	input := flag.String("input", "", "path to input CSV file (required)")
	output := flag.String("output", "", "path to output directory (required)")
	topK := flag.Int("topk", 10, "number of top campaigns to include in reports (default: 10)")
	benchmark := flag.Bool("benchmark", false, "enable benchmark timing logs on stderr")
	flag.Parse()

	if *benchmark {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	if *input == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "usage: csvagg --input <csv_path> --output <output_dir> [--topk <number>] [--benchmark]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := run(*input, *output, *topK); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(input, output string, topK int) error {
	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer f.Close()

	start := time.Now()
	fmt.Fprintf(os.Stderr, "processing %s ...\n", input)

	svc := aggregator.NewService(
		aggregator.NewCSVProcessor(),
		aggregator.NewFileReportWriter(output, topK),
	)

	if err := svc.Run(f); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "done in %s\n", time.Since(start))
	fmt.Fprintf(os.Stderr, "reports written to %s/\n", output)
	return nil
}
