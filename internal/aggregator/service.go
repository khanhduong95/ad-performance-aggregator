package aggregator

import (
	"io"
	"log"
	"time"
)

// Service orchestrates CSV processing and report generation.
type Service struct {
	processor Processor
	writer    ReportWriter
	benchmark bool
}

// NewService creates a Service with the given processor and writer.
func NewService(p Processor, w ReportWriter, benchmark bool) *Service {
	return &Service{processor: p, writer: w, benchmark: benchmark}
}

// Run processes CSV data from r and writes the generated reports.
func (s *Service) Run(r io.Reader) error {
	store := NewInMemoryMetricsStore()

	t0 := time.Now()
	if err := s.processor.Process(r, store); err != nil {
		return err
	}
	if s.benchmark {
		log.Printf("benchmark: processing phase completed in %s", time.Since(t0))
	}

	t1 := time.Now()
	if err := s.writer.WriteReports(store); err != nil {
		return err
	}
	if s.benchmark {
		log.Printf("benchmark: report writing phase completed in %s", time.Since(t1))
	}

	return nil
}
