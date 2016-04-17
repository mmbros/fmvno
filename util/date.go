package util

import (
	"time"
	"unicode/utf8"
)

// simulation date
var simulDate time.Time

// SetSimulDate imposta la data di simulazione.
// Tronca la data ai giorni, scartando ore, min, sec, ...
// Restituisce la data di simulazione aggiornata.
func SetSimulDate(date time.Time) time.Time {
	simulDate = date.Truncate(24 * time.Hour)
	return simulDate
}

// IncSimulDate incrementa la data di simulazione di delta giorni.
// Restituisce la data di simulazione aggiornata.
func IncSimulDate(delta int) time.Time {
	simulDate = AddDay(simulDate, delta)
	return simulDate
}

// SimulDate restituisce la data di simulazione corrente.
func SimulDate() time.Time {
	return simulDate
}

// YYYYMMDD restituisce una stringa che rappresenta la data nel formato YYYY-MM-DD
func YYYYMMDD(date time.Time) string {
	return date.Format("2006-01-02")
}

// AddDay returns date + N days.
func AddDay(date time.Time, N int) time.Time {
	return date.Add(time.Duration(N) * 24 * time.Hour)
}

// WeekDayInitial restituisce il carattere iniziale del WeekDay
func WeekDayInitial(date time.Time) string {
	r, _ := utf8.DecodeRuneInString(date.Weekday().String())
	return string(r)
}

// Date return a time.Time date with given year, month, day
func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}
