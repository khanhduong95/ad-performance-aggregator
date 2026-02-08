package aggregator

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type fileReportWriter struct {
	outputDir string
	topK      int
}

// NewFileReportWriter returns a ReportWriter that writes CSV reports
// to the given directory. The topK parameter controls how many top campaigns
// to include in each report (defaults to 10 if <= 0).
func NewFileReportWriter(outputDir string, topK int) ReportWriter {
	if topK <= 0 {
		topK = 10
	}
	return &fileReportWriter{outputDir: outputDir, topK: topK}
}

// WriteReports produces top{K}_ctr.csv and top{K}_cpa.csv inside the output directory.
func (w *fileReportWriter) WriteReports(store MetricsStore) error {
	if err := os.MkdirAll(w.outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	ctrData := store.TopKByCTR(w.topK)
	ctrPath := filepath.Join(w.outputDir, fmt.Sprintf("top%d_ctr.csv", w.topK))
	if err := writeMetricsFile(ctrPath, ctrHeader, ctrData, ctrRow); err != nil {
		return err
	}
	log.Printf("benchmark: wrote %d campaigns to %s", len(ctrData), ctrPath)

	cpaData := store.TopKByCPA(w.topK)
	cpaPath := filepath.Join(w.outputDir, fmt.Sprintf("top%d_cpa.csv", w.topK))
	if err := writeMetricsFile(cpaPath, cpaHeader, cpaData, cpaRow); err != nil {
		return err
	}
	log.Printf("benchmark: wrote %d campaigns to %s", len(cpaData), cpaPath)

	return nil
}

// writeMetricsFile creates a file at path and writes header + rows as CSV.
func writeMetricsFile(path string, header []string, rows []*CampaignMetrics, toRow func(*CampaignMetrics) []string) error {
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
