# Prompts Used During Development

This document records the prompts used to build this project with AI assistance.
It is intended to give reviewers insight into how the problem was broken down,
the iterative communication style, and the overall problem-solving approach.

---

## 1. Build Streaming CSV Processor CLI

> Scaffold a production-quality Go CLI application that processes a very large
> CSV file (~1GB) in a streaming, memory-efficient way.
>
> Requirements:
> - CLI flags: `--input <csv_path>` and `--output <output_dir>`
> - Stream CSV line-by-line (do NOT load the whole file into memory)
> - Aggregate metrics by `campaign_id`:
>   - total_impressions
>   - total_clicks
>   - total_spend
>   - total_conversions
> - Compute CTR and CPA only after aggregation
> - Produce:
>   - `top10_ctr.csv` (highest CTR)
>   - `top10_cpa.csv` (lowest CPA, exclude zero conversions)
>
> Constraints:
> - Readable, idiomatic Go
> - Explicit error handling
> - Minimal memory usage
> - No unnecessary abstractions or dependencies
>
> Output:
> - Suggested folder structure
> - `main.go` with CLI wiring
> - Core aggregation structs and functions
> - TODOs for validation, sorting, and tests

**Why this prompt:** Established the full scope upfront — data flow, output
format, and quality constraints — so the initial scaffold would be close to
production-ready rather than a throwaway prototype.

---

## 2. Add Interfaces and Dependency Injection

> Add minimal interfaces, DI, and unit tests.

**Why this prompt:** After the working scaffold existed, this was the
smallest step that would make the code testable in isolation without
over-engineering the design.

---

## 3. Make Top-K Parameter Configurable via CLI

> top-K (default to 10) should be configurable via CLI optional param.

**Why this prompt:** A single, focused change — keeps the diff reviewable and
avoids mixing feature work with refactoring.

---

## 4. Support Alternative Storage for Aggregated Metrics

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

**Why this prompt:** Explicitly scoped as a zero-behavior-change refactor to
keep the review focused on structure, not logic. Calling out what *not* to do
(extra backends, new flags) prevented scope creep.

---

## 5. Refine MetricsStore Interface

> **Prompt 1:** (same as #4 — the interface was introduced here after the
> earlier attempt needed rework)
>
> **Follow-up prompts:** Remove unnecessary functions, rename implementations.

**Why these prompts:** The first pass introduced more surface area than needed.
Follow-up prompts trimmed the interface to only what callers actually use,
keeping the API minimal.

---

## 6. Create Dockerfile

> Create Dockerfile for application.

**Why this prompt:** Short and direct — Docker best practices are well-known,
so a concise prompt was sufficient.

---

## 7. Add Benchmark Logs for Performance Tracking

> **Prompt 1:** Add minimal benchmark logs at key execution points to
> demonstrate performance awareness, without impacting behavior or readability.
>
> **Prompt 2:** Benchmarks should be enabled by CLI flag.
>
> **Follow-up prompts:** Prevent unnecessary global `Benchmark` bool.

**Why these prompts:** Started with the simplest thing (always log), then made
it opt-in via a flag when it became clear the output was noisy. The follow-up
caught a design smell (global mutable state) before it spread.

---

## 8. Create README with Setup and Usage Instructions

> **Prompt 1:** Generate a README.md including:
> - Setup instructions
> - How to run the program
> - Libraries used
> - Processing time for the 1GB file
> - Peak memory usage (if measured)
>
> **Prompt 2:** How can I get Peak memory (RSS) myself?
>
> **Follow-up prompts:** Remove unnecessary table of packages in
> "Libraries used" section.

**Why these prompts:** The first prompt produced a comprehensive draft. The
second was a learning question to fill a knowledge gap. Follow-ups trimmed
boilerplate — listing every stdlib package adds noise, not value.

---

## 9. Review Code and Refactor

> **Prompt 1:** Check these points for possible refactoring (do not blindly
> change just because I suggest them):
> - Add username to package namespace to match repository full name
> - Format multi-param functions/structs into multi lines
> - "Service" struct name seems vague
> - Most comments in code files are not needed
>
> **Follow-up prompts:** Update `go.mod` to use the oldest Go version
> supporting the current codebase; update Dockerfile to use the latest
> stable Go version.

**Why these prompts:** The parenthetical "do not blindly change" was
intentional — it forced evaluation of each suggestion on merit rather than
applying every change mechanically. The Go version follow-ups ensured
compatibility without requiring bleeding-edge toolchains.

---

## Problem-Solving Approach

1. **Start with a working vertical slice** — prompt 1 delivered an end-to-end
   pipeline before any refactoring.
2. **Refactor in isolated steps** — each subsequent prompt targeted one concern
   (testability, configurability, storage abstraction).
3. **Constrain scope explicitly** — phrases like "zero-behavior-change",
   "do not add extra backends", and "do not blindly change" kept changes
   focused.
4. **Iterate on output quality** — follow-up prompts trimmed unnecessary code,
   renamed vague identifiers, and removed noisy documentation.
5. **Ask learning questions inline** — prompt 8's "how can I get peak memory
   myself" turned a documentation task into a knowledge-building moment.
