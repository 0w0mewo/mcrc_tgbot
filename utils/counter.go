package utils

import (
	"strconv"
	"strings"
	"sync"
)

type Counter struct {
	counters map[string]int
	lock     *sync.RWMutex
}

func NewCounter() *Counter {
	return &Counter{
		counters: make(map[string]int),
		lock:     &sync.RWMutex{},
	}
}

func (mc *Counter) Inc(id string) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mc.counters[id]++
}

func (mc *Counter) Get(id string) int {
	mc.lock.RLock()
	defer mc.lock.RUnlock()

	return mc.counters[id]
}

func (mc *Counter) Reset(id string) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mc.counters[id] = 0
}

func (mc *Counter) String() string {

	strBuilder := &strings.Builder{}
	strBuilder.WriteString("counter status: ")

	mc.lock.RLock()
	for k, v := range mc.counters {
		strBuilder.WriteString(k + ", " + strconv.Itoa(v) + "\n")
	}
	mc.lock.RUnlock()

	return strBuilder.String()
}
