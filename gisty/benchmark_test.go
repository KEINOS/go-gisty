package gisty

import (
	"os"
	"testing"
)

// Before enhancement it is recommended to use the following command to run the benchmark:
//
// To benchmark, run:
//   go test -bench . ./gisty/... -count 10 -benchmem > before.txt
//
// To compare after enhancement, run:
//   benchstat before.txt after.txt

func Benchmark_sanitizeGistID(b *testing.B) {
	// The Complete Works of William Shakespeare by William Shakespeare
	// http://www.gutenberg.org/files/100/100-0.txt
	data, err := os.ReadFile("testdata/shakespare.txt")
	if err != nil {
		b.Fatal(err)
	}

	strShakespeare := string(data)

	b.ResetTimer()

	for N := 0; N < b.N; N++ {
		sanitizeGistID(strShakespeare)
	}
}
