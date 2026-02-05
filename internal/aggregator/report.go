package aggregator

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const topN = 10

// WriteReports produces top10_ctr.csv and top10_cpa.csv inside outputDir.
func WriteReports(metrics map[string]*CampaignMetrics, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	all := metricsSlice(metrics)

	if err := writeTopCTR(all, outputDir); err != nil {
		return err
	}
	if err := writeTopCPA(all, outputDir); err != nil {
		return err
	}
	return nil
}

// metricsSlice converts the map to a slice for sorting.
func metricsSlice(m map[string]*CampaignMetrics) []*CampaignMetrics {
	s := make([]*CampaignMetrics, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

// writeTopCTR writes the top-10 campaigns by CTR (descending).
func writeTopCTR(all []*CampaignMetrics, outputDir string) error {
	// TODO: consider filtering out campaigns with zero impressions
	sort.Slice(all, func(i, j int) bool {
		return all[i].CTR() > all[j].CTR()
	})

	n := topN
	if len(all) < n {
		n = len(all)
	}

	path := filepath.Join(outputDir, "top10_ctr.csv")
	return writeCSV(path, ctrHeader, all[:n], ctrRow)
}

// writeTopCPA writes the top-10 campaigns by CPA (ascending),
// excluding campaigns with zero conversions.
func writeTopCPA(all []*CampaignMetrics, outputDir string) error {
	// Filter to campaigns that actually have conversions.
	eligible := make([]*CampaignMetrics, 0, len(all))
	for _, m := range all {
		if m.TotalConversions > 0 {
			eligible = append(eligible, m)
		}
	}

	sort.Slice(eligible, func(i, j int) bool {
		return eligible[i].CPA() < eligible[j].CPA()
	})

	n := topN
	if len(eligible) < n {
		n = len(eligible)
	}

	path := filepath.Join(outputDir, "top10_cpa.csv")
	return writeCSV(path, cpaHeader, eligible[:n], cpaRow)
}

var ctrHeader = []string{"campaign_id", "impressions", "clicks", "ctr"}
var cpaHeader = []string{"campaign_id", "spend", "conversions", "cpa"}

func ctrRow(m *CampaignMetrics) []string {
	return []string{
		m.CampaignID,
		strconv.FormatInt(m.TotalImpressions, 10),
		strconv.FormatInt(m.TotalClicks, 10),
		strconv.FormatFloat(m.CTR(), 'f', 6, 64),
	}
}

func cpaRow(m *CampaignMetrics) []string {
	return []string{
		m.CampaignID,
		strconv.FormatFloat(m.TotalSpend, 'f', 2, 64),
		strconv.FormatInt(m.TotalConversions, 10),
		strconv.FormatFloat(m.CPA(), 'f', 2, 64),
	}
}

// writeCSV is a small helper to write a header + N data rows.
func writeCSV(path string, header []string, rows []*CampaignMetrics, toRow func(*CampaignMetrics) []string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	if err := w.Write(header); err != nil {
		return fmt.Errorf("write header to %s: %w", path, err)
	}
	for _, m := range rows {
		if err := w.Write(toRow(m)); err != nil {
			return fmt.Errorf("write row to %s: %w", path, err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("flush %s: %w", path, err)
	}
	return nil
}
