package fw

import (
	"fmt"
	"math/rand"

	"github.com/mmbros/fmvno/distribution"
	"github.com/mmbros/fmvno/model"
	"github.com/spf13/viper"
)

// NewMobile return a random mobile object
func NewMobile() *model.Mobile {

	return &model.Mobile{
		Msisdn:         newMsisdn(),
		Status:         randMobileStatus(),
		ProfiloOfferta: randProfiloOfferta(),
	}
}

// distribuzioni inizializzate da InitConfigEngineMobile
var mDistProfiloOfferta, mDistStatus *distribution.Discrete

// InitConfigEngineMobile ...
func InitConfigEngineMobile() {
	var err error
	mDistProfiloOfferta, err = distribution.NewDiscreteFromStrings(
		viper.GetStringSlice("mobile.profilo_offerta"))
	if err != nil {
		panic(fmt.Sprintf("Reading config [mobile.profilo_offerta]: %s", err))
	}

	mDistStatus, err = distribution.NewDiscreteFromStrings(
		viper.GetStringSlice("mobile.status"))
	if err != nil {
		panic(fmt.Sprintf("Reading config [mobile.status]: %s", err))
	}

}

func randProfiloOfferta() model.TipoProfiloOfferta {
	if mDistProfiloOfferta == nil {
		panic("mDistProfiloOfferta == nil: InitConfigEngineMobile non called?")
	}
	n := mDistProfiloOfferta.RandInt()
	return model.TipoProfiloOfferta(n)
}

func randMobileStatus() model.MobileStatus {
	if mDistStatus == nil {
		panic("mDistStatus == nil: InitConfigEngineMobile non called?")
	}
	n := mDistStatus.RandInt()
	return model.MobileStatus(n)
}

// module variable usata solo da newMsisdn
var seqMsisdn int

// restituisce un nuovo Msisdn
// fa in modo che i numeri non si ripetano
// prefisso casuale e il resto da una sequence
func newMsisdn() model.MsisdnType {
	const multi int = 1000000
	prefix := 320 + rand.Intn(80)
	seqMsisdn++
	n := prefix*multi + seqMsisdn
	return model.MsisdnType(n)
}
