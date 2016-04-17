package model

// spedizConsegna
type spedizConsegna struct {
	sped *Spedizione
	cons *Consegna
}

// SpedizConsegnaList ...
type SpedizConsegnaList []*spedizConsegna

// AddConsegna ...
func (l *SpedizConsegnaList) AddConsegna(c *Consegna) {
	if c == nil {
		panic("Consegna must not be nil")
	}
	if l == nil {
		l = &SpedizConsegnaList{}
	}

	*l = append(*l, &spedizConsegna{cons: c})
}

// AddSpedizione ...
func (l *SpedizConsegnaList) AddSpedizione(s *Spedizione) {
	if s == nil {
		panic("Spedizione must not be nil")
	}
	if l == nil {
		l = &SpedizConsegnaList{}
	}

	*l = append(*l, &spedizConsegna{sped: s})
}

// GetConsegna ...
func (sc *spedizConsegna) GetConsegna() *Consegna { return sc.cons }

// GetSpedizione ...
func (sc *spedizConsegna) GetSpedizione() *Spedizione { return sc.sped }
