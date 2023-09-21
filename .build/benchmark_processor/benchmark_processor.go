/*
Render benchmark output from stdin to markdown document
*/
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var bmRx = regexp.MustCompile(`Benchmark(?P<collection>\w+)/(?P<type>\w+)-(?P<operation>\w+)-(?P<elements>\d+)-(?P<threadsafe>\w+)-(?P<presize>\w+)-(?P<concurrency>\w+)-\d+\s+\d+\s+(?P<nanoseconds>[\d]+(\.\d+)?)\sns/op\s+(?P<bytes>\d+)\sB/op\s+(?P<alloc>\d+)\s+allocs/op`)

type presize int

const (
	NotApplicableP presize = iota
	NoPresize
	Presize
)

type threadSafety int

const (
	NotApplicableT threadSafety = iota
	NoThreadSafe
	ThreadSafe
)

type concurrency int

const (
	NotApplicableC concurrency = iota
	NoConcurrent
	Concurrent
)

type Result struct {
	collection  string
	typ         string
	op          string
	elements    int
	threadsafe  threadSafety
	presize     presize
	concurrent  concurrency
	nanoseconds float64
	bytes       int
	allocs      int
}

var indent = "    "

func (*Result) toMarkdownHeader() string {
	h1 := indent + "| Collection | Operation | Elements | TS | PS | CC | ns/op | B/op | allocs/op |"
	h2 := indent + "|------------|-----------|---------:|:--:|:--:|:--:|------:|-----:|----------:|"
	return h1 + "\n" + h2
}

func (r *Result) toMarkdown(maxWidthNs int) string {
	return fmt.Sprintf("%s| %s | %s | %s | %s | %s | %s | %s | %s | %s |",
		indent,
		r.collection,
		//r.typ,
		r.op,
		formatInt(r.elements),
		Iif(r.threadsafe == NotApplicableT, "", Iif(r.threadsafe == NoThreadSafe, ":x:", ":heavy_check_mark:")),
		Iif(r.presize == NotApplicableP, "", Iif(r.presize == NoPresize, ":x:", ":heavy_check_mark:")),
		Iif(r.concurrent == NotApplicableC, "", Iif(r.concurrent == NoConcurrent, ":x:", ":heavy_check_mark:")),
		formatFloat(r.nanoseconds, maxWidthNs),
		formatInt(r.bytes),
		formatInt(r.allocs))
}

func renderMarkdown(rs []*Result, maxWidthNs int) {
	var currentType, nextType string

	for i, r := range rs {
		nextType = r.typ

		if currentType != nextType {
			if i > 0 {
				fmt.Printf("\n%s</details>\n", indent)
				fmt.Println()
			}

			fmt.Printf("* %ss\n\n", r.typ)
			fmt.Println(indent + "<details>")
			fmt.Println(indent + "<summary>Expand</summary>\n")
			fmt.Println(r.toMarkdownHeader())
			currentType = nextType
		}

		fmt.Println(r.toMarkdown(maxWidthNs))
	}

	fmt.Printf("\n%s</details>\n", indent)
}

func formatInt(i any) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", i)
}

func formatFloat(f float64, maxWidth int) string {
	var str string
	p := message.NewPrinter(language.English)

	switch {
	case f < 10:
		str = p.Sprintf("%.3f", f)
	case f < 100:
		str = p.Sprintf("%.2f ", f)
	case f < 1000:
		str = p.Sprintf("%.1f  ", f)
	default:
		str = p.Sprintf("%d    ", int64(f))
	}

	return fmt.Sprintf("`%*s`", maxWidth+4, str)
}

func maxWidth(results []*Result) (nsWidth int) {

	floats := []float64{}

	for _, v := range results {
		floats = append(floats, v.nanoseconds)
	}

	nsWidth = maxWidthFloat(floats)
	return
}

func maxWidthFloat(nums []float64) int {
	max := nums[0]

	for _, v := range nums {
		if v > max {
			max = v
		}
	}

	l := math.Floor(math.Log10(max))
	commas := math.Floor(l / 3)
	return int(l+commas) + 1
}

