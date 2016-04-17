package distribution

import (
	"fmt"
	"math/rand"
)

// Normal is ..
type Normal struct {
	mean   float64
	stddev float64
}

// NewNormal returns a Normal aleatory variable
// with given mean and standard deviation
func NewNormal(mean, stddev float64) *Normal {
	return &Normal{
		mean:   mean,
		stddev: stddev,
	}
}

// Mean returns the mean value
func (d *Normal) Mean() float64 {
	return d.mean
}

// StdDev returns the standard deviation value
func (d *Normal) StdDev() float64 {
	return d.stddev
}

// RandFloat64 returns a normally distributed float64 in the range
// [-math.MaxFloat64, +math.MaxFloat64] with
// normal distribution with given mean and stddev
// from the default Source.
func (d *Normal) RandFloat64() float64 {
	return rand.NormFloat64()*d.stddev + d.mean
}

// RandInt returns int(RandFloat64())
func (d *Normal) RandInt() int {
	x := d.RandFloat64()
	return int(x)
}

func (d *Normal) String() string {
	return fmt.Sprintf("distribution.Normal{Mean:%v\n, StdDev:%f}", d.mean, d.stddev)
}

/*

// PrintTest genera `iter` numeri casuali richiamando RandInt
// e stampa la distribuzione ottenuta (assoluta e relativa)
func (d *Normal) PrintTest(iter int) {

	res := make(map[int]int)

	vmin := int(d.mean)
	vmax := vmin

	for j := 0; j < iter; j++ {
		v := d.RandInt()
		count, ok := res[v]
		if ok {
			res[v] = count + 1
		} else {
			res[v] = 1
		}

		if v < vmin {
			vmin = v
		}
		if v > vmax {
			vmax = v
		}
	}

	fmt.Printf("--- total = %d ---\n", iter)
	for v := vmin; v <= vmax; v++ {
		count, _ := res[v]
		fmt.Printf("freq[%v]: %d\n", v, count)
	}
}


*/
