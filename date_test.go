package date

import "testing"

func TestIsLeapYear(t *testing.T) {
	cases := []struct {
		year int
		want bool
	}{
		{2024, true}, {2023, false}, {2000, true}, {1900, false}, {2100, false},
	}
	for _, c := range cases {
		if got := IsLeapYear(c.year); got != c.want {
			t.Errorf("IsLeapYear(%d) = %v, want %v", c.year, got, c.want)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	cases := []struct {
		year, month int
		want        int
	}{
		{2024, 2, 29}, {2023, 2, 28}, {2026, 4, 30}, {2026, 1, 31},
		{2026, 9, 30}, {2026, 12, 31}, {2026, 11, 30},
	}
	for _, c := range cases {
		if got := DaysInMonth(c.year, c.month); got != c.want {
			t.Errorf("DaysInMonth(%d, %d) = %d, want %d", c.year, c.month, got, c.want)
		}
	}
}

func TestWeekday(t *testing.T) {
	cases := []struct {
		year, month, day int
		want             int // 0 = domingo
	}{
		{1970, 1, 1, 4},   // jueves — ancla del epoch
		{2025, 1, 1, 3},   // miércoles
		{2024, 2, 29, 4},  // jueves (bisiesto)
		{2026, 8, 1, 6},   // sábado
		{2026, 8, 11, 2},  // martes
		{2026, 12, 25, 5}, // viernes
		{2024, 1, 1, 1},   // lunes
	}
	for _, c := range cases {
		if got := Weekday(c.year, c.month, c.day); got != c.want {
			t.Errorf("Weekday(%d, %d, %d) = %d, want %d", c.year, c.month, c.day, got, c.want)
		}
	}
}

func TestAddMonths(t *testing.T) {
	cases := []struct {
		y, m, delta int
		wy, wm      int
	}{
		{2026, 1, -1, 2025, 12},
		{2026, 12, 1, 2027, 1},
		{2026, 1, -13, 2024, 12},
		{2026, 8, 3, 2026, 11},
		{2026, 8, 0, 2026, 8},
	}
	for _, c := range cases {
		gy, gm := AddMonths(c.y, c.m, c.delta)
		if gy != c.wy || gm != c.wm {
			t.Errorf("AddMonths(%d, %d, %d) = (%d, %d), want (%d, %d)", c.y, c.m, c.delta, gy, gm, c.wy, c.wm)
		}
	}
}

func TestParseMonthKey(t *testing.T) {
	cases := []struct {
		in   string
		want [2]int
	}{
		{"2026-08", [2]int{2026, 8}},
		{"2026-08-11", [2]int{2026, 8}}, // tolera la forma fecha completa
		{"2026-8", [2]int{0, 0}},
		{"2026-13", [2]int{0, 0}},
		{"2026-00", [2]int{0, 0}},
		{"2026", [2]int{0, 0}},
		{"agosto-2026", [2]int{0, 0}},
		{"", [2]int{0, 0}},
	}
	for _, c := range cases {
		y, m := ParseMonthKey(c.in)
		if y != c.want[0] || m != c.want[1] {
			t.Errorf("ParseMonthKey(%q) = (%d, %d), want (%d, %d)", c.in, y, m, c.want[0], c.want[1])
		}
	}
}

func TestMonthKey(t *testing.T) {
	if got := MonthKey(2026, 8); got != "2026-08" {
		t.Errorf("MonthKey(2026, 8) = %q, want 2026-08", got)
	}
	if got := MonthKey(2026, 1); got != "2026-01" {
		t.Errorf("MonthKey(2026, 1) = %q, want 2026-01", got)
	}
}

func TestDateKey(t *testing.T) {
	if got := DateKey(2026, 8, 1); got != "2026-08-01" {
		t.Errorf("DateKey(2026, 8, 1) = %q, want 2026-08-01", got)
	}
}

func TestMonthName(t *testing.T) {
	if MonthName(1) != "January" {
		t.Error("MonthName(1) should be January")
	}
	if MonthName(8) != "August" {
		t.Error("MonthName(8) should be August")
	}
	if MonthName(12) != "December" {
		t.Error("MonthName(12) should be December")
	}
	if MonthName(13) != "" {
		t.Error("MonthName(13) should be empty")
	}
}

func TestWeekdayName(t *testing.T) {
	if WeekdayName(0) != "Sunday" {
		t.Error("WeekdayName(0) should be Sunday")
	}
	if WeekdayName(1) != "Monday" {
		t.Error("WeekdayName(1) should be Monday")
	}
	if WeekdayName(6) != "Saturday" {
		t.Error("WeekdayName(6) should be Saturday")
	}
	if WeekdayName(7) != "" {
		t.Error("WeekdayName(7) should be empty")
	}
}

func TestParseDateKey(t *testing.T) {
	y, m, d := ParseDateKey("2026-08-18")
	if y != 2026 || m != 8 || d != 18 {
		t.Errorf("ParseDateKey(2026-08-18) = %d, %d, %d, want 2026, 8, 18", y, m, d)
	}
	if y, m, d := ParseDateKey("2026-02-30"); y != 0 || m != 0 || d != 0 {
		t.Errorf("ParseDateKey(2026-02-30) = %d, %d, %d, want 0, 0, 0 (Feb has no 30th)", y, m, d)
	}
	if y, m, d := ParseDateKey("not-a-date"); y != 0 || m != 0 || d != 0 {
		t.Errorf("ParseDateKey(not-a-date) = %d, %d, %d, want 0, 0, 0", y, m, d)
	}
}

func TestDaysSinceEpochAndRoundTrip(t *testing.T) {
	cases := []struct {
		y, m, d  int
		wantDays int64
	}{
		{1970, 1, 1, 0},
		{1970, 1, 2, 1},
		{1969, 12, 31, -1},
		{2000, 1, 1, 10957},
		{2024, 2, 29, 19782}, // Leap year
		{2025, 1, 6, 20094},
		{2026, 3, 2, 20514},
	}
	for _, c := range cases {
		got := DaysSinceEpoch(c.y, c.m, c.d)
		if got != c.wantDays {
			t.Errorf("DaysSinceEpoch(%d, %d, %d) = %d, want %d", c.y, c.m, c.d, got, c.wantDays)
		}
		ry, rm, rd := FromDaysSinceEpoch(got)
		if ry != c.y || rm != c.m || rd != c.d {
			t.Errorf("FromDaysSinceEpoch(%d) = (%d, %d, %d), want (%d, %d, %d)", got, ry, rm, rd, c.y, c.m, c.d)
		}
	}
}
