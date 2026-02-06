package aggregator

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileReportWriter_WriteReports(t *testing.T) {
	metrics := map[string]*CampaignMetrics{
		"camp1": {CampaignID: "camp1", TotalImpressions: 1000, TotalClicks: 100, TotalSpend: 500.00, TotalConversions: 10},
		"camp2": {CampaignID: "camp2", TotalImpressions: 2000, TotalClicks: 50, TotalSpend: 200.00, TotalConversions: 20},
	}

	dir := t.TempDir()
	w := NewFileReportWriter(dir)

	if err := w.WriteReports(metrics); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify files were created and are non-empty.
	for _, name := range []string{"top10_ctr.csv", "top10_cpa.csv"} {
		path := filepath.Join(dir, name)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected %s to exist: %v", name, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("expected %s to be non-empty", name)
		}
	}
}

func TestWriteTopCTR_Ranking(t *testing.T) {
	all := []*CampaignMetrics{
		{CampaignID: "low", TotalImpressions: 1000, TotalClicks: 10},
		{CampaignID: "high", TotalImpressions: 1000, TotalClicks: 100},
		{CampaignID: "mid", TotalImpressions: 1000, TotalClicks: 50},
	}

	var buf bytes.Buffer
	if err := writeTopCTR(&buf, all); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 { // header + 3 data rows
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}

	// First data row should be "high" (highest CTR).
	if !strings.HasPrefix(lines[1], "high,") {
		t.Errorf("expected first data row to be 'high', got %s", lines[1])
	}
	// Last data row should be "low" (lowest CTR).
	if !strings.HasPrefix(lines[3], "low,") {
		t.Errorf("expected last data row to be 'low', got %s", lines[3])
	}
}

func TestWriteTopCPA_ExcludesZeroConversions(t *testing.T) {
	all := []*CampaignMetrics{
		{CampaignID: "has_conv", TotalSpend: 100.00, TotalConversions: 10},
		{CampaignID: "no_conv", TotalSpend: 200.00, TotalConversions: 0},
	}

	var buf bytes.Buffer
	if err := writeTopCPA(&buf, all); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "no_conv") {
		t.Error("expected zero-conversion campaign to be excluded")
	}
	if !strings.Contains(output, "has_conv") {
		t.Error("expected campaign with conversions to be included")
	}
}

func TestWriteTopCPA_Ranking(t *testing.T) {
	all := []*CampaignMetrics{
		{CampaignID: "expensive", TotalSpend: 1000.00, TotalConversions: 10}, // CPA = 100
		{CampaignID: "cheap", TotalSpend: 100.00, TotalConversions: 10},     // CPA = 10
		{CampaignID: "mid", TotalSpend: 500.00, TotalConversions: 10},       // CPA = 50
	}

	var buf bytes.Buffer
	if err := writeTopCPA(&buf, all); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 { // header + 3 data rows
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}

	// First data row should be "cheap" (lowest CPA).
	if !strings.HasPrefix(lines[1], "cheap,") {
		t.Errorf("expected first data row to be 'cheap', got %s", lines[1])
	}
	// Last data row should be "expensive" (highest CPA).
	if !strings.HasPrefix(lines[3], "expensive,") {
		t.Errorf("expected last data row to be 'expensive', got %s", lines[3])
	}
}
