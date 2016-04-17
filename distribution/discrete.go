package distribution

import (
	"fmt"
	"math/rand"
	"strconv"
)

// Discrete represents a discrete probability distribution
// from 0 to L-1, with L = Len(weights)
type Discrete struct {
	Weights   *[]float64
	Total     float64
	Freq      *[]float64
	FreqCumul *[]float64
}

func arrayStringToFloat64(in []string) ([]float64, error) {
	var err error
	out := make([]float64, len(in))

	for j, s := range in {
		if out[j], err = strconv.ParseFloat(s, 64); err != nil {
			return nil, err
		}

	}
	return out, nil
}

// NewDiscreteFromStrings return a new Discrete distribution with given weights (as string)
func NewDiscreteFromStrings(weights []string) (*Discrete, error) {
	arrFloat, err := arrayStringToFloat64(weights)
	if err != nil {
		return nil, err
	}
	return NewDiscrete(arrFloat...)
}

// NewDiscrete return a new Discrete distribution with given weights
func NewDiscrete(weights ...float64) (*Discrete, error) {

	total := 0.0
	for _, v := range weights {
		total += v
	}

	L := len(weights)
	freq := make([]float64, L)
	freqcumul := make([]float64, L)

	for j, v := range weights {
		freq[j] = v / total
	}

	prec := 0.0
	for j, v := range weights {
		freq[j] = v / total
		prec += freq[j]
		freqcumul[j] = prec
	}

	dd := &Discrete{
		Weights:   &weights,
		Total:     total,
		Freq:      &freq,
		FreqCumul: &freqcumul,
	}

	return dd, nil
}

// Len returns then number of weights of the distribution
func (d *Discrete) Len() int {
	return len(*d.Weights)
}

// RandInt returns, as a Int, a pseudo-random number in [0 .. Len()-1]
// The random number has the given probability distribution
func (d *Discrete) RandInt() int {
	x := rand.Float64()

	for j, v := range *d.FreqCumul {
		if v >= x {
			return j
		}
	}
	return len(*d.FreqCumul) - 1
}

func (d *Discrete) String() string {
	return fmt.Sprintf("distribution.Discrete{\n  Weights:  %v\n  Total: %f\n  Freq:  %v\n  Cumul: %v\n}\n", *d.Weights, d.Total, *d.Freq, *d.FreqCumul)
}

/*

// PrintTest genera `iter` numeri casuali richiamando RandInt
// e stampa la distribuzione ottenuta (assoluta e relativa)
func (d *Discrete) PrintTest(iter int) {

	res := make(map[int]int)

	for j := 0; j < iter; j++ {
		v := d.RandInt()
		count, ok := res[v]
		if ok {
			res[v] = count + 1
		} else {
			res[v] = 1
		}
	}

	fmt.Printf("--- total = %d ---\n", iter)
	w := 1.0 / float64(iter)
	for k := 0; k < d.Len(); k++ {
		fmt.Printf("freq[%v]: %0.4f - count[%v]: %v\n", k, float64(res[k])*w, k, res[k])
	}
}

*/
