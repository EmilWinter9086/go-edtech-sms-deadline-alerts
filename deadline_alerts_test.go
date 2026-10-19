package main

import (
	"testing"
	"time"
)

func TestReminderWindow(t *testing.T) {
	now := time.Date(2026, 9, 9, 9, 0, 0, 0, time.UTC)
	cases := []struct {
		name  string
		hours float64
		want  bool
	}{{"inside", 6, true}, {"outside", 30, false}, {"past", -1, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, got := Reminder(Course{Title: "Algebra", Deadline: now.Add(time.Duration(tc.hours) * time.Hour)}, Learner{Name: "Mina"}, now)
			if got != tc.want {
				t.Fatalf("want %v, got %v", tc.want, got)
			}
		})
	}
}
