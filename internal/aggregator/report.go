package aggregator

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
func (w *fileReportWriter) WriteReports(metrics map[string]*CampaignMetrics) error {
	if err := os.MkdirAll(w.outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	all := metricsSlice(metrics)

	ctrPath := filepath.Join(w.outputDir, fmt.Sprintf("top%d_ctr.csv", w.topK))
	if err := w.writeToFile(ctrPath, all, w.writeTopCTR); err != nil {
		return err
	}

	cpaPath := filepath.Join(w.outputDir, fmt.Sprintf("top%d_cpa.csv", w.topK))
	if err := w.writeToFile(cpaPath, all, w.writeTopCPA); err != nil {
		return err
	}

	return nil
}

// writeToFile creates a file at path and delegates writing to fn.
func (w *fileReportWriter) writeToFile(path string, all []*CampaignMetrics, fn func(io.Writer, []*CampaignMetrics) error) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	if err := fn(f, all); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
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

// writeTopCTR writes the top-K campaigns by CTR (descending) to out.
func (w *fileReportWriter) writeTopCTR(out io.Writer, all []*CampaignMetrics) error {
	sort.Slice(all, func(i, j int) bool {
		return all[i].CTR() > all[j].CTR()
	})

	n := w.topK
	if len(all) < n {
		n = len(all)
	}

	return writeCSV(out, ctrHeader, all[:n], ctrRow)
}

// writeTopCPA writes the top-K campaigns by CPA (ascending) to out,
// excluding campaigns with zero conversions.
func (w *fileReportWriter) writeTopCPA(out io.Writer, all []*CampaignMetrics) error {
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

	n := w.topK
	if len(eligible) < n {
		n = len(eligible)
	}

	return writeCSV(out, cpaHeader, eligible[:n], cpaRow)
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
