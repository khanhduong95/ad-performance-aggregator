# ad-performance-aggregator

A streaming CLI tool that processes large CSV files of advertising campaign performance data and produces ranked reports by CTR (Click-Through Rate) and CPA (Cost Per Acquisition).

The tool aggregates metrics (impressions, clicks, spend, conversions) by campaign ID using a streaming approach, so memory usage is proportional to the number of distinct campaigns rather than the input file size.

## Setup

**Prerequisites:** Go 1.21+

```bash
# Clone the repository
git clone https://github.com/khanhduong95/ad-performance-aggregator.git
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
./csvagg --input ad_data.csv --output ./results --topk 10 --benchmark
```

Docker:

```bash
docker run --rm -v "$PWD":/data csvagg --input /data/ad_data.csv --output /data/results
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

Only the Go standard library -- no external dependencies.

## Performance

Benchmarked on the given 1 GB CSV file:

| Metric             | Value     |
|--------------------|-----------|
| Processing time    | ~5s      |
| Peak memory (RSS)  | ~8 MB    |
