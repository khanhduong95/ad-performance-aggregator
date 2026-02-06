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

const topN = 10

type fileReportWriter struct {
	outputDir string
}

// NewFileReportWriter returns a ReportWriter that writes CSV reports
// to the given directory.
func NewFileReportWriter(outputDir string) ReportWriter {
	return &fileReportWriter{outputDir: outputDir}
}

// WriteReports produces top10_ctr.csv and top10_cpa.csv inside the output directory.
func (w *fileReportWriter) WriteReports(metrics map[string]*CampaignMetrics) error {
	if err := os.MkdirAll(w.outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	all := metricsSlice(metrics)

	ctrPath := filepath.Join(w.outputDir, "top10_ctr.csv")
	if err := writeToFile(ctrPath, all, writeTopCTR); err != nil {
		return err
	}

	cpaPath := filepath.Join(w.outputDir, "top10_cpa.csv")
	if err := writeToFile(cpaPath, all, writeTopCPA); err != nil {
		return err
	}

	return nil
}

// writeToFile creates a file at path and delegates writing to fn.
func writeToFile(path string, all []*CampaignMetrics, fn func(io.Writer, []*CampaignMetrics) error) error {
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

// writeTopCTR writes the top-10 campaigns by CTR (descending) to w.
func writeTopCTR(w io.Writer, all []*CampaignMetrics) error {
	sort.Slice(all, func(i, j int) bool {
		return all[i].CTR() > all[j].CTR()
	})

	n := topN
	if len(all) < n {
		n = len(all)
	}

	return writeCSV(w, ctrHeader, all[:n], ctrRow)
}

// writeTopCPA writes the top-10 campaigns by CPA (ascending) to w,
// excluding campaigns with zero conversions.
func writeTopCPA(w io.Writer, all []*CampaignMetrics) error {
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

	return writeCSV(w, cpaHeader, eligible[:n], cpaRow)
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
