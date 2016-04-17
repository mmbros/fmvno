package oplold

import (
	"fmt"
	"time"

	"github.com/mmbros/fmvno/model"
	"github.com/mmbros/fmvno/util"
)

// SpedizioneStatus ...
type SpedizioneStatus int

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

// Spedizione ...
type Spedizione struct {
	ID               int              `json:"spedizione_id"`
	Status           SpedizioneStatus `json:"stato"`
	Esito            SpedizioneEsito  `json:"esito"`
	EsitoForzato     bool             `json:"esito_forzato,omitempty"`
	EsitoDescr       string           `json:"esito_descr,omitempty"`
	DataCreazione    time.Time        `json:"data_creaz"`
	DataSpedizione   time.Time        `json:"data_spediz"`
	DataEsito        time.Time        `json:"data_esito"`
	RequestID        int              `json:"request_id"`
	MobileNumbers    model.MobileList `json:"mobile_numbers"`
	UsimNumbers      model.UsimList   `json:"usims"`
	ElapsedGiorniLav int              `json:"elapsed_gg_lav"`
}

// SpedizioneList ...
type SpedizioneList []*Spedizione

// SpedizioneMap ...
type SpedizioneMap map[int]*Spedizione

var nextSpedizioneID = util.NewIntSequence()

// newSpedizione restituisce una nuova spedizione
// associata alla requestID e ai mobili della lista
// NESSUNA MODIFICA VIENE FATTA SUI MOBILI
// Sono fatti solo controlli di congruenza che possono portare ad errori di validazione
// Nel caso di errore di validazione, la spedizione restituita ha ID == -1
func newSpedizione(requestID int, mobiles model.MobileList) *Spedizione {

	// funzione ausiliaria per restituire una spedizione con errore di validazione
	fnErrValidazione := func(descrErr string) *Spedizione {
		// crea la nuova spedizione
		sped := &Spedizione{
			ID:            -1,
			RequestID:     requestID,
			MobileNumbers: mobiles,
			Status:        SpedizioneErrore,
			Esito:         ErrValidazione,
			EsitoDescr:    descrErr,
			DataCreazione: util.SimulDate(),
			DataEsito:     util.SimulDate(),
		}
		return sped
	}

	// check mobiles list
	if len(mobiles) == 0 {
		return fnErrValidazione("mobile list must not be empty")
	}

	// check mobiles attributes
	for _, mob := range mobiles {
		if mob.MigrStatus != model.MigrDaSpedireOPL {
			msg := fmt.Sprintf("Invalid mobile status: expected MigrDaSpedireOPL (%d) found (%d)", model.MigrDaSpedireOPL, mob.MigrStatus)
			return fnErrValidazione(msg)
		}
		if mob.UsimFMVNO != nil {
			return fnErrValidazione("UsimFMVNO must be nil")
		}
	}

	// crea la nuova spedizione
	sped := &Spedizione{
		ID:            nextSpedizioneID(),
		RequestID:     requestID,
		Status:        SpedizioneDaInviare,
		Esito:         EsitoNonDisponibile,
		DataCreazione: util.SimulDate(),
		MobileNumbers: mobiles,
		UsimNumbers:   make(model.UsimList, len(mobiles)),
	}

	// associa (crea) una nuova USIM a ciascun numero mobile
	for j := range mobiles {
		// storicizza la usim associata al mobile per questa spedizione
		sped.UsimNumbers[j] = model.NewUsim()
	}

	return sped
}
