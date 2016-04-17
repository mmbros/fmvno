package opl

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/mmbros/fmvno/model"
	"github.com/mmbros/fmvno/util"
)

type spedInCorsoType struct {
	sped                 *model.Spedizione
	saveDataEsito        time.Time
	saveElapsedGiorniLav int
}

type listaSpedInCorsoType []*spedInCorsoType

// OPL type
type OPL struct {
	MaxNumSpedAlGiorno int // must be >= 0
	MinNumSpedAlGiorno int // must be >= 0
	NumSpedAlGiorno    util.RandIntByDateFunc
	GiorniLavConsegne  util.RandIntByDateFunc
	EsitoConsegne      util.RandIntByDateFunc
	CalSpedizioni      util.Calendar
	RespChan           chan *model.Spedizione

	listSpedDaInviare SpedizioneList
	listSpedInCorso   listaSpedInCorsoType
	mu                *sync.RWMutex
	reqChanClosed     bool
	respChanClosed    bool
}

/*
func (opl *OPL) HandleRichiestaSpedizioni(requests <-chan SpedRequest)
func (opl *OPL) HandleInvioSpedizioni() int
func (opl *OPL) HandleEsitoSpedizioni() (numOK, numErr int)
*/

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// DefaultEsitoConsegne is ...
func DefaultEsitoConsegne(time.Time) int {
	var esito int // 0 = ConsegaOK
	n := rand.Intn(1000)
	if n > 900 { // 90% di esiti positivi
		esito = 1 + rand.Intn(5)
	}
	//fmt.Println("DefaultEsitoConsegne =", esito)
	return esito
}

// NewOPL ...
func NewOPL(
	minNumSpedAlGiorno, maxNumSpedAlGiorno int,
	numSpedAlGiorno, giorniLavConsegne, esitoConsegne util.RandIntByDateFunc,
	calSpedizioni util.Calendar) *OPL {
	if minNumSpedAlGiorno < 0 {
		minNumSpedAlGiorno = 0
	}
	if maxNumSpedAlGiorno <= 0 {
		maxNumSpedAlGiorno = 1
	}
	if maxNumSpedAlGiorno < minNumSpedAlGiorno {
		maxNumSpedAlGiorno = minNumSpedAlGiorno
	}

	if esitoConsegne == nil {
		esitoConsegne = DefaultEsitoConsegne
	}

	opl := &OPL{
		MaxNumSpedAlGiorno: maxNumSpedAlGiorno,
		MinNumSpedAlGiorno: minNumSpedAlGiorno,
		NumSpedAlGiorno:    numSpedAlGiorno,
		GiorniLavConsegne:  giorniLavConsegne,
		EsitoConsegne:      esitoConsegne,
		CalSpedizioni:      calSpedizioni,
		listSpedDaInviare:  make(SpedizioneList, 0, 1000),
		listSpedInCorso:    make(listaSpedInCorsoType, 0, 1000),
		mu:                 &sync.RWMutex{},
		RespChan:           make(chan *Spedizione),
	}

	return opl

}

// HandleRichiestaSpedizioni gestisce le richieste di spedizione da parte dei client
func (opl *OPL) HandleRichiestaSpedizioni(requests <-chan *model.Spedizione) {
	for sped := range requests {
		fmt.Printf("opl.HandleRichiestaSpedizioni %v\n", sped)

		// aggiunge la nuova spedizione alle spedizioni da inviare
		opl.mu.Lock()
		opl.listSpedDaInviare = append(opl.listSpedDaInviare, sped)
		opl.mu.Unlock()

		// aggiorna la data invio all'OPL
		sped.DataInvio = util.SimulDate()

		// associa (crea) una nuova USIM a ciascun numero mobile
		sped.UsimNumbers = make(model.UsimList, len(sped.MobileNumbers))
		for _, mob := range sped.MobileNumbers {
			// seleziona (crea) una nuova usim
			usim := model.NewUsim()
			// associa la usim al numero
			mob.UsimFMVNO = usim
			// storicizza la usim associata al mobile per questa spedizione
			sped.UsimNumbers[j] = usim
		}

		opl.RespChan <- sped
	}
	opl.reqChanClosed = true
}

