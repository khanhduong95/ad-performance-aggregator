package aggregator

import (
	"strings"
	"testing"
)

func TestAggregate_BasicAggregation(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
camp1,1000,50,100.00,10
camp1,500,25,50.00,5
`
	metrics, err := Aggregate(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("expected 1 campaign, got %d", len(metrics))
	}

	m := metrics["camp1"]
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

func TestAggregate_MultipleCampaigns(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
camp1,1000,50,100.00,10
camp2,2000,100,200.00,20
camp3,3000,150,300.00,30
`
	metrics, err := Aggregate(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(metrics) != 3 {
		t.Fatalf("expected 3 campaigns, got %d", len(metrics))
	}

	for _, id := range []string{"camp1", "camp2", "camp3"} {
		if metrics[id] == nil {
			t.Errorf("campaign %s not found", id)
		}
	}
}

func TestAggregate_MissingHeader(t *testing.T) {
	input := `campaign_id,impressions,clicks
camp1,1000,50
`
	_, err := Aggregate(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for missing columns")
	}
}

func TestAggregate_BadImpressions(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
camp1,not_a_number,50,100.00,10
`
	_, err := Aggregate(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for bad impressions value")
	}
}

func TestAggregate_EmptyCampaignID(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
,1000,50,100.00,10
`
	_, err := Aggregate(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for empty campaign_id")
	}
}

func TestAggregate_HeaderOnly(t *testing.T) {
	input := "campaign_id,impressions,clicks,spend,conversions\n"
	metrics, err := Aggregate(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metrics) != 0 {
		t.Errorf("expected 0 campaigns, got %d", len(metrics))
	}
}

func TestAggregate_DerivedMetrics(t *testing.T) {
	input := `campaign_id,impressions,clicks,spend,conversions
camp1,1000,100,500.00,50
`
	metrics, err := Aggregate(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := metrics["camp1"]
	wantCTR := 0.1 // 100/1000
	if m.CTR() != wantCTR {
		t.Errorf("CTR: got %f, want %f", m.CTR(), wantCTR)
	}

	wantCPA := 10.0 // 500/50
	if m.CPA() != wantCPA {
		t.Errorf("CPA: got %f, want %f", m.CPA(), wantCPA)
	}
}
