package aggregator

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// expectedHeader defines the columns we require, in order.
var expectedHeader = []string{
	"campaign_id", "impressions", "clicks", "spend", "conversions",
}

// Process streams the CSV at path line-by-line and returns aggregated
// metrics keyed by campaign_id. Memory usage is proportional to the
// number of distinct campaign IDs, not the file size.
func Process(path string) (map[string]*CampaignMetrics, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.ReuseRecord = true // reuse the backing array across Read calls

	// --- read and validate header ---
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	colIndex, err := mapColumns(header)
	if err != nil {
		return nil, err
	}

	// --- stream rows ---
	metrics := make(map[string]*CampaignMetrics)
	lineNum := 1 // 1 = header already read

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum+1, err)
		}
		lineNum++

		if err := accumulateRow(metrics, record, colIndex, lineNum); err != nil {
			// TODO: make error-on-bad-row vs skip-and-log configurable
			return nil, err
		}
	}

	return metrics, nil
}

// columnIndex holds the positional index for each required column.
type columnIndex struct {
	campaignID  int
	impressions int
	clicks      int
	spend       int
	conversions int
}

// mapColumns resolves header names to column positions.
func mapColumns(header []string) (columnIndex, error) {
	idx := columnIndex{-1, -1, -1, -1, -1}
	for i, name := range header {
		switch name {
		case "campaign_id":
			idx.campaignID = i
		case "impressions":
			idx.impressions = i
		case "clicks":
			idx.clicks = i
		case "spend":
			idx.spend = i
		case "conversions":
			idx.conversions = i
		}
	}
	// Verify all required columns were found.
	if idx.campaignID < 0 || idx.impressions < 0 || idx.clicks < 0 ||
		idx.spend < 0 || idx.conversions < 0 {
		return idx, fmt.Errorf("missing required columns; need %v, got %v", expectedHeader, header)
	}
	return idx, nil
}

// accumulateRow parses a single CSV record and merges it into metrics.
func accumulateRow(metrics map[string]*CampaignMetrics, record []string, col columnIndex, lineNum int) error {
	campaignID := record[col.campaignID]
	if campaignID == "" {
		return fmt.Errorf("line %d: empty campaign_id", lineNum)
	}

	impressions, err := strconv.ParseInt(record[col.impressions], 10, 64)
	if err != nil {
		return fmt.Errorf("line %d: bad impressions %q: %w", lineNum, record[col.impressions], err)
	}

	clicks, err := strconv.ParseInt(record[col.clicks], 10, 64)
	if err != nil {
		return fmt.Errorf("line %d: bad clicks %q: %w", lineNum, record[col.clicks], err)
	}

	spend, err := strconv.ParseFloat(record[col.spend], 64)
	if err != nil {
		return fmt.Errorf("line %d: bad spend %q: %w", lineNum, record[col.spend], err)
	}

	conversions, err := strconv.ParseInt(record[col.conversions], 10, 64)
	if err != nil {
		return fmt.Errorf("line %d: bad conversions %q: %w", lineNum, record[col.conversions], err)
	}

	// TODO: validate non-negative values

	m, ok := metrics[campaignID]
	if !ok {
		m = &CampaignMetrics{CampaignID: campaignID}
		metrics[campaignID] = m
	}
	m.TotalImpressions += impressions
	m.TotalClicks += clicks
	m.TotalSpend += spend
	m.TotalConversions += conversions

	return nil
}
