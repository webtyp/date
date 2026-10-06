// Package date is pure calendar arithmetic for Go and TinyGo: weekday of a
// date, days in a month, leap years, month names, and the "YYYY-MM"/
// "YYYY-MM-DD" key format used across the webtyp ecosystem to identify a
// month or a day. No time zone, no wall clock, no JS interop — for that, use
// webtyp.com/time. Every function here is a pure computation over
// (year, month, day) ints, so it runs identically in WASM and on the
// backend, and needs no build-tag split.
//
// Every calendar-unit name this package returns (MonthName, WeekdayName) is
// English — the canonical, untranslated form and the translation key. The
// translations ship as data in this module's lang.json (merged into the page by
// sitec); a consumer shows a name with lang.Translate (webtyp.com/lang). This
// package does not import lang or choose a language itself.
package date

import "webtyp.com/fmt"

// IsLeapYear reports whether year is a leap year.
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysInMonth returns the number of days in month (1-12) for year.
func DaysInMonth(year, month int) int {
	switch month {
	case 2:
		if IsLeapYear(year) {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 31
	}
}

// Weekday returns the day of week for a date (0 = Sunday … 6 = Saturday)
// using the Sakamoto algorithm: pure arithmetic, no time zone, no parsing —
// deterministic in WASM and on the backend.
func Weekday(year, month, day int) int {
	t := [12]int{0, 3, 2, 5, 0, 3, 5, 1, 4, 6, 2, 4}
	y := year
	if month < 3 {
		y--
	}
	w := (y + y/4 - y/100 + y/400 + t[month-1] + day) % 7
	if w < 0 {
		w += 7
	}
	return w
}

// AddMonths adds delta months to (year, month), rolling the year over as
// needed — the only date arithmetic that depends on neither time zone nor
// the calendar (unlike day arithmetic, a month always has exactly 12 steps
// per year).
func AddMonths(year, month, delta int) (int, int) {
	t := year*12 + (month - 1) + delta
	y := t / 12
	m := t % 12
	if m < 0 {
		m += 12
		y--
	}
	return y, m + 1
}

// MonthName returns month's English name ("January".."December"), or "" if
// month is out of range. English is the canonical, untranslated form — see
// the package doc comment for why. Switch, not a map — TinyGo.
func MonthName(month int) string {
	switch month {
	case 1:
		return "January"
	case 2:
		return "February"
	case 3:
		return "March"
	case 4:
		return "April"
	case 5:
		return "May"
	case 6:
		return "June"
	case 7:
		return "July"
	case 8:
		return "August"
	case 9:
		return "September"
	case 10:
		return "October"
	case 11:
		return "November"
	case 12:
		return "December"
	default:
		return ""
	}
}

// WeekdayName returns w's English name ("Sunday".."Saturday"), or "" if w
// is out of range. w follows Weekday's convention: 0 = Sunday … 6 =
// Saturday. English is the canonical, untranslated form, same reasoning as
// MonthName. Switch, not a map — TinyGo.
func WeekdayName(w int) string {
	switch w {
	case 0:
		return "Sunday"
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	case 6:
		return "Saturday"
	default:
		return ""
	}
}

// Weekday indexes, the convention Weekday returns and WeekdayName reads.
// They exist so a call site says which day it means: SetFirstWeekday(Monday)
// instead of SetFirstWeekday(1).
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// firstWeekday is the day a week starts on for presentation. Monday by
// default: most of the world — and every Spanish-speaking locale — reads a
// calendar that way, and a grid that opens on Sunday reads as wrong there.
// The ISO-8601 week starts on Monday too.
//
// This is the ONE place that decision lives. Before it existed, each consumer
// re-derived it inline (components/calendarslider carried a bare
// "(Weekday(y,m,1) + 6) % 7 // 0 = lunes"), which is how two calendars in one
// app end up disagreeing and why changing it meant editing every one of them.
var firstWeekday = Monday

// FirstWeekday reports the weekday a week starts on. Monday unless the app
// changed it.
func FirstWeekday() int {
	return firstWeekday
}

// SetFirstWeekday changes the day a week starts on, for a locale that reads
// its calendars differently — SetFirstWeekday(Sunday) for the US. Out-of-range
// values are ignored, so a bad argument leaves the default standing instead of
// producing a calendar with no first column.
//
// Call it once, where the app configures its locale (the same file that calls
// lang.OutLang), never from a library: a library that sets it decides for every
// app that imports it.
func SetFirstWeekday(w int) {
	if w < Sunday || w > Saturday {
		return
	}
	firstWeekday = w
}

// WeekColumn returns the 0-based column w occupies in a week laid out from
// FirstWeekday — 0 for the first column, 6 for the last. With the Monday
// default, Monday is 0 and Sunday is 6.
//
// This is what a calendar grid needs to place a month's first day, and what a
// day picker needs to sort its chips.
func WeekColumn(w int) int {
	c := (w - firstWeekday) % 7
	if c < 0 {
		c += 7
	}
	return c
}

// WeekOrder returns the seven weekday indexes in presentation order, starting
// at FirstWeekday — the inverse of WeekColumn. With the Monday default:
// [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday].
//
// A component renders its day headers or day chips by ranging over this
// instead of counting 0..6, which is what makes the order follow the app's
// locale rather than each component's own assumption.
func WeekOrder() [7]int {
	var out [7]int
	for i := 0; i < 7; i++ {
		out[i] = (firstWeekday + i) % 7
	}
	return out
}

// ParseMonthKey reads "YYYY-MM" (or the "YYYY-MM-DD" form, ignoring the
// day); returns (0, 0) if s is not a valid month key.
func ParseMonthKey(s string) (year, month int) {
	if len(s) < 7 || s[4] != '-' {
		return 0, 0
	}
	y, err := fmt.Convert(s[:4]).Int()
	if err != nil {
		return 0, 0
	}
	m, err := fmt.Convert(s[5:7]).Int()
	if err != nil || m < 1 || m > 12 {
		return 0, 0
	}
	return y, m
}

// ParseDateKey reads "YYYY-MM-DD"; returns (0, 0, 0) if s is not a valid
// date key — including a day that does not exist in that month (e.g.
// "2026-02-30").
func ParseDateKey(s string) (year, month, day int) {
	y, m := ParseMonthKey(s)
	if y == 0 || len(s) < 10 || s[7] != '-' {
		return 0, 0, 0
	}
	d, err := fmt.Convert(s[8:10]).Int()
	if err != nil || d < 1 || d > DaysInMonth(y, m) {
		return 0, 0, 0
	}
	return y, m, d
}

// MonthKey formats (year, month) as "YYYY-MM".
func MonthKey(year, month int) string {
	return fmt.Sprintf("%04d-%02d", year, month)
}

// DateKey formats (year, month, day) as "YYYY-MM-DD".
func DateKey(year, month, day int) string {
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}
