package model

import (
	"fmt"
	"time"

	"github.com/mmbros/fmvno/util"
)

var nextConsegnaID = util.NewIntSequence()

// Consegna is ...
type Consegna struct {
	ID           int
	MobileNumber *Mobile
	Comsy        *Comsy
	DataConsegna time.Time
	UsimFMVNO    *Usim
}

// NewConsegna restituisce una nuova consegna
// Sono fatti dei controlli di congruenza che possono portare ad errori
// In caso di esito positivo, il mobile viene posto in stato MigrConsegnato
func NewConsegna(mobile *Mobile, comsy *Comsy, date time.Time) (*Consegna, error) {

	// contolli di congruenza
	if comsy == nil {
		return nil, fmt.Errorf("comsy must be NOT nil")
	}
	if mobile == nil {
		return nil, fmt.Errorf("mobile must be NOT nil")
	}

	if mobile.MigrStatus != MigrDaConsegnare {
		return nil, fmt.Errorf("Invalid mobile migration status")
	}

	// crea una nuova consegna
	c := &Consegna{
		ID:           nextConsegnaID(),
		MobileNumber: mobile,
		Comsy:        comsy,
		DataConsegna: date,
		UsimFMVNO:    NewUsim(),
	}

	// aggiuna usim e migration status del numero
	mobile.UsimFMVNO = c.UsimFMVNO
	mobile.MigrStatus = MigrConsegnatoDealer

	return c, nil
}
