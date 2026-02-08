package aggregator

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

type fileReportWriter struct {
	outputDir string
	topK      int
}

func NewFileReportWriter(outputDir string, topK int) ReportWriter {
	if topK <= 0 {
		topK = 10
	}
	return &fileReportWriter{outputDir: outputDir, topK: topK}
}

func (w *fileReportWriter) WriteReports(store MetricsStore) error {
	if err := os.MkdirAll(w.outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	ctrData := store.TopKByCTR(w.topK)
	ctrPath := filepath.Join(w.outputDir, fmt.Sprintf("top%d_ctr.csv", w.topK))
	if err := writeMetricsFile(ctrPath, reportHeader, ctrData, fullRow); err != nil {
		return err
	}
	slog.Debug("wrote report", "path", ctrPath, "campaigns", len(ctrData))

	cpaData := store.TopKByCPA(w.topK)
	cpaPath := filepath.Join(w.outputDir, fmt.Sprintf("top%d_cpa.csv", w.topK))
	if err := writeMetricsFile(cpaPath, reportHeader, cpaData, fullRow); err != nil {
		return err
	}
	slog.Debug("wrote report", "path", cpaPath, "campaigns", len(cpaData))

	return nil
}

func writeMetricsFile(
	path string,
	header []string,
	rows []*CampaignMetrics,
	toRow func(*CampaignMetrics) []string,
) error {
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

var reportHeader = []string{
	"campaign_id", "total_impressions", "total_clicks",
	"total_spend", "total_conversions", "CTR", "CPA",
}

func fullRow(m *CampaignMetrics) []string {
	cpa := ""
	if m.TotalConversions > 0 {
		cpa = strconv.FormatFloat(m.CPA(), 'f', 2, 64)
	}
	return []string{
		m.CampaignID,
		strconv.FormatInt(m.TotalImpressions, 10),
		strconv.FormatInt(m.TotalClicks, 10),
		strconv.FormatFloat(m.TotalSpend, 'f', 2, 64),
		strconv.FormatInt(m.TotalConversions, 10),
		strconv.FormatFloat(m.CTR(), 'f', 4, 64),
		cpa,
	}
}

func writeCSV(
	w io.Writer,
	header []string,
	rows []*CampaignMetrics,
	toRow func(*CampaignMetrics) []string,
) error {
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
