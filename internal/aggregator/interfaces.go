package aggregator

import "io"

// CSVProcessor reads and aggregates ad performance data from CSV input.
type CSVProcessor interface {
	Process(r io.Reader) (map[string]*CampaignMetrics, error)
}

// ReportWriter generates reports from aggregated campaign metrics.
type ReportWriter interface {
	WriteReports(metrics map[string]*CampaignMetrics) error
}
