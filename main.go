//go:generate stringer -output model/model_string.go -type AccountStatus,MigrationAccountStatus,SpedizioneStatus,SpedizioneEsito model
package main

import (
	"fmt"

	"github.com/mmbros/fmvno/fw"
	"github.com/mmbros/fmvno/opl"
	"github.com/mmbros/fmvno/util"
	"github.com/spf13/viper"
)

func initConfig(name string) {
	viper.SetConfigName(name) // name of config file (without extension)
	viper.AddConfigPath(".")  // path to look for the config file in

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	util.SetSimulDate(viper.GetTime("spedizioni.data_inizio"))
	fmt.Println(util.SimulDate())

	fw.InitConfigEngineMobile()

}

func main() {
	initConfig("config")

	accounts := fw.NewAccounts(
		viper.GetInt("account.tot_account"),
		viper.GetInt("account.tot_usim"),
		viper.GetInt("account.tot_cluster"))
	fmt.Printf("%v\n", accounts)

	opl := opl.InitConfigOPL()

	go opl.HandleRichiestaSpedizioni(accounts.ChanSpedRequest)
	go accounts.HandleRispostaSpedizioni(opl.RespChan)

	var totInvio, totOK, totErr int

	for iter := 0; iter < 40; iter++ {
		accounts.Spedisci(5)
		nInvio := opl.HandleInvioSpedizioni()
		nOK, nErr := opl.HandleEsitoSpedizioni()

		totInvio += nInvio
		totOK += nOK
		totErr += nErr

		fmt.Printf("ITER #%2d: %s - invio=%2d, ok=%2d, err=%2d, Tok=%2d, Terr=%2d, I=%3d, T=%2d\n",
			iter, util.YYYYMMDD(util.SimulDate()),
			nInvio, nOK, nErr, totOK, totErr,
			totInvio, totOK+totErr)

		util.IncSimulDate(1)
	}

}
