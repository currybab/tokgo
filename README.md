# 🚀 TokGo - Go Tokenizer

[![Go Report Card](https://goreportcard.com/badge/github.com/currybab/tokgo)](https://goreportcard.com/report/github.com/currybab/tokgo)
<!-- [![GoDoc](https://godoc.org/github.com/currybab/tokgo?status.svg)](https://godoc.org/github.com/currybab/tokgo) -->
[![License: MIT](https://img.shields.io/github/license/currybab/tokgo)](https://opensource.org/license/mit/)
<!-- ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/currybab/tokgo/build-publish.yml)
<!-- Consider adding a CI badge once set up -->

Welcome to TokGo, a Go tokenizer library designed for use with OpenAI models inspired by [tiktoken](https://github.com/openai/tiktoken) and [jtokkit](https://github.com/knuddelsgmbh/jtokkit).

## ⚡ Performance

`tokgo` is designed to be highly performant, with a particular focus on memory efficiency. Here's a comparison with another popular Go tokenizer library, [`tiktoken-go`](https://github.com/pkoukk/tiktoken-go), based on benchmarks run on an Apple M1 Pro.

The benchmark measures the performance of encoding a large text file ([udhr concatted txt](https://research.ics.aalto.fi/cog/data/udhr/)) line by line.

| Library         | ns/op (lower is better) | B/op (lower is better) | allocs/op (lower is better) |
| --------------- | ----------------------- | ---------------------- | --------------------------- |
| **tokgo**       | 91,650                  | **33,782**             | **445**                     |
| `tiktoken-go`   | 91,211                  | 45,511                 | 564                         |

### Key Advantages:

*   **⚡ Comparable Speed:** `tokgo` offers processing speeds nearly identical to `tiktoken-go`.
*   **🧠 Superior Memory Efficiency:** `tokgo` uses approximately **26% less memory** per operation.
*   **🗑️ Fewer Allocations:** It also results in about **21% fewer memory allocations**, reducing garbage collection overhead.

This makes `tokgo` an excellent choice for memory-sensitive applications or high-throughput systems where minimizing GC pressure is critical.

## 📚 Usage

```go
package main

import (
	"fmt"

	tokmod "github.com/currybab/tokgo/mod"
	tokgo "github.com/currybab/tokgo/registry"
)

func main() {
	enc, _ := tokgo.NewDefaultEncodingRegistry().GetEncodingByType(tokmod.CL100K_BASE)

	// this should print a list of token ids
	ids := enc.EncodeToIntArray("The quick brown fox jumps over the lazy dog 😎")
	fmt.Println(ids)

	// this should print the original string back
	text := enc.Decode(ids)
	fmt.Println(text)
}
```
