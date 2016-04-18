package model

import (
	"fmt"

	"github.com/mmbros/fmvno/util"
)

// AccountStatus ...
type AccountStatus int

// AccountStatus values
const (
	AccountAttivo AccountStatus = iota
	AccountSospeso
)

// MigrationAccountStatus ...
type MigrationAccountStatus int

// MigrationAccountStatus values
const (
	MigrAccDaSpedire MigrationAccountStatus = iota
	MigrAccInCorso
	MigrAccTerminata
)

// Account ...
type Account struct {
	ID            int                `json:"account_id"`
	Status        AccountStatus      `json:"status"`
	MobileNumbers MobileList         `json:"mobile_numbers"`
	SpedizCons    SpedizConsegnaList `json:"spediz_cons"`
	MigrAccStatus MigrationAccountStatus
	Cluster       int
}

var nextAccountID = util.NewIntSequence()

// NewAccount return a random mobile object
func NewAccount() *Account {

	return &Account{
		ID:            nextAccountID(),
		Status:        AccountAttivo,
		MobileNumbers: MobileList{},
	}
}

// AddMobileNumber add the mobile number n to the account's numbers
func (a *Account) AddMobileNumber(m *Mobile) {
	a.MobileNumbers = append(a.MobileNumbers, m)
}

func (a *Account) String() string {
	/*
		msisdns := make([]MsisdnType, len(a.MobileNumbers))
		for j, m := range a.MobileNumbers {
			msisdns[j] = m.Msisdn
		}

		return fmt.Sprintf("Account{ID: %d, St:%d, Clu:%d, Mob:%v}", a.ID, a.Status, a.Cluster, msisdns)
	*/
	return fmt.Sprintf("Account{ID:%d, St:%d, Clu:%d, Mob:%d}", a.ID, a.Status, a.Cluster, len(a.MobileNumbers))
}
