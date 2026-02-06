package aggregator

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// WriteReports produces top-K CTR and CPA CSV reports inside outputDir.
func WriteReports(outputDir string, metrics map[string]*CampaignMetrics, k int) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	topCTR := TopKByCTR(metrics, k)
	ctrPath := filepath.Join(outputDir, "top10_ctr.csv")
	if err := writeToFile(ctrPath, topCTR, ctrHeader, ctrRow); err != nil {
		return err
	}

	topCPA := TopKByLowestCPA(metrics, k)
	cpaPath := filepath.Join(outputDir, "top10_cpa.csv")
	if err := writeToFile(cpaPath, topCPA, cpaHeader, cpaRow); err != nil {
		return err
	}

	return nil
}

// writeToFile creates a file at path and writes ranked rows as CSV.
func writeToFile(path string, rows []*CampaignMetrics, header []string, toRow func(*CampaignMetrics) []string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	if err := writeCSV(f, header, rows, toRow); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
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

// writeCSV writes a header + data rows to w.
func writeCSV(w io.Writer, header []string, rows []*CampaignMetrics, toRow func(*CampaignMetrics) []string) error {
	cw := csv.NewWriter(w)

	if err := cw.Write(header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	for _, m := range rows {
		if err := cw.Write(toRow(m)); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
