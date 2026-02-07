package aggregator

import (
	"strings"
	"testing"
)

func TestInMemoryMetricsStore_AddAndCount(t *testing.T) {
	s := NewInMemoryMetricsStore()
	if n := len(s.TopKByCTR(100)); n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}

	s.Add("camp1", 100, 10, 50.0, 5)
	if n := len(s.TopKByCTR(100)); n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}

	// Adding to the same campaign does not increase count.
	s.Add("camp1", 200, 20, 100.0, 10)
	if n := len(s.TopKByCTR(100)); n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}

	s.Add("camp2", 300, 30, 150.0, 15)
	if n := len(s.TopKByCTR(100)); n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}
}

func TestInMemoryMetricsStore_Accumulation(t *testing.T) {
	s := NewInMemoryMetricsStore()
	s.Add("camp1", 1000, 50, 100.00, 10)
	s.Add("camp1", 500, 25, 50.00, 5)

	all := s.TopKByCTR(100)
	if len(all) != 1 {
		t.Fatalf("expected 1 campaign, got %d", len(all))
	}

	m := all[0]
	if m.CampaignID != "camp1" {
		t.Errorf("expected camp1, got %s", m.CampaignID)
	}
	if m.TotalImpressions != 1500 {
		t.Errorf("impressions: got %d, want 1500", m.TotalImpressions)
	}
	if m.TotalClicks != 75 {
		t.Errorf("clicks: got %d, want 75", m.TotalClicks)
	}
	if m.TotalSpend != 150.0 {
		t.Errorf("spend: got %f, want 150.0", m.TotalSpend)
	}
	if m.TotalConversions != 15 {
		t.Errorf("conversions: got %d, want 15", m.TotalConversions)
	}
}

func TestInMemoryMetricsStore_TopKByCTR_Ranking(t *testing.T) {
	s := NewInMemoryMetricsStore()
	s.Add("low", 1000, 10, 0, 0)   // CTR 0.01
	s.Add("high", 1000, 100, 0, 0)  // CTR 0.10
	s.Add("mid", 1000, 50, 0, 0)    // CTR 0.05

	top := s.TopKByCTR(10)
	if len(top) != 3 {
		t.Fatalf("expected 3, got %d", len(top))
	}
	if top[0].CampaignID != "high" {
		t.Errorf("expected first to be 'high', got %s", top[0].CampaignID)
	}
	if top[2].CampaignID != "low" {
		t.Errorf("expected last to be 'low', got %s", top[2].CampaignID)
	}
}

func TestInMemoryMetricsStore_TopKByCTR_Limit(t *testing.T) {
	s := NewInMemoryMetricsStore()
	for i := 0; i < 5; i++ {
		s.Add(string(rune('A'+i)), 1000, int64((i+1)*10), 0, 0)
	}

	top := s.TopKByCTR(2)
	if len(top) != 2 {
		t.Fatalf("expected 2, got %d", len(top))
	}
}

func TestInMemoryMetricsStore_TopKByCPA_Ranking(t *testing.T) {
	s := NewInMemoryMetricsStore()
	s.Add("expensive", 0, 0, 1000.00, 10) // CPA = 100
	s.Add("cheap", 0, 0, 100.00, 10)      // CPA = 10
	s.Add("mid", 0, 0, 500.00, 10)        // CPA = 50

	top := s.TopKByCPA(10)
	if len(top) != 3 {
		t.Fatalf("expected 3, got %d", len(top))
	}
	if top[0].CampaignID != "cheap" {
		t.Errorf("expected first to be 'cheap', got %s", top[0].CampaignID)
	}
	if top[2].CampaignID != "expensive" {
		t.Errorf("expected last to be 'expensive', got %s", top[2].CampaignID)
	}
}

func TestInMemoryMetricsStore_TopKByCPA_ExcludesZeroConversions(t *testing.T) {
	s := NewInMemoryMetricsStore()
	s.Add("has_conv", 0, 0, 100.00, 10)
	s.Add("no_conv", 0, 0, 200.00, 0)

	top := s.TopKByCPA(10)
	if len(top) != 1 {
		t.Fatalf("expected 1, got %d", len(top))
	}
	if top[0].CampaignID != "has_conv" {
		t.Errorf("expected 'has_conv', got %s", top[0].CampaignID)
	}
}

func TestInMemoryMetricsStore_TopKByCPA_Limit(t *testing.T) {
	s := NewInMemoryMetricsStore()
	for i := 0; i < 5; i++ {
		s.Add(string(rune('A'+i)), 0, 0, float64((i+1)*100), 10)
	}

	top := s.TopKByCPA(2)
	if len(top) != 2 {
		t.Fatalf("expected 2, got %d", len(top))
	}
}

func TestInMemoryMetricsStore_TopK_Empty(t *testing.T) {
	s := NewInMemoryMetricsStore()

	ctr := s.TopKByCTR(10)
	if len(ctr) != 0 {
		t.Errorf("expected empty CTR result, got %d", len(ctr))
	}

	cpa := s.TopKByCPA(10)
	if len(cpa) != 0 {
		t.Errorf("expected empty CPA result, got %d", len(cpa))
	}
}

func TestInMemoryMetricsStore_DerivedMetrics(t *testing.T) {
	s := NewInMemoryMetricsStore()
	s.Add("camp1", 1000, 100, 500.00, 50)

	all := s.TopKByCTR(1)
	m := all[0]

	wantCTR := 0.1 // 100/1000
	if m.CTR() != wantCTR {
		t.Errorf("CTR: got %f, want %f", m.CTR(), wantCTR)
	}

	wantCPA := 10.0 // 500/50
	if m.CPA() != wantCPA {
		t.Errorf("CPA: got %f, want %f", m.CPA(), wantCPA)
	}
}

// findByCampaignID is a test helper that locates a campaign in a slice.
func findByCampaignID(metrics []*CampaignMetrics, id string) *CampaignMetrics {
	for _, m := range metrics {
		if m.CampaignID == id {
			return m
		}
	}
	return nil
}

func TestInMemoryMetricsStore_SatisfiesInterface(t *testing.T) {
	// Compile-time check that InMemoryMetricsStore implements MetricsStore.
	var _ MetricsStore = NewInMemoryMetricsStore()

	// Also verify via strings.NewReader round-trip with processor.
	s := NewInMemoryMetricsStore()
	p := NewCSVProcessor()
	input := "campaign_id,impressions,clicks,spend,conversions\ncamp1,100,10,50.00,5\n"
	if err := p.Process(strings.NewReader(input), s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n := len(s.TopKByCTR(100)); n != 1 {
		t.Errorf("expected 1 campaign, got %d", n)
	}
}
