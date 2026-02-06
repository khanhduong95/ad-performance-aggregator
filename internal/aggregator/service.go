package aggregator

import "io"

// Service orchestrates CSV processing and report generation.
type Service struct {
	processor CSVProcessor
	writer    ReportWriter
}

// NewService creates a Service with the given processor and writer.
func NewService(p CSVProcessor, w ReportWriter) *Service {
	return &Service{processor: p, writer: w}
}

// Run processes CSV data from r and writes the generated reports.
func (s *Service) Run(r io.Reader) (map[string]*CampaignMetrics, error) {
	metrics, err := s.processor.Process(r)
	if err != nil {
		return nil, err
	}
	if err := s.writer.WriteReports(metrics); err != nil {
		return nil, err
	}
	return metrics, nil
}
