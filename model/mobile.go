package model

// MobileStatus ...
type MobileStatus int

// MsisdnType ...
type MsisdnType int

// MigrationStatus ...
type MigrationStatus int

// TipoProfiloOfferta ...
type TipoProfiloOfferta int

// UsimStatus values
const (
	MobileAttivo MobileStatus = iota
	MobileSospeso
	MobileSilente
)

// TipoProfiloOfferta values
const (
	VoceDati TipoProfiloOfferta = iota
	SoloDati
)

// MigrationStatus values
const (
	MigrDaConsegnare MigrationStatus = iota
	MigrDaSpedireOPL
	MigrInviatoOPL
	MigrInSpedizioneOPL
	MigrConsegnatoOPL
	MigrErroreOPL
	MigrConsegnatoDealer
	MigrTecnRichiesta
	MigrTecnInCorso
	MigrTecnCompletata
)

// AmmetteNuovaSpedizione restituisce true se lo stato è compatibile con una nuova spedizione
func (s MigrationStatus) AmmetteNuovaSpedizione() bool {
	switch s {
	case MigrDaConsegnare, MigrDaSpedireOPL, MigrConsegnatoOPL, MigrErroreOPL, MigrConsegnatoDealer:
		return true
	}
	return false
}

// Mobile is ...
type Mobile struct {
	Msisdn         MsisdnType         `json:"msisdn"`
	Status         MobileStatus       `json:"stato_numero,omitempty"`
	MigrStatus     MigrationStatus    `json:"stato_migraz,omitempty"`
	ProfiloOfferta TipoProfiloOfferta `json:"profilo_offerta,omitempty"`
	UsimFMVNO      *Usim              `json:"usim,omitempty"`
	//	ID             int                `json:"mobile_id"`
	//	AccountID      int                `json:"account_id"`
}

// MobileList is ...
type MobileList []*Mobile
