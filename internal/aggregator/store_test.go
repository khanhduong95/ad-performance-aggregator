package aggregator

import (
	"strings"
	"testing"
)

func TestInMemoryStore_AddAndLen(t *testing.T) {
	s := NewInMemoryStore()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}

	s.Add("camp1", 100, 10, 50.0, 5)
	if s.Len() != 1 {
		t.Fatalf("expected 1, got %d", s.Len())
	}

	// Adding to the same campaign does not increase Len.
	s.Add("camp1", 200, 20, 100.0, 10)
	if s.Len() != 1 {
		t.Fatalf("expected 1, got %d", s.Len())
	}

	s.Add("camp2", 300, 30, 150.0, 15)
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestInMemoryStore_Accumulation(t *testing.T) {
	s := NewInMemoryStore()
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

func TestInMemoryStore_TopKByCTR_Ranking(t *testing.T) {
	s := NewInMemoryStore()
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

func TestInMemoryStore_TopKByCTR_Limit(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 5; i++ {
		s.Add(string(rune('A'+i)), 1000, int64((i+1)*10), 0, 0)
	}

	top := s.TopKByCTR(2)
	if len(top) != 2 {
		t.Fatalf("expected 2, got %d", len(top))
	}
}

func TestInMemoryStore_TopKByCPA_Ranking(t *testing.T) {
	s := NewInMemoryStore()
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

func TestInMemoryStore_TopKByCPA_ExcludesZeroConversions(t *testing.T) {
	s := NewInMemoryStore()
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

func TestInMemoryStore_TopKByCPA_Limit(t *testing.T) {
	s := NewInMemoryStore()
	for i := 0; i < 5; i++ {
		s.Add(string(rune('A'+i)), 0, 0, float64((i+1)*100), 10)
	}

	top := s.TopKByCPA(2)
	if len(top) != 2 {
		t.Fatalf("expected 2, got %d", len(top))
	}
}

func TestInMemoryStore_TopK_Empty(t *testing.T) {
	s := NewInMemoryStore()

	ctr := s.TopKByCTR(10)
	if len(ctr) != 0 {
		t.Errorf("expected empty CTR result, got %d", len(ctr))
	}

	cpa := s.TopKByCPA(10)
	if len(cpa) != 0 {
		t.Errorf("expected empty CPA result, got %d", len(cpa))
	}
}

func TestInMemoryStore_DerivedMetrics(t *testing.T) {
	s := NewInMemoryStore()
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

func TestInMemoryStore_SatisfiesInterface(t *testing.T) {
	// Compile-time check that InMemoryStore implements MetricsStore.
	var _ MetricsStore = NewInMemoryStore()

	// Also verify via strings.NewReader round-trip with processor.
	s := NewInMemoryStore()
	p := NewCSVProcessor()
	input := "campaign_id,impressions,clicks,spend,conversions\ncamp1,100,10,50.00,5\n"
	if err := p.Process(strings.NewReader(input), s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Errorf("expected 1 campaign, got %d", s.Len())
	}
}
