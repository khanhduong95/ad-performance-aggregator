package aggregator

import "testing"

func TestTopKByCTR_Ranking(t *testing.T) {
	metrics := map[string]*CampaignMetrics{
		"low":  {CampaignID: "low", TotalImpressions: 1000, TotalClicks: 10},
		"high": {CampaignID: "high", TotalImpressions: 1000, TotalClicks: 100},
		"mid":  {CampaignID: "mid", TotalImpressions: 1000, TotalClicks: 50},
	}

	top := TopKByCTR(metrics, 10)
	if len(top) != 3 {
		t.Fatalf("expected 3 results, got %d", len(top))
	}
	if top[0].CampaignID != "high" {
		t.Errorf("expected first to be 'high', got %s", top[0].CampaignID)
	}
	if top[2].CampaignID != "low" {
		t.Errorf("expected last to be 'low', got %s", top[2].CampaignID)
	}
}

func TestTopKByCTR_LimitsToK(t *testing.T) {
	metrics := map[string]*CampaignMetrics{
		"a": {CampaignID: "a", TotalImpressions: 1000, TotalClicks: 100},
		"b": {CampaignID: "b", TotalImpressions: 1000, TotalClicks: 50},
		"c": {CampaignID: "c", TotalImpressions: 1000, TotalClicks: 10},
	}

	top := TopKByCTR(metrics, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
}

func TestTopKByLowestCPA_Ranking(t *testing.T) {
	metrics := map[string]*CampaignMetrics{
		"expensive": {CampaignID: "expensive", TotalSpend: 1000.00, TotalConversions: 10}, // CPA = 100
		"cheap":     {CampaignID: "cheap", TotalSpend: 100.00, TotalConversions: 10},      // CPA = 10
		"mid":       {CampaignID: "mid", TotalSpend: 500.00, TotalConversions: 10},         // CPA = 50
	}

	top := TopKByLowestCPA(metrics, 10)
	if len(top) != 3 {
		t.Fatalf("expected 3 results, got %d", len(top))
	}
	if top[0].CampaignID != "cheap" {
		t.Errorf("expected first to be 'cheap', got %s", top[0].CampaignID)
	}
	if top[2].CampaignID != "expensive" {
		t.Errorf("expected last to be 'expensive', got %s", top[2].CampaignID)
	}
}

func TestTopKByLowestCPA_ExcludesZeroConversions(t *testing.T) {
	metrics := map[string]*CampaignMetrics{
		"has_conv": {CampaignID: "has_conv", TotalSpend: 100.00, TotalConversions: 10},
		"no_conv":  {CampaignID: "no_conv", TotalSpend: 200.00, TotalConversions: 0},
	}

	top := TopKByLowestCPA(metrics, 10)
	if len(top) != 1 {
		t.Fatalf("expected 1 result, got %d", len(top))
	}
	if top[0].CampaignID != "has_conv" {
		t.Errorf("expected 'has_conv', got %s", top[0].CampaignID)
	}
}

func TestTopKByLowestCPA_LimitsToK(t *testing.T) {
	metrics := map[string]*CampaignMetrics{
		"a": {CampaignID: "a", TotalSpend: 100.00, TotalConversions: 10},
		"b": {CampaignID: "b", TotalSpend: 200.00, TotalConversions: 10},
		"c": {CampaignID: "c", TotalSpend: 300.00, TotalConversions: 10},
	}

	top := TopKByLowestCPA(metrics, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
}
