package util

import (
	"fmt"
	"time"
)

// Calendar is ...
type Calendar struct {
	SaturdayWorkingDay bool
	SundayWorkingDay   bool
}

// AddWorkDay returns date + N working days.
func (cal *Calendar) AddWorkDay(date time.Time, days int) time.Time {

	for {
		date = AddDay(date, 1)
		if cal.WorkingDay(date) {
			days--
			if days <= 0 {
				return date
			}
		}
	}

}

// WorkingDay returns true if the date is a working day in the calendar
// see: http://www.giorni-festivi.it/
func (cal *Calendar) WorkingDay(date time.Time) bool {
	wd := date.Weekday()
	if (wd == time.Sunday) && !cal.SundayWorkingDay {
		return false
	}
	if (wd == time.Saturday) && !cal.SaturdayWorkingDay {
		return false
	}

	year, month, day := date.Date()

	switch month {
	case time.January:
		if day == 1 || day == 6 { // capodanno o befana
			return false
		}
	case time.April:
		if day == 25 { // festa della liberazione
			return false
		}
	case time.May:
		if day == 1 { // festa del lavoro
			return false
		}
	case time.June:
		if day == 2 { // festa della repubblica
			return false
		}
	case time.August:
		if day == 15 { // ferragosto
			return false
		}
	case time.November:
		if day == 1 { // ognisanti
			return false
		}
	case time.December:
		if day == 8 || day == 25 || day == 26 { // immacolata, natale, s.stefano
			return false
		}
	}

	// pasqua e pasquetta
	switch year {
	case 2016:
		if (month == time.March) && (day == 27 || day == 28) {
			return false
		}
	case 2017:
		if (month == time.April) && (day == 16 || day == 17) {
			return false
		}
	case 2018:
		if (month == time.April) && (day == 1 || day == 2) {
			return false
		}
	case 2019:
		if (month == time.April) && (day == 21 || day == 22) {
			return false
		}
	case 2020:
		if (month == time.April) && (day == 12 || day == 13) {
			return false
		}
	}

	return true
}

// PrintTest show the calendar test
func (cal *Calendar) PrintTest(date time.Time, days int) {
	for i := 0; i < days; i++ {
		wd := cal.WorkingDay(date)

		d2 := cal.AddWorkDay(date, 2)

		fmt.Printf("%02d: %s [%v] %-5v %s\n", i, YYYYMMDD(date), WeekDayInitial(date), wd, YYYYMMDD(d2))
		date = AddDay(date, 1)
	}
}
