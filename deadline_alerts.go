package main

import (
	"fmt"
	"time"
)

type Learner struct{ Name, Phone string }
type Course struct {
	Title    string
	Deadline time.Time
}

func Reminder(course Course, learner Learner, now time.Time) (string, bool) {
	remaining := course.Deadline.Sub(now)
	if remaining < 0 || remaining > 24*time.Hour {
		return "", false
	}
	return fmt.Sprintf("%s: %s is due by %s", learner.Name, course.Title, course.Deadline.Format("Jan 2 15:04 MST")), true
}

func SendDeadlineReminder(client *SMSClient, course Course, learner Learner, now time.Time) (string, error) {
	message, ok := Reminder(course, learner, now)
	if !ok {
		return "", nil
	}
	requestID := fmt.Sprintf("deadline:%s:%s", learner.Phone, course.Deadline.UTC().Format(time.RFC3339))
	return client.Send(learner.Phone, message, requestID)
}
