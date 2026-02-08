package aggregator

import (
	"strings"
	"testing"
)

func TestCSVProcessor_BasicAggregation(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
camp1,1000,50,100.00,10
camp1,500,25,50.00,5
`
	p := NewCSVProcessor(false)
	store := NewInMemoryMetricsStore()
	if err := p.Process(strings.NewReader(input), store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	all := store.TopKByCTR(100)
	if len(all) != 1 {
		t.Fatalf("expected 1 campaign, got %d", len(all))
	}

	m := findByCampaignID(all, "camp1")
	if m == nil {
		t.Fatal("camp1 not found")
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

func TestCSVProcessor_MultipleCampaigns(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
camp1,1000,50,100.00,10
camp2,2000,100,200.00,20
camp3,3000,150,300.00,30
`
	p := NewCSVProcessor(false)
	store := NewInMemoryMetricsStore()
	if err := p.Process(strings.NewReader(input), store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	all := store.TopKByCTR(100)
	if len(all) != 3 {
		t.Fatalf("expected 3 campaigns, got %d", len(all))
	}
	for _, id := range []string{"camp1", "camp2", "camp3"} {
		if findByCampaignID(all, id) == nil {
			t.Errorf("campaign %s not found", id)
		}
	}
}

func TestCSVProcessor_MissingHeader(t *testing.T) {
	input := `campaign_id,impressions,clicks
camp1,1000,50
`
	p := NewCSVProcessor(false)
	store := NewInMemoryMetricsStore()
	err := p.Process(strings.NewReader(input), store)
	if err == nil {
		t.Fatal("expected error for missing columns")
	}
}

func TestCSVProcessor_BadImpressions(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
camp1,not_a_number,50,100.00,10
`
	p := NewCSVProcessor(false)
	store := NewInMemoryMetricsStore()
	err := p.Process(strings.NewReader(input), store)
	if err == nil {
		t.Fatal("expected error for bad impressions value")
	}
}

func TestCSVProcessor_EmptyCampaignID(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
,1000,50,100.00,10
`
	p := NewCSVProcessor(false)
	store := NewInMemoryMetricsStore()
	err := p.Process(strings.NewReader(input), store)
	if err == nil {
		t.Fatal("expected error for empty campaign_id")
	}
}

func TestCSVProcessor_HeaderOnly(t *testing.T) {
	input := "campaign_id,impressions,clicks,spend,conversions\n"
	p := NewCSVProcessor(false)
	store := NewInMemoryMetricsStore()
	if err := p.Process(strings.NewReader(input), store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n := len(store.TopKByCTR(100)); n != 0 {
		t.Errorf("expected 0 campaigns, got %d", n)
	}
}

func TestCSVProcessor_DerivedMetrics(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
camp1,1000,100,500.00,50
`
	p := NewCSVProcessor(false)
	store := NewInMemoryMetricsStore()
	if err := p.Process(strings.NewReader(input), store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := findByCampaignID(store.TopKByCTR(100), "camp1")
	wantCTR := 0.1 // 100/1000
	if m.CTR() != wantCTR {
		t.Errorf("CTR: got %f, want %f", m.CTR(), wantCTR)
	}

	wantCPA := 10.0 // 500/50
	if m.CPA() != wantCPA {
		t.Errorf("CPA: got %f, want %f", m.CPA(), wantCPA)
	}
}
