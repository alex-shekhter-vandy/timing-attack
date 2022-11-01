package main

import (
	"time"
)

func Max(ress []result) (found result) {
	for _, v := range ress {
		if found.duration < v.duration {
			found = v
		}
	}
	// log.Printf("Max ress: %+v; found: %+v", ress, found)
	return found
}

func Min(ress []result) (found result) {
	if len(ress) == 0 {
		return found
	}
	found = ress[0]
	found.duration = 999999999999999999
	for _, v := range ress {
		if found.duration > v.duration {
			found = v
		}
	}
	return found
}

func Avg(ds []time.Duration) (d time.Duration) {
	lds := int64(len(ds))
	if lds == 0 {
		return d
	}
	for _, v := range ds {
		d += v
	}
	return time.Duration(int64(d) / lds)
}
