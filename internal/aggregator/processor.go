package aggregator

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
)

// expectedHeader defines the columns we require, in order.
var expectedHeader = []string{
	"campaign_id", "impressions", "clicks", "spend", "conversions",
}

type csvProcessor struct {
	benchmark bool
}

// NewCSVProcessor returns a Processor that parses and aggregates
// ad performance CSV data.
func NewCSVProcessor(benchmark bool) Processor {
	return &csvProcessor{benchmark: benchmark}
}

// Process streams the CSV from r line-by-line and accumulates
// metrics into store. Memory usage is proportional to the
// number of distinct campaign IDs, not the input size.
func (p *csvProcessor) Process(r io.Reader, store MetricsStore) error {
	reader := csv.NewReader(r)
	reader.ReuseRecord = true // reuse the backing array across Read calls

	// --- read and validate header ---
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	colIndex, err := mapColumns(header)
	if err != nil {
		return err
	}

	// --- stream rows ---
	lineNum := 1 // 1 = header already read

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNum+1, err)
		}
		lineNum++

		if err := accumulateRow(store, record, colIndex, lineNum); err != nil {
			return err
		}
	}

	if p.benchmark {
		log.Printf("benchmark: parsed %d data rows", lineNum-1)
	}
	return nil
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

// accumulateRow parses a single CSV record and merges it into the store.
func accumulateRow(store MetricsStore, record []string, col columnIndex, lineNum int) error {
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

	store.Add(campaignID, impressions, clicks, spend, conversions)

	return nil
}
