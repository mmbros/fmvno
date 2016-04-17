package util

import "time"

// RandInter is implemented by any value that has a RandInt method
type RandInter interface {
	RandInt() int
}

// RandFloat64er is implemented by any value that has a RandFloat64 method
type RandFloat64er interface {
	RandFloat64() float64
}

// RandIntByDateFunc type
type RandIntByDateFunc func(time.Time) int
