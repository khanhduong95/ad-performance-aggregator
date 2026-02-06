package aggregator

import "sort"

// TopKByCTR returns up to k campaigns with the highest click-through rate,
// sorted in descending CTR order.
func TopKByCTR(metrics map[string]*CampaignMetrics, k int) []*CampaignMetrics {
	all := metricsSlice(metrics)
	sort.Slice(all, func(i, j int) bool {
		return all[i].CTR() > all[j].CTR()
	})
	if k > len(all) {
		k = len(all)
	}
	return all[:k]
}

// TopKByLowestCPA returns up to k campaigns with the lowest cost per
// acquisition, sorted in ascending CPA order. Campaigns with zero
// conversions are excluded since CPA is undefined for them.
func TopKByLowestCPA(metrics map[string]*CampaignMetrics, k int) []*CampaignMetrics {
	eligible := make([]*CampaignMetrics, 0, len(metrics))
	for _, m := range metrics {
		if m.TotalConversions > 0 {
			eligible = append(eligible, m)
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

// metricsSlice converts the map to a slice for sorting.
func metricsSlice(m map[string]*CampaignMetrics) []*CampaignMetrics {
	s := make([]*CampaignMetrics, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