func mustAtoi(val string) int {
	res, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return res
}

func mustAtof(val string) float64 {
	res, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return res

}

func mustDecodeThreadSafe(val string) threadSafety {
	if val == "ThreadSafe" {
		return ThreadSafe
	}

	if val == "NoThreadSafe" {
		return NoThreadSafe
	}

	if val == "NA" {
		return NotApplicableT
	}

	panic(fmt.Sprintf("Unrecognized thread safety: %s", val))
}

func mustDecodePresize(val string) presize {
	if val == "Presize" {
		return Presize
	}

	if val == "NoPresize" {
		return NoPresize
	}

	if val == "NA" {
		return NotApplicableP
	}

	panic(fmt.Sprintf("Unrecognized presize: %s", val))
}

func mustDecodeConcurrency(val string) concurrency {
	if val == "Concurrent" {
		return Concurrent
	}

	if val == "NoConcurrent" {
		return NoConcurrent
	}

	if val == "NA" {
		return NotApplicableC
	}

	panic(fmt.Sprintf("Unrecognized presize: %s", val))
}

func decodeBenchmark(str string) (*Result, error) {
	match := bmRx.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		key := bmRx.SubexpNames()[i]
		if len(key) != 0 {
			// Only named groups
			results[key] = name
		}
	}

	if len(results) == 0 {
		return nil, nil
	}

	if len(results) != 10 {
		return nil, fmt.Errorf("failed to parse: %s. len = %d", str, len(results))
	}

	res := &Result{
		collection:  results["collection"],
		typ:         results["type"],
		op:          results["operation"],
		elements:    mustAtoi(results["elements"]),
		threadsafe:  mustDecodeThreadSafe(results["threadsafe"]),
		presize:     mustDecodePresize(results["presize"]),
		concurrent:  concurrency(mustDecodeConcurrency(results["concurrency"])),
		nanoseconds: mustAtof(results["nanoseconds"]),
		bytes:       mustAtoi(results["bytes"]),
		allocs:      mustAtoi(results["alloc"]),
	}

	return res, nil
}

func sortResults(results []*Result) {

	opWeights := map[string]int{
		"Push":     0,
		"Enqueue":  0,
		"Add":      0,
		"Pop":      1,
		"Dequeue":  1,
		"Remove":   1,
		"Sort":     2,
		"Min":      3,
		"Max":      4,
		"Contains": 5,
	}

	sort.Slice(results, func(i, j int) bool {
		// Class
		if results[i].typ != results[j].typ {
			return results[i].typ < results[j].typ
		}

		// then Collection
		if results[i].collection != results[j].collection {
			return results[i].collection < results[j].collection
		}

		// then Operation
		opWeightI, ok := opWeights[results[i].op]
		if !ok {
			opWeightI = 10
		}

		opWeightJ, ok := opWeights[results[j].op]
		if !ok {
			opWeightJ = 10
		}

		if opWeightI != opWeightJ {
			return opWeightI < opWeightJ
		}

		// then threadsafe
		if results[i].threadsafe != results[j].threadsafe {
			return results[i].threadsafe < results[j].threadsafe
		}

		// then presize
		if results[i].presize != results[j].presize {
			return results[i].presize < results[j].presize
		}

		// then size
		if results[i].elements != results[j].elements {
			return results[i].elements < results[j].elements
		}

		// then concurrent
		return results[i].concurrent < results[j].concurrent
	})
}

func Iif[T any](predicate bool, trueVal, falseVal T) T {
	if predicate {
		return trueVal
	}

	return falseVal
}

func main() {
	var results []*Result

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		fmt.Fprint(os.Stderr, text)
		result, err := decodeBenchmark(text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}

		if result != nil {
			results = append(results, result)
		}
	}

	sortResults(results)
	maxWidthNs := maxWidth(results)

	// fmt.Println(results[0].toMarkdownHeader())
	// for _, v := range results {
	// 	fmt.Println(v.toMarkdown(maxWidthNs))
	// }

	renderMarkdown(results, maxWidthNs)
}
