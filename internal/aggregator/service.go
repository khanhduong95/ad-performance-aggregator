package aggregator

import (
	"io"
	"log/slog"
	"time"
)

// Service orchestrates CSV processing and report generation.
type Service struct {
	processor Processor
	writer    ReportWriter
}

func NewService(p Processor, w ReportWriter) *Service {
	return &Service{processor: p, writer: w}
}

func (s *Service) Run(r io.Reader) error {
	store := NewInMemoryMetricsStore()

	t0 := time.Now()
	if err := s.processor.Process(r, store); err != nil {
		return err
	}
	slog.Debug("processing phase complete", "elapsed", time.Since(t0))

	t1 := time.Now()
	if err := s.writer.WriteReports(store); err != nil {
		return err
	}
	slog.Debug("report writing phase complete", "elapsed", time.Since(t1))

	return nil
}
