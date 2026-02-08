# ad-performance-aggregator

A streaming CLI tool that processes large CSV files of advertising campaign performance data and produces ranked reports by CTR (Click-Through Rate) and CPA (Cost Per Acquisition).

The tool aggregates metrics (impressions, clicks, spend, conversions) by campaign ID using a streaming approach, so memory usage is proportional to the number of distinct campaigns rather than the input file size.

## Setup

**Prerequisites:** Go 1.24+

```bash
# Clone the repository
git clone https://github.com/ad-performance-aggregator.git
cd ad-performance-aggregator

# Build the binary
go build -o csvagg ./cmd/csvagg
```

### Docker

```bash
docker build -t csvagg .
```

## Usage

```bash
csvagg --input <csv_path> --output <output_dir> [--topk <number>] [--benchmark]
```

| Flag          | Type   | Default | Description                                    |
|---------------|--------|---------|------------------------------------------------|
| `--input`     | string | *required* | Path to input CSV file                      |
| `--output`    | string | *required* | Directory for output reports                |
| `--topk`      | int    | 10      | Number of top campaigns per report             |
| `--benchmark` | bool   | false   | Enable debug-level timing logs on stderr       |

### Example

```bash
./csvagg --input data.csv --output ./results --topk 10 --benchmark
```

Docker:

```bash
docker run --rm -v "$PWD":/data csvagg --input /data/data.csv --output /data/results
```

### Input format

CSV with the following required columns (order-independent):

```
campaign_id,impressions,clicks,spend,conversions
```

Rows with the same `campaign_id` are summed together.

### Output

Two CSV reports are written to the output directory:

- **`top{K}_ctr.csv`** -- Top K campaigns ranked by CTR (clicks / impressions), descending.
- **`top{K}_cpa.csv`** -- Top K campaigns ranked by CPA (spend / conversions), ascending. Campaigns with zero conversions are excluded.

## Running tests

```bash
go test ./...
```

## Libraries used

This project uses **only the Go standard library** -- no external dependencies.

| Package          | Purpose                            |
|------------------|------------------------------------|
| `encoding/csv`   | Streaming CSV parsing              |
| `flag`           | Command-line argument parsing      |
| `fmt`            | Formatted I/O                      |
| `io`             | I/O interfaces                     |
| `log/slog`       | Structured logging                 |
| `os`             | File system operations             |
| `path/filepath`  | File path manipulation             |
| `sort`           | Sorting for top-K ranking          |
| `strconv`        | String-to-number conversion        |
| `strings`        | String utilities                   |
| `time`           | Elapsed time measurement           |

## Performance

Benchmarked on a 1.1 GB CSV file (~30 million rows, 10,000 distinct campaigns):

| Metric             | Value     |
|--------------------|-----------|
| Processing time    | ~15s      |
| Peak memory (RSS)  | ~24 MB    |

Memory usage stays low because the tool streams CSV rows one at a time (`csv.Reader.ReuseRecord = true`) and only stores per-campaign aggregates in memory. The memory footprint scales with the number of distinct campaign IDs, not the size of the input file.

Use `--benchmark` to see per-phase timing on stderr:

```
processing data.csv ...
time=... level=DEBUG msg="parsed csv input" rows=30149188
time=... level=DEBUG msg="processing phase complete" elapsed=14.719s
time=... level=DEBUG msg="report writing phase complete" elapsed=5.644ms
done in 14.725s
reports written to ./results/
```
