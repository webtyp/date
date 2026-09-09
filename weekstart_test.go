package date

import "testing"

// restore puts the package default back, so one test's SetFirstWeekday never
// leaks into the next.
func restore(t *testing.T) {
	t.Helper()
	t.Cleanup(func() { SetFirstWeekday(Monday) })
}

// TestFirstWeekday_DefaultsToMonday is the point of the whole file: a calendar
// that opens on Sunday reads as wrong in every Spanish-speaking locale, and
// ISO-8601 agrees. Nobody should have to configure this to get it right.
func TestFirstWeekday_DefaultsToMonday(t *testing.T) {
	if got := FirstWeekday(); got != Monday {
		t.Errorf("FirstWeekday() = %d (%s), want Monday", got, WeekdayName(got))
	}
}

// TestWeekOrder_StartsAtMondayAndEndsAtSunday: the order a component renders.
func TestWeekOrder_StartsAtMondayAndEndsAtSunday(t *testing.T) {
	restore(t)
	want := [7]int{Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday}
	if got := WeekOrder(); got != want {
		t.Errorf("WeekOrder() = %v, want %v", got, want)
	}
}

// TestWeekColumn_IsWeekOrdersInverse: the two must agree, or a calendar's
// header row and its day cells drift apart by a column.
func TestWeekColumn_IsWeekOrdersInverse(t *testing.T) {
	restore(t)
	for col, day := range WeekOrder() {
		if got := WeekColumn(day); got != col {
			t.Errorf("WeekColumn(%s) = %d, want %d", WeekdayName(day), got, col)
		}
	}
}

// TestSetFirstWeekday_MovesBothInStep: an app that reads its calendars the US
// way changes one call and both helpers follow.
func TestSetFirstWeekday_MovesBothInStep(t *testing.T) {
	restore(t)
	SetFirstWeekday(Sunday)

	if got := FirstWeekday(); got != Sunday {
		t.Fatalf("FirstWeekday() = %d, want Sunday", got)
	}
	want := [7]int{Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday}
	if got := WeekOrder(); got != want {
		t.Errorf("WeekOrder() = %v, want %v", got, want)
	}
	if got := WeekColumn(Sunday); got != 0 {
		t.Errorf("WeekColumn(Sunday) = %d, want 0", got)
	}
	if got := WeekColumn(Saturday); got != 6 {
		t.Errorf("WeekColumn(Saturday) = %d, want 6", got)
	}
}

// TestSetFirstWeekday_IgnoresOutOfRange: a bad argument leaves a working
// calendar standing instead of producing one with no first column.
func TestSetFirstWeekday_IgnoresOutOfRange(t *testing.T) {
	restore(t)
	for _, bad := range []int{-1, 7, 99} {
		SetFirstWeekday(bad)
		if got := FirstWeekday(); got != Monday {
			t.Errorf("SetFirstWeekday(%d) changed the default to %d", bad, got)
		}
	}
}

// TestWeekColumn_MatchesTheInlineFormulaItReplaces: components/calendarslider
// carried "(Weekday(y,m,1) + 6) %% 7" with a "0 = lunes" comment. WeekColumn
// must produce exactly that under the Monday default, or migrating the caller
// shifts every month grid by a column.
func TestWeekColumn_MatchesTheInlineFormulaItReplaces(t *testing.T) {
	restore(t)
	for w := Sunday; w <= Saturday; w++ {
		if got, want := WeekColumn(w), (w+6)%7; got != want {
			t.Errorf("WeekColumn(%s) = %d, inline formula gives %d", WeekdayName(w), got, want)
		}
	}
}
