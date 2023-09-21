package util

import (
	"math/rand"
	"sort"
	"time"
)

/*
Functions for generating test data
*/

func CreateRandInt(seed *int64) int {
	rand.Seed(*seed)
	*seed++
	return rand.Int()
}

func CreateIntListData(arraySize int, seed *int64) (headItems, tailItems, headItemsReverse, tailItemsReverse []int) {

	headItems = make([]int, arraySize)
	tailItems = make([]int, arraySize)
	headItemsReverse = make([]int, arraySize)
	tailItemsReverse = make([]int, arraySize)

	for i := 0; i < arraySize; i++ {
		index := (arraySize - 1) - i
		head := CreateRandInt(seed)
		tail := CreateRandInt(seed)
		headItems[i] = head
		headItemsReverse[index] = head
		tailItems[i] = tail
		tailItemsReverse[index] = tail
	}

	return
}

func CreateSingleIntListData(arraySize int, seed *int64) []int {

	items := make([]int, arraySize)

	for i := 0; i < arraySize; i++ {
		items[i] = CreateRandInt(seed)
	}

	return items
}

func GetTestMinMax(data []int) (int, int) {
	sorted := make([]int, len(data))
	copy(sorted, data)
	sort.Ints(sorted)
	return sorted[0], sorted[len(data)-1]
}

func CreateMinMaxTestData(arraySize int, seed *int64) (data []int, min int, max int) {
	data = CreateSingleIntListData(arraySize, seed)
	min, max = GetTestMinMax(data)
	return
}

func CreateSerialMinMaxTestData(arraySize int, seed *int64) (data []int, min int, max int) {
	data = CreateSerialIntListData(arraySize, seed)
	min = data[0]
	max = data[arraySize-1]
	return
}

func CreateSmallIntListData(arraySize int, seed *int64) (headItems, tailItems, headItemsReverse, tailItemsReverse []int) {

	headItems = make([]int, arraySize)
	tailItems = make([]int, arraySize)
	headItemsReverse = make([]int, arraySize)
	tailItemsReverse = make([]int, arraySize)

	for i := 0; i < arraySize; i++ {
		index := (arraySize - 1) - i
		head := CreateRandInt(seed) & 0xffff
		tail := CreateRandInt(seed) & 0xffff
		headItems[i] = head
		headItemsReverse[index] = head
		tailItems[i] = tail
		tailItemsReverse[index] = tail
	}

	return
}

func CreateSerialIntListData(arraySize int, seed *int64) (headItems []int) {
	headItems = make([]int, arraySize)

	for i := 0; i < arraySize; i++ {
		headItems[i] = int(*seed)
		*seed++
	}

	return
}

func CreateSerialSmallIntListData(arraySize int, seed *int64) (headItems, tailItems, headItemsReverse, tailItemsReverse []int) {
	headItems = make([]int, arraySize)
	tailItems = make([]int, arraySize)
	headItemsReverse = make([]int, arraySize)
	tailItemsReverse = make([]int, arraySize)

	for i := 0; i < arraySize; i++ {
		index := (arraySize - 1) - i
		head := *seed % 0xffff
		tail := *seed % 0xffff
		headItems[i] = int(head)
		headItemsReverse[index] = int(head)
		tailItems[i] = int(tail)
		tailItemsReverse[index] = int(tail)
		*seed++
	}

	return
}

func CreateTimeListData(arraySize int, seed *int64) []time.Time {
	var dayFactor int
	items := make([]time.Time, arraySize)
	year := 2023
	loc, _ := time.LoadLocation("Europe/London")
	for i := 0; i < arraySize; i++ {
		val1 := CreateRandInt(seed)
		val2 := CreateRandInt(seed)
		ns := (val1 & 0x3FFFFFFF) % 1000000000
		sec := (val2 & 0x3f) % 60
		min := ((val2 >> 6) & 0x3f) % 60
		hr := ((val2 >> 12) & 0x1f) % 24
		mon := (((val2 >> 17) & 0x1f) % 12) + 1

		switch mon {
		case 1, 3, 5, 7, 8, 10, 12:
			dayFactor = 31
		case 4, 6, 9, 11:
			dayFactor = 30
		default:
			dayFactor = 28
		}

		day := (((val2 >> 22) & 0x1f) % dayFactor) + 1

		items[i] = time.Date(year, time.Month(mon), day, hr, min, sec, ns, loc)
	}

	return items
}
