package fw

import (
	"math/rand"
	"sort"

	"github.com/mmbros/fmvno/model"
)

// accountsDaSpedireType
type accountsDaSpedireType []*model.Account

func newAccountsDaSpedire(accounts map[int]*model.Account) accountsDaSpedireType {
	totAccounts := len(accounts)

	// popola la lista degli account da spedire
	list := make(accountsDaSpedireType, totAccounts)

	j := 0
	for _, acc := range accounts {
		list[j] = acc
		j++
	}

	return list
}

// clusterizza effettua una nuova clusterizzazione MKTG degli account ancora da spedire
// Non modifica il cluster degli account già spediti la prima volta
// NOTA BENE: la gestione della ListDaSpedire non è thread safe
func (a accountsDaSpedireType) clusterizza(totCluster int) {
	for _, acc := range a {
		// cluster casuale fra 1 .. totCluster
		acc.Cluster = 1 + rand.Intn(totCluster)
	}

	// ordina la lista degli account da spedire
	sort.Sort(byCluster(a))
}

// byCluster implements sort.Interface for []*Account
// usata da Clusterizza
type byCluster []*model.Account

func (a byCluster) Len() int      { return len(a) }
func (a byCluster) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a byCluster) Less(i, j int) bool {
	cmp := a[i].Cluster - a[j].Cluster
	if cmp == 0 {
		cmp = a[i].ID - a[j].ID
	}
	return cmp < 0

}
