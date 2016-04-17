package distribution

import (
	"fmt"
	"math/rand"
	"sort"
)

// Geometric is ..
type Geometric struct {
	mean float64
	p    float64
	q    float64
}

// NewGeometric is ...
func NewGeometric(mean float64) *Geometric {
	if mean <= 0.0 {
		panic("NewGeometric: mean must be positive")
	}
	p := 1.0 / mean
	return &Geometric{
		mean: mean,
		p:    p,
		q:    1 - p,
	}
}

func (g *Geometric) String() string {
	return fmt.Sprintf("distribution.Geometric{Mean: %0.4f}", g.mean)
}

// Mean returns the mean value
func (g *Geometric) Mean() float64 {
	return g.mean
}

// P returns the p value
func (g *Geometric) P() float64 {
	return g.p
}

// Q returns the q = 1-p value
func (g *Geometric) Q() float64 {
	return g.q
}

// RandInt is ...
func (g *Geometric) RandInt() int {
	// 1-x
	umx := rand.Float64()
	q := g.q
	qk := q

	// 100 is an upper bound
	for k := 1; k < 100; k++ {
		// q^k < 1-x
		if qk < umx {
			return k
		}
		qk *= q
	}

	return 1

}

// PrintTest genera `iter` numeri casuali richiamando RandInt
// e stampa la distribuzione ottenuta (assoluta e relativa)
func (g *Geometric) PrintTest(iter int) {

	res := make(map[int]int)

	for j := 0; j < iter; j++ {
		v := g.RandInt()
		count, ok := res[v]
		if ok {
			res[v] = count + 1
		} else {
			res[v] = 1
		}
	}

	L := len(res)
	v := make([]int, L)
	j := 0
	for k := range res {
		v[j] = k
		j++
	}
	sort.Ints(v)

	fmt.Printf("--- %v ---\n", g)
	w := 1.0 / float64(iter)
	sum := 0
	for j := 0; j < L; j++ {
		k := v[j]
		sum += k * res[k]
		fmt.Printf("freq[%2v]: %0.4f - count[%2v]: %v\n", k, float64(res[k])*w, k, res[k])
	}
	fmt.Printf("%8s: %0.4f - %9s: %d\n", "MEDIA", float64(sum)/float64(iter), "TOT", iter)
}
