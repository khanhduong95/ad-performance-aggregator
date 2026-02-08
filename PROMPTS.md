# Prompts Used During Development

Prompts used to build this project with AI assistance, showing how the
problem was broken down iteratively.

---

## 1. Build Streaming CSV Processor CLI

> Scaffold a production-quality Go CLI application that processes a very large
> CSV file (~1GB) in a streaming, memory-efficient way.
>
> Requirements:
> - CLI flags: `--input <csv_path>` and `--output <output_dir>`
> - Stream CSV line-by-line (do NOT load the whole file into memory)
> - Aggregate metrics by `campaign_id`:
>   total_impressions, total_clicks, total_spend, total_conversions
> - Compute CTR and CPA only after aggregation
> - Produce `top10_ctr.csv` (highest CTR) and `top10_cpa.csv` (lowest CPA, exclude zero conversions)
>
> Constraints: Readable idiomatic Go, explicit error handling, minimal memory
> usage, no unnecessary abstractions or dependencies.
>
> Output: Suggested folder structure, `main.go` with CLI wiring, core
> aggregation structs and functions, TODOs for validation/sorting/tests.

## 2. Add Interfaces and Dependency Injection

> Add minimal interfaces, DI, and unit tests.

## 3. Make Top-K Parameter Configurable via CLI

> top-K (default to 10) should be configurable via CLI optional param.

## 4. Introduce MetricsStore Interface

> Introduce a `MetricsStore` interface that replaces the implicit
> `map[string]*CampaignMetrics` coupling between Processor, Service,
> and ReportWriter.
>
> The store should own accumulation (write path) and top-K retrieval
> (read path). Provide an in-memory implementation as the only concrete
> backend. Refactor Processor to write into the store, ReportWriter to
> read from it, and Service to wire them together.
>
> This is a zero-behavior-change refactor. Update all existing tests
> and add `store_test.go`. Do not add extra backends, CLI flags, or
> context parameters. Run `go test ./...` before committing.

Iterated to trim the interface surface and rename implementations.

## 5. Create Dockerfile

> Create Dockerfile for application.

## 6. Add Benchmark Logs

> Add minimal benchmark logs at key execution points to demonstrate
> performance awareness, without impacting behavior or readability.

Iterated to make benchmarks opt-in via CLI flag and avoid global state.

## 7. Create README

> Generate a README.md including: setup instructions, how to run the program,
> libraries used, processing time for the 1GB file, peak memory usage
> (if measured).

Iterated to learn how to measure peak RSS and trim verbose package listings.

## 8. Review Code and Refactor

> Check these points for possible refactoring (do not blindly change just
> because I suggest them):
> - Add username to package namespace to match repository full name
> - Format multi-param functions/structs into multi lines
> - "Service" struct name seems vague
> - Most comments in code files are not needed

Iterated to align `go.mod` with the oldest supported Go version and
Dockerfile with the latest stable.
