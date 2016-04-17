package fw

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/mmbros/fmvno/model"
)

type typeAccSped struct {
	acc  *model.Account
	sped *model.Spedizione
}

// Spedizioni ...
type Spedizioni struct {
	accsped map[int]typeAccSped
	mu      sync.Mutex
}

// NewSpedizioni ..
func NewSpedizioni() *Spedizioni {
	return &Spedizioni{
		accsped: make(map[int]typeAccSped),
		mu:      sync.Mutex{},
	}
}

// Add ..
func (s *Spedizioni) Add(acc *model.Account, sped *model.Spedizione) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accsped[sped.ID] = typeAccSped{acc, sped}
}

// Update ..
func (s *Spedizioni) Update(sped *model.Spedizione) {
	s.mu.Lock()
	defer s.mu.Unlock()
	accsped, ok := s.accsped[sped.ID]
	if !ok {
		panic(fmt.Sprintf("Spedizioni.Update: spedizione non trovata (ID=%d)", sped.ID))
	}
	accsped.sped = sped
}

// Stat restituisce lo storico delle spedizioni raggruppate per stato
//	0 - SpedizioneDaInviare
//	1 - SpedizioneInCorso
//	2 - SpedizioneConsegnata
//	3 - SpedizioneErrore
func (s *Spedizioni) Stat() [4][]*model.Spedizione {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res [4][]*model.Spedizione
	var initSize = len(s.accsped) / 8
	for j := 0; j < 4; j++ {
		res[j] = make([]*model.Spedizione, 0, initSize)
	}

	for _, accsped := range s.accsped {
		status := accsped.sped.Status
		res[status] = append(res[status], accsped.sped)
	}

	return res
}

// gestione prima spedizione dell'account
func firstGetMobiliDaSpedire(acc *model.Account) model.MobileList {

	// check account status
	if acc.MigrAccStatus != model.MigrAccDaSpedire {
		panic(fmt.Sprintf("Invalid account status: expected %s, found %s", model.MigrAccDaSpedire, acc.MigrAccStatus))
	}

	L := len(acc.MobileNumbers)
	mobiles := make(model.MobileList, 0, L)

	for _, mob := range acc.MobileNumbers {
		if mob.MigrStatus == model.MigrDaConsegnare {
			mob.MigrStatus = model.MigrDaSpedireOPL
			mobiles = append(mobiles, mob)
		}
	}

	return mobiles
}

// cerca totMobiles da spedire fra quelli dell'account.
// se totMobiles < 0, restituisce tutti i mobili (comaptibili con la spedizione)
// Spedisce solo in mobili che hanno uno stato copatibile con la spedizione.
// Aggiorna lo stato migrazione dei mobili spediti (MigrDaSpedireOPL)
func nextGetMobiliDaSpedire(acc *model.Account, totMobiles int) model.MobileList {

	// check account status
	if acc.MigrAccStatus != model.MigrAccInCorso {
		panic(fmt.Sprintf("Invalid account status: expected %s, found %s", model.MigrAccInCorso, acc.MigrAccStatus))
	}

	if totMobiles == 0 {
		return nil
	}

	L := len(acc.MobileNumbers)
	// se totMobiles < 0, prova a spedire tutti i mobili dell'account
	// verifica che totMobiles non sia maggiore del numero di mobili dell'account
	if totMobiles < 0 || totMobiles > L {
		totMobiles = L
	}

	mobiles := make(model.MobileList, 0, totMobiles)

	// inizia dall'elemento n-esimo (casuale)
	n := rand.Intn(L)

	count := 0 // conta i mobili trovati da spedire

loop:
	for j := 0; j < L; j++ {
		mob := acc.MobileNumbers[n]

		if mob.MigrStatus.AmmetteNuovaSpedizione() {
			mob.MigrStatus = model.MigrDaSpedireOPL
			mobiles = append(mobiles, mob)
			count++
			if count >= totMobiles {
				break loop
			}
		}
		n = (n + 1) % L
	}

	return mobiles
}
