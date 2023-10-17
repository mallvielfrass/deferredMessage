package clocker

import "time"

type Message struct {
	Message string
	Time    time.Time
	Id      string
}
