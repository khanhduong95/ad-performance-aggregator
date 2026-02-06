package aggregator

import "io"

// Processor reads and aggregates ad performance data from input.
type Processor interface {
	Process(r io.Reader) (map[string]*CampaignMetrics, error)
}

// ReportWriter generates reports from aggregated campaign metrics.
type ReportWriter interface {
	WriteReports(metrics map[string]*CampaignMetrics) error
}
