package fw

import (
	"fmt"
	"math/rand"

	"github.com/mmbros/fmvno/model"
	"github.com/mmbros/mmsync"
)

// Accounts type
type Accounts struct {
	Acc               map[int]*model.Account
	TotMobiles        int
	TotCluster        int
	StoricoSpedizioni *Spedizioni
	ListDaSpedire     accountsDaSpedireType

	poolOfMutex  *mmsync.PoolMutex // used by mutexByAccID
	mutexByAccID mmsync.MutexInt

	ChanSpedRequest chan *model.Spedizione
}

// NewAccounts is ...
func NewAccounts(totAccounts, totMobiles, totCluster int) *Accounts {
	if totMobiles < 0 {
		panic("engine.NewAccounts: totMobiles must be positive")
	}
	if totMobiles < totAccounts {
		panic("engine.NewAccounts: totMobiles must be greater or equal than totAccounts")
	}

	pool, _ := mmsync.NewPoolMutex(10, 100)
	muxint := mmsync.NewPoolMutexInt(pool)

	a := &Accounts{
		Acc:               newAccountsMap(totAccounts, totMobiles),
		TotMobiles:        totMobiles,
		TotCluster:        totCluster,
		poolOfMutex:       pool,
		mutexByAccID:      muxint,
		ChanSpedRequest:   make(chan *model.Spedizione),
		StoricoSpedizioni: NewSpedizioni(),
	}

	// inizializza la lista degli account da spedire
	a.ListDaSpedire = newAccountsDaSpedire(a.Acc)
	a.ListDaSpedire.clusterizza(totCluster)

	return a
}

// crea
func newAccountsMap(totAccounts, totMobiles int) map[int]*model.Account {

	accounts := make(map[int]*model.Account, totAccounts)

	// lista temporanea degli account ID
	// NB: per generalità suppone che l'ID non sia
	// una sequenza consecutiva di int
	listIDs := make([]int, totAccounts)

	// crea gli account ognuno con un solo mobile
	for j := 0; j < totAccounts; j++ {
		acc := model.NewAccount()

		acc.AddMobileNumber(NewMobile())
		accounts[acc.ID] = acc

		listIDs[j] = acc.ID
	}

	// crea i mobili mancanti
	for j := 0; j < totMobiles-totAccounts; j++ {
		// seleziona una posizione a caso nella lista degli id
		n := rand.Intn(totAccounts)
		// recupera l'account
		acc := accounts[listIDs[n]]
		// aggiunge un mobile
		acc.AddMobileNumber(NewMobile())
	}

	return accounts
}

func (a *Accounts) String() string {
	return fmt.Sprintf("Accounts{accounts: %d, mobiles:%d, cluster:%d}", len(a.Acc), a.TotMobiles, a.TotCluster)
}

// Spedisci effettua fino a maxSped spedizioni.
// Il numero può essere inferiore in base alle spedizioni rimaste da effettuare.
// NOTA BENE: la gestione della ListDaSpedire non è thread safe
func (a *Accounts) Spedisci(maxSped int) {
	num := 0
	for _, acc := range a.ListDaSpedire {
		if num >= maxSped {
			break
		}
		// XXX da provare la goroutine
		// go a.doSpedisci(acc)
		a.doSpedisci(acc, -1)
		num++
	}
	// i primi 'num' sono stati spediti
	a.ListDaSpedire = a.ListDaSpedire[num:]
}

// Effettua la spedizione dell'account `acc`.
// Nel caso di prima spedizione, cerca di spedire TUTTI i mobili dell'account
// Nel caso di spedizione successive, cerca di spedire `totMobiles` dell'account.
// Il numero di mobili effettivamente spedito puo' essere inferiore a quanto
// richiesto in base allo stato dei mobili stessi.
//
// La funzione è thread safe (l'account viene lockato)
func (a *Accounts) doSpedisci(acc *model.Account, totMobiles int) {
	fmt.Printf("doSpedisci: acc=%v, totMobiles=%d\n", acc, totMobiles)

	// lock account
	a.mutexByAccID.Lock(acc.ID)
	// defer unlock account
	defer a.mutexByAccID.Unlock(acc.ID)

	var mobiles model.MobileList
	if acc.MigrAccStatus == model.MigrAccDaSpedire {
		// prima spedizione
		mobiles = firstGetMobiliDaSpedire(acc)
		// avanza lo stato migrazione dell'account
		acc.MigrAccStatus = model.MigrAccInCorso
	} else {
		// spedizioni successive
		mobiles = nextGetMobiliDaSpedire(acc, totMobiles)
	}

	if len(mobiles) == 0 {
		// nessun mobile dell'account può essere spedito
		// in ogni caso, se necc, la prima spedizione si considera fatta
		return
	}

	// crea la nuova spedizione
	sped := model.NewSpedizione(mobiles)

	// check errori creazione spedizione (eventuali errori di validazione)
	if sped.Status == model.SpedizioneErrore {
		panic(fmt.Sprintf("Errore spedizione: %d - %s", sped.Esito, sped.EsitoDescr))
	}

	// aggiunge la spedizione alla lista delle spedizioni-consegne dell'account
	acc.SpedizCons.AddSpedizione(sped)

	// aggiunge la spedizione allo storico spedizioni
	a.StoricoSpedizioni.Add(acc, sped)

	// invia la spedizione
	fmt.Printf("Accounts.doSpedisci: %v\n", sped)
	a.ChanSpedRequest <- sped
}

// HandleRispostaSpedizioni gestisce gli esiti delle spedizione inviate dall'OPL
func (a *Accounts) HandleRispostaSpedizioni(response <-chan *model.Spedizione) {
	for sped := range response {

		a.StoricoSpedizioni.Update(sped)

		fmt.Printf("Accounts.HandleRispostaSpedizioni %v\n", sped)
	}
}
