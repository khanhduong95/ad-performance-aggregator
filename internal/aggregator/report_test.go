package aggregator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileReportWriter_WriteReports(t *testing.T) {
	store := NewInMemoryMetricsStore()
	store.Add("camp1", 1000, 100, 500.00, 10)
	store.Add("camp2", 2000, 50, 200.00, 20)

	dir := t.TempDir()
	w := NewFileReportWriter(dir, 10)

	if err := w.WriteReports(store); err != nil {
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
	store := NewInMemoryMetricsStore()
	store.Add("low", 1000, 10, 5.00, 1)
	store.Add("high", 1000, 100, 50.00, 10)
	store.Add("mid", 1000, 50, 25.00, 5)

	dir := t.TempDir()
	w := NewFileReportWriter(dir, 10)

	if err := w.WriteReports(store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "top10_ctr.csv"))
	if err != nil {
		t.Fatalf("read ctr file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 4 { // header + 3 data rows
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}

	// Verify header has all 7 columns.
	if lines[0] != "campaign_id,total_impressions,total_clicks,total_spend,total_conversions,CTR,CPA" {
		t.Errorf("unexpected header: %s", lines[0])
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
	store := NewInMemoryMetricsStore()
	store.Add("has_conv", 0, 0, 100.00, 10)
	store.Add("no_conv", 0, 0, 200.00, 0)

	dir := t.TempDir()
	w := NewFileReportWriter(dir, 10)

	if err := w.WriteReports(store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "top10_cpa.csv"))
	if err != nil {
		t.Fatalf("read cpa file: %v", err)
	}

	output := string(data)
	if strings.Contains(output, "no_conv") {
		t.Error("expected zero-conversion campaign to be excluded")
	}
	if !strings.Contains(output, "has_conv") {
		t.Error("expected campaign with conversions to be included")
	}
}

func TestWriteTopCPA_Ranking(t *testing.T) {
	store := NewInMemoryMetricsStore()
	store.Add("expensive", 1000, 50, 1000.00, 10) // CPA = 100
	store.Add("cheap", 1000, 50, 100.00, 10)      // CPA = 10
	store.Add("mid", 1000, 50, 500.00, 10)        // CPA = 50

	dir := t.TempDir()
	w := NewFileReportWriter(dir, 10)

	if err := w.WriteReports(store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "top10_cpa.csv"))
	if err != nil {
		t.Fatalf("read cpa file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 4 { // header + 3 data rows
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}

	// Verify header has all 7 columns.
	if lines[0] != "campaign_id,total_impressions,total_clicks,total_spend,total_conversions,CTR,CPA" {
		t.Errorf("unexpected header: %s", lines[0])
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

func TestWriteTopCTR_CPANullForZeroConversions(t *testing.T) {
	store := NewInMemoryMetricsStore()
	store.Add("has_conv", 1000, 100, 500.00, 50) // CPA = 10.00
	store.Add("no_conv", 1000, 200, 300.00, 0)   // CPA should be empty

	dir := t.TempDir()
	w := NewFileReportWriter(dir, 10)

	if err := w.WriteReports(store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "top10_ctr.csv"))
	if err != nil {
		t.Fatalf("read ctr file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 { // header + 2 data rows
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	// "no_conv" has higher CTR (0.2 vs 0.1) so appears first.
	// Its CPA column (last field) should be empty.
	if !strings.HasSuffix(lines[1], ",") {
		t.Errorf("expected empty CPA for zero-conversion campaign, got %s", lines[1])
	}

	// "has_conv" should have CPA = 10.00.
	if !strings.HasSuffix(lines[2], ",10.00") {
		t.Errorf("expected CPA 10.00 for campaign with conversions, got %s", lines[2])
	}
}

func TestConfigurableTopK(t *testing.T) {
	store := NewInMemoryMetricsStore()
	store.Add("camp1", 1000, 10, 0, 0)  // CTR: 0.01
	store.Add("camp2", 1000, 20, 0, 0)  // CTR: 0.02
	store.Add("camp3", 1000, 30, 0, 0)  // CTR: 0.03
	store.Add("camp4", 1000, 40, 0, 0)  // CTR: 0.04
	store.Add("camp5", 1000, 50, 0, 0)  // CTR: 0.05

	dir := t.TempDir()
	// Test with topK = 2, should only return top 2 campaigns.
	w := NewFileReportWriter(dir, 2)

	if err := w.WriteReports(store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "top2_ctr.csv"))
	if err != nil {
		t.Fatalf("read ctr file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 { // header + 2 data rows
		t.Fatalf("expected 3 lines (header + 2 data rows), got %d", len(lines))
	}

	// Verify that only top 2 campaigns are included.
	if !strings.HasPrefix(lines[1], "camp5,") {
		t.Errorf("expected first data row to be 'camp5', got %s", lines[1])
	}
	if !strings.HasPrefix(lines[2], "camp4,") {
		t.Errorf("expected second data row to be 'camp4', got %s", lines[2])
	}
}

func TestConfigurableTopK_FileNames(t *testing.T) {
	store := NewInMemoryMetricsStore()
	store.Add("camp1", 1000, 100, 500.00, 10)

	dir := t.TempDir()
	w := NewFileReportWriter(dir, 5)

	if err := w.WriteReports(store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify files were created with the correct topK value in the name.
	for _, name := range []string{"top5_ctr.csv", "top5_cpa.csv"} {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected %s to exist: %v", name, err)
		}
	}

	// Verify old file names don't exist.
	for _, name := range []string{"top10_ctr.csv", "top10_cpa.csv"} {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			t.Errorf("expected %s to NOT exist", name)
		}
	}
}
