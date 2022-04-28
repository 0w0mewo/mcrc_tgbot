package utils

import (
	"sync"
)

type RandomMap struct {
	randmap    *sync.Map
	start, end int
}

func NewRandomMap(start, end int) *RandomMap {
	return &RandomMap{
		start:   start,
		end:     end,
		randmap: &sync.Map{},
	}
}

func (rm *RandomMap) Get(key any) int {
	if rn, exist := rm.randmap.Load(key); exist {
		return rn.(int)
	}

	return rm.Generate(key)
}

func (rm *RandomMap) Generate(key any) int {
	rn := RandomBetween(rm.start, rm.end)
	rm.randmap.Store(key, rn)

	return rn
}

func (rm *RandomMap) SetStart(start int) {
	rm.start = start
}

func (rm *RandomMap) SetEnd(end int) {
	rm.end = end
}

func (rm *RandomMap) Low() int {

	return rm.start
}

func (rm *RandomMap) High() int {
	return rm.end
}
