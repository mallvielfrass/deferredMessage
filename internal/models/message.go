package models

import "time"

type Message struct {
	Message     string
	Time        time.Time
	Id          string
	ChatId      string
	CreatorId   string
	IsProcessed bool
	IsSended    bool
	Error       string
}
