package aggregator

import "sort"

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

// InMemoryStore is an in-memory implementation of MetricsStore backed
// by a plain map. It is the only concrete backend.
type InMemoryStore struct {
	m map[string]*CampaignMetrics
}

// NewInMemoryStore creates a new empty InMemoryStore.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{m: make(map[string]*CampaignMetrics)}
}

func (s *InMemoryStore) Add(campaignID string, impressions, clicks int64, spend float64, conversions int64) {
	cm, ok := s.m[campaignID]
	if !ok {
		cm = &CampaignMetrics{CampaignID: campaignID}
		s.m[campaignID] = cm
	}
	cm.TotalImpressions += impressions
	cm.TotalClicks += clicks
	cm.TotalSpend += spend
	cm.TotalConversions += conversions
}

func (s *InMemoryStore) TopKByCTR(k int) []*CampaignMetrics {
	all := s.all()
	sort.Slice(all, func(i, j int) bool {
		return all[i].CTR() > all[j].CTR()
	})
	if k > len(all) {
		k = len(all)
	}
	return all[:k]
}

func (s *InMemoryStore) TopKByCPA(k int) []*CampaignMetrics {
	eligible := make([]*CampaignMetrics, 0, len(s.m))
	for _, cm := range s.m {
		if cm.TotalConversions > 0 {
			eligible = append(eligible, cm)
		}
	}
	sort.Slice(eligible, func(i, j int) bool {
		return eligible[i].CPA() < eligible[j].CPA()
	})
	if k > len(eligible) {
		k = len(eligible)
	}
	return eligible[:k]
}

func (s *InMemoryStore) Len() int {
	return len(s.m)
}

func (s *InMemoryStore) all() []*CampaignMetrics {
	result := make([]*CampaignMetrics, 0, len(s.m))
	for _, v := range s.m {
		result = append(result, v)
	}
	return result
}
