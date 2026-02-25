# Go JSON to Excel/CSV Exporter (j2e) 

[![Go CI](https://github.com/hosseineddin/go-json2excel/actions/workflows/ci.yml/badge.svg)](https://github.com/hosseineddin/go-json2excel/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hosseineddin/go-json2excel)](https://goreportcard.com/report/github.com/hosseineddin/go-json2excel)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hosseineddin/go-json2excel)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


A highly optimized, enterprise-grade, zero-copy streaming library for converting massive dynamic JSON payloads into `XLSX` (Excel) or `CSV` formats. Built specifically for large-scale Go backend systems to directly pipe HTTP requests/responses without exhausting server RAM.

## Under the Hood (Architecture)

Unlike traditional libraries that load the entire JSON array into memory (causing Out-Of-Memory/OOM kills), `go-json2excel` operates on a strict **Streaming Pipeline**:
1. **Zero-Copy concept**: Data is read from an `io.Reader` (e.g., a Database cursor or HTTP request) and written directly to an `io.Writer` (e.g., HTTP Response).
2. **Aggressive Sync.Pool**: Slice allocations are reused via `sync.Pool`. This drops Garbage Collection (GC) pauses to near zero, ensuring a flat $O(1)$ memory footprint relative to the number of rows.
3. **Dynamic Type Strategy**: Uses the Strategy Pattern to infer JSON types (including RFC3339 timestamps) on the fly without heavy reflection overhead.
4. **Auto-Pagination**: Excel has a hard limit of ~1M rows. The engine automatically creates new sheets to prevent file corruption.

## Key Features
- **$O(1)$ Memory Complexity**: Process 10,000 or 10,000,000 records with the exact same RAM usage (~70MB peak).
- **Format Agnostic**: Seamlessly switch between `.xlsx` and `.csv` using the Factory pattern.
- **Concurrent Decoding**: JSON decoding and Excel encoding run on separate Goroutines communicating via buffered channels.
- **Production Ready**: 100% test coverage with built-in race-condition checks.

## Installation
```bash
go get [github.com/hosseineddin/go-json2excel](https://github.com/hosseineddin/go-json2excel)