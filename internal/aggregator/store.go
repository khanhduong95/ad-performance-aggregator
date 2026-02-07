package aggregator

import "sort"

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
