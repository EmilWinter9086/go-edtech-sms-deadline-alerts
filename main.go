package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	client, err := NewSMSClient()
	if err != nil {
		log.Fatal(err)
	}
	phone := os.Getenv("LEARNER_PHONE")
	if phone == "" {
		log.Fatal("LEARNER_PHONE is required")
	}
	course := Course{Title: "Algebra fundamentals", Deadline: time.Now().Add(8 * time.Hour)}
	messageID, err := SendDeadlineReminder(client, course, Learner{Name: "Mina", Phone: phone}, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	if messageID == "" {
		fmt.Println("no reminder needed")
	} else {
		fmt.Println("deadline reminder sent:", messageID)
	}
}
