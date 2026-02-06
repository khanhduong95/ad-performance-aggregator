package aggregator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteReports(t *testing.T) {
	metrics := map[string]*CampaignMetrics{
		"camp1": {CampaignID: "camp1", TotalImpressions: 1000, TotalClicks: 100, TotalSpend: 500.00, TotalConversions: 10},
		"camp2": {CampaignID: "camp2", TotalImpressions: 2000, TotalClicks: 50, TotalSpend: 200.00, TotalConversions: 20},
	}

	dir := t.TempDir()
	if err := WriteReports(dir, metrics, 10); err != nil {
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
