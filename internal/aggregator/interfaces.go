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

// MetricsStore owns the accumulation (write path) and top-K retrieval
// (read path) of campaign metrics.
type MetricsStore interface {
	// Add accumulates a single row of metrics for the given campaign.
	Add(campaignID string, impressions, clicks int64, spend float64, conversions int64)

	// TopKByCTR returns the top k campaigns sorted by CTR descending.
	TopKByCTR(k int) []*CampaignMetrics

	// TopKByCPA returns the top k campaigns sorted by CPA ascending,
	// excluding campaigns with zero conversions.
	TopKByCPA(k int) []*CampaignMetrics

	// Len returns the number of distinct campaigns in the store.
	Len() int
}
