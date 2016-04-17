package model

import "github.com/mmbros/fmvno/util"

// UsimStatus ...
type UsimStatus int

// UsimStatus values
const (
	UsimDaAttivare UsimStatus = iota
	UsimAttiva
	UsimDisattiva
)

// Usim ...
type Usim struct {
	ID     int        `json:"usim_id"`
	Status UsimStatus `json:"stato"`
	//	Iccid  int        `json:"iccid"`
	//	Imsi   int        `json:"imsi"` //  International Mobile Subscriber Identity
}

var nextUsimID = util.NewIntSequence()

// NewUsim return a random usim object
func NewUsim() *Usim {
	return &Usim{
		ID:     nextUsimID(),
		Status: UsimDaAttivare,
	}
}

// MCC MNC MSIN
// 222 08  0...9
// http://www.mobileworld.it/2016/01/12/fastweb-tim-full-mvno-62709/

// prefissi msisdn
// fastweb full: 3755 e 3756
// fastweb esp:  373

// UsimList is ...
type UsimList []*Usim
