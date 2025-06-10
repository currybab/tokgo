package main

import (
	"log"
	"os"
	"strings"
	"testing"

	tokmod "github.com/currybab/tokgo/mod"
	tokgo "github.com/currybab/tokgo/registry"
)

// go test -benchmem -run=^$ -bench ^BenchmarkEncodingInFullLanguage$ -benchtime=100000x github.com/currybab/tokgo/benchmark
func BenchmarkEncodingInFullLanguage(b *testing.B) {
	data, err := os.ReadFile("../tmp/udhr.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")
	encoder, err := tokgo.NewDefaultEncodingRegistry().GetEncodingForModelType(tokmod.GPT_4O)
	lineCount := len(lines)
	if err != nil {
		log.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		encoder.EncodeOrdinaryToIntArray(lines[n%lineCount])
	}
}