// HandleInvioSpedizioni gestisce l'invio delle spedizioni
// gestisce (o almeno dovrebbe ...) l'accesso concorrente alle risorse
// - opl.listSpedDaInviare : solo rimozione all'inizio della lista
// - opl.listSpedInCorso : solo estensione in fondo alla lista (append)
// restituisce il numero di spedizioni effettuate
func (opl *OPL) HandleInvioSpedizioni() int {
	if opl.respChanClosed {
		return 0
	}

	// data attuale
	date := util.SimulDate()

	// se non siamo in un giorno lavorativo per le spedizioni,
	// non spedisce nulla
	if !opl.CalSpedizioni.WorkingDay(date) {
		return 0
	}

	// determina il numero di spedizioni da effettuare nel giorno
	totDaSpedire := opl.NumSpedAlGiorno(date)
	if totDaSpedire > opl.MaxNumSpedAlGiorno {
		// le spedizioni eccedono il limite massimo:
		// abbassa il numero al massimo consentito
		totDaSpedire = opl.MaxNumSpedAlGiorno
	}

	// se il numero di spedizioni è inferiore al minimo consentito,
	// non spedisce nulla
	// NOTA si suppone che: opl.MinNumSpedAlGiorno >= 0
	if totDaSpedire < opl.MinNumSpedAlGiorno {
		return 0
	}

	// pone il lock per la modifica della opl.listSpedDaInviare
	opl.mu.Lock()
	if L := len(opl.listSpedDaInviare); totDaSpedire > L {
		totDaSpedire = L
	}
	// determina la lista delle spedizioni da inviare a questo giro
	currListSpedDaInviare := opl.listSpedDaInviare[:totDaSpedire]

	// elimina gli elementi inviati dalla lista delle spedizioni da inviare
	opl.listSpedDaInviare = opl.listSpedDaInviare[totDaSpedire:]
	opl.mu.Unlock()

	currListSpedInCorso := make(listaSpedInCorsoType, totDaSpedire)

	// loop delle spedizioni
	for j, sped := range currListSpedDaInviare {

		// aggiorna la spedizione
		sped.Status = model.SpedizioneInCorso
		sped.DataSpedizione = date

		// determina subito i giorni lavorativi e la data per la consegna
		// per far bene le cose :) non li salva direttamente nella spedizione,
		// ma in una struttura ausiliaria
		gglav := opl.GiorniLavConsegne(date)
		if gglav <= 0 {
			panic(fmt.Sprint("giorni lavorativi deve essere maggiore di zero: valore restituito =", gglav))
		}
		dataEsito := opl.CalSpedizioni.AddWorkDay(date, gglav)

		currListSpedInCorso[j] = &spedInCorsoType{
			sped:                 sped,
			saveDataEsito:        dataEsito,
			saveElapsedGiorniLav: gglav,
		}

		// notifica al client l'invio della spedizione
		opl.RespChan <- sped

	}

	// estende la lista deglle spedizioni in corso
	opl.mu.Lock()
	opl.listSpedInCorso = append(opl.listSpedInCorso, currListSpedInCorso...)
	opl.mu.Unlock()

	opl.checkCloseRespChan()
	return totDaSpedire
}

// HandleEsitoSpedizioni ...
// opl.listSpedInCorso: accorcia la lista, gestendo la possibilità che nel frattempo sia stata estesa
func (opl *OPL) HandleEsitoSpedizioni() (numOK, numErr int) {

	if opl.respChanClosed {
		return 0, 0
	}

	// data attuale
	date := util.SimulDate()

	Lorig := len(opl.listSpedInCorso)
	L := Lorig

	// ciclo a ritroso dall'elemento L-esimo - 1
	// gli elementi da 0 a L-1 non possono cambiare al di fuori di questa funzione
	for j := L - 1; j >= 0; j-- {
		sic := opl.listSpedInCorso[j]

		// se l'esito non è della data corrente, passa alla prossima spedizione
		if sic.saveDataEsito != date {
			continue
		}

		sped := sic.sped
		sped.DataEsito = sic.saveDataEsito
		sped.ElapsedGiorniLav = sic.saveElapsedGiorniLav

		// l'elemento deve essere eliminato dalla lista
		// 1. diminuisco la lunghezza della lista
		L--
		if j < L {
			// 2. sovrascrivo l'elemento da eliminare, con l'ultimo elemento da mantenere
			opl.listSpedInCorso[j] = opl.listSpedInCorso[L]
		}

		// esito della spedizione
		sped.Esito = SpedizioneEsito(opl.EsitoConsegne(date))
		sped.DataEsito = date

		// determina lo stato della spedizione in base all'esito
		if sped.Esito == ConsegnaOK {
			sped.Status = SpedizioneConsegnata
			numOK++
		} else {
			sped.Status = SpedizioneErrore
			numErr++
		}

		// notifica al client l'esito della spedizione
		opl.RespChan <- sped
	}

	// 3. ridimensiona la lista
	opl.mu.Lock()
	if len(opl.listSpedInCorso) > Lorig {
		// la lista è stata estesa (in modo concorrente) da un'altra funzione (es. HandleInvioSpedizioni)
		opl.listSpedInCorso = append(opl.listSpedInCorso[0:L], opl.listSpedInCorso[Lorig:]...)
	} else {
		opl.listSpedInCorso = opl.listSpedInCorso[0:L]
	}
	opl.mu.Unlock()

	opl.checkCloseRespChan()
	return
}

// checkCloseRespChan verifica se chiudere il RespChan. Devono essere verificare
// le seguenti condizioni:
// 1. ReqChan è già stato chiuso e perciò non sono previste nuove richieste
// 2. le code delle spedizioni da inviare e in corso sono vuote
func (opl *OPL) checkCloseRespChan() {
	opl.mu.Lock()
	defer opl.mu.Unlock()
	if opl.reqChanClosed && !opl.respChanClosed {
		if (len(opl.listSpedDaInviare) == 0) && (len(opl.listSpedInCorso) == 0) {
			opl.respChanClosed = true
			close(opl.RespChan)
		}
	}
}
