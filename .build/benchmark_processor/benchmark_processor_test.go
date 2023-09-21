package main

import "testing"
import "os"
import "fmt"

func TestDecodeBenchmark(t *testing.T) {
	line := "BenchmarkRingBuffer/Queue-Dequeue-1000-ThreadSafe-NA-NA-12                 2123             2.945 ns/op               0 B/op          0 allocs/op"

	r, err := decodeBenchmark(line)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if r == nil {
		fmt.Fprint(os.Stderr, "Line was not parsed")
	} else {

		md := r.toMarkdown(10)
		_ = md
	}
}

func TestFormatFloat(t *testing.T) {
	floats := []float64{2.423, 24.23, 26945, 280900, 2792286}

	maxWidth := maxWidthFloat(floats)
	for _, f := range floats {
		s := formatFloat(f, maxWidth)
		l := len(s)
		fmt.Println(s)
		_ = l
	}
}
