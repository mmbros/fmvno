package model

import (
	"fmt"
	"time"

	"github.com/mmbros/fmvno/util"
)

var nextSpedizioneID = util.NewIntSequence()

// SpedizioneStatus ...
type SpedizioneStatus int

// Spedizione is ...
type Spedizione struct {
	ID               int
	Status           SpedizioneStatus
	Esito            SpedizioneEsito
	EsitoForzato     bool
	EsitoDescr       string
	DataCreazione    time.Time
	DataInvio        time.Time
	DataSpedizione   time.Time
	DataEsito        time.Time
	MobileNumbers    MobileList
	UsimNumbers      UsimList
	ElapsedGiorniLav int
}

// SpedizioneList ...
type SpedizioneList []*Spedizione

// SpedizioneMap ...
type SpedizioneMap map[int]*Spedizione

// NewSpedizione restituisce una nuova spedizione
// associata ai mobili della lista
// NESSUNA MODIFICA VIENE FATTA SUI MOBILI
// Sono fatti solo controlli di congruenza che possono portare ad errori di validazione
// Nel caso di errore di validazione, la spedizione ha:
// Status = errore, Esito = errore di validazione, EsitoDescr = descrizione errore
func NewSpedizione(mobiles MobileList) *Spedizione {

	sped := &Spedizione{
		ID:            nextSpedizioneID(),
		DataCreazione: util.SimulDate(),
		MobileNumbers: mobiles,
	}

	// funzione ausiliaria per restituire una spedizione con errore di validazione
	fnErrValidazione := func(descrErr string) *Spedizione {
		sped.Status = SpedizioneErrore
		sped.Esito = ErrValidazione
		sped.EsitoDescr = descrErr
		sped.DataEsito = util.SimulDate()
		return sped
	}

	// check mobiles list
	if len(mobiles) == 0 {
		return fnErrValidazione("mobile list must not be empty")
	}

	// check mobiles attributes
	for _, mob := range mobiles {
		if mob.MigrStatus != MigrDaSpedireOPL {
			msg := fmt.Sprintf("Invalid mobile status: expected MigrDaSpedireOPL (%d) found (%d)", MigrDaSpedireOPL, mob.MigrStatus)
			return fnErrValidazione(msg)
		}
		if mob.UsimFMVNO != nil {
			return fnErrValidazione("UsimFMVNO must be nil")
		}
	}

	return sped
}

// SpedizioneStatus values
const (
	SpedizioneDaInviare  SpedizioneStatus = iota // --> Esito = EsitoNonDisponibile
	SpedizioneInCorso                            // --> Esito = EsitoNonDisponibile
	SpedizioneConsegnata                         // --> Esito = ConsegnaOK
	SpedizioneErrore                             // --> Esito = ErrConsegnaXXX
)

// SpedizioneEsito ...
type SpedizioneEsito int

// SpedizioneEsito values
const (
	EsitoNonDisponibile            SpedizioneEsito = -1 + iota
	ConsegnaOK                                     // 0
	ErrConsegnaInGiacenza                          // 1
	ErrConsegnaIndirizzoNonTrovato                 // 2
	ErrConsegnaGenerico                            // 3
	ErrValidazione                                 // 4
)
