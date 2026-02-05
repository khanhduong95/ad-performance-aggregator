package aggregator

import "fmt"

// CampaignMetrics holds the running totals for a single campaign_id.
// All fields are accumulated during the streaming pass; derived metrics
// (CTR, CPA) are computed on demand after aggregation is complete.
type CampaignMetrics struct {
	CampaignID       string
	TotalImpressions int64
	TotalClicks      int64
	TotalSpend       float64
	TotalConversions int64
}

// CTR returns the click-through rate (clicks / impressions).
// Returns 0 if there are no impressions.
func (m *CampaignMetrics) CTR() float64 {
	if m.TotalImpressions == 0 {
		return 0
	}
	return float64(m.TotalClicks) / float64(m.TotalImpressions)
}

// CPA returns cost per acquisition (spend / conversions).
// Returns 0 if there are no conversions. Callers that need to exclude
// zero-conversion campaigns should check TotalConversions before calling.
func (m *CampaignMetrics) CPA() float64 {
	if m.TotalConversions == 0 {
		return 0
	}
	return m.TotalSpend / float64(m.TotalConversions)
}

// String implements fmt.Stringer for debugging.
func (m *CampaignMetrics) String() string {
	return fmt.Sprintf("campaign=%s imp=%d click=%d spend=%.2f conv=%d ctr=%.6f cpa=%.2f",
		m.CampaignID, m.TotalImpressions, m.TotalClicks, m.TotalSpend,
		m.TotalConversions, m.CTR(), m.CPA())
}
