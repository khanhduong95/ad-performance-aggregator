package aggregator

import "io"

// Service orchestrates CSV processing and report generation.
type Service struct {
	processor Processor
	writer    ReportWriter
}

// NewService creates a Service with the given processor and writer.
func NewService(p Processor, w ReportWriter) *Service {
	return &Service{processor: p, writer: w}
}

// Run processes CSV data from r and writes the generated reports.
func (s *Service) Run(r io.Reader) error {
	store := NewInMemoryMetricsStore()

	if err := s.processor.Process(r, store); err != nil {
		return err
	}
	return s.writer.WriteReports(store)
}
