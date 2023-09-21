package util

import "sync"

type CollectionBase struct {
	version int
	lock    sync.RWMutex
}
