package aggregator

import "io"

// Processor reads and aggregates ad performance data from input
// into the provided MetricsStore.
type Processor interface {
	Process(r io.Reader, store MetricsStore) error
}

// ReportWriter generates reports from the provided MetricsStore.
type ReportWriter interface {
	WriteReports(store MetricsStore) error
}
