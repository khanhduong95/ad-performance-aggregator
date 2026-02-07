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
// It returns the number of distinct campaigns aggregated.
func (s *Service) Run(r io.Reader) (int, error) {
	store := NewInMemoryStore()

	if err := s.processor.Process(r, store); err != nil {
		return 0, err
	}
	if err := s.writer.WriteReports(store); err != nil {
		return 0, err
	}
	return store.Len(), nil
}
