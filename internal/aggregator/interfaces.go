package aggregator

import "io"

type Processor interface {
	Process(r io.Reader, store MetricsStore) error
}

type ReportWriter interface {
	WriteReports(store MetricsStore) error
}

// MetricsStore owns the accumulation (write path) and top-K retrieval
// (read path) of campaign metrics.
type MetricsStore interface {
	Add(
		campaignID string,
		impressions, clicks int64,
		spend float64,
		conversions int64,
	)

	// TopKByCTR returns the top k campaigns sorted by CTR descending.
	TopKByCTR(k int) []*CampaignMetrics

	// TopKByCPA returns the top k campaigns sorted by CPA ascending,
	// excluding campaigns with zero conversions.
	TopKByCPA(k int) []*CampaignMetrics
}
