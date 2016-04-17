package model

import "github.com/mmbros/fmvno/util"

// Comsy is ...
type Comsy struct {
	ID   int
	Name string
}

var nextComsyID = util.NewIntSequence()

// NewComsy restituisce una nuova consegna
func NewComsy(name string) (*Comsy, error) {
	c := &Comsy{
		ID:   nextComsyID(),
		Name: name,
	}
	return c, nil
}
