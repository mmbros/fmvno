package util

// NewIntSequence return a new seguence
func NewIntSequence() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
