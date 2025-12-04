package main

import (
	"time"
)

type Message struct {
	FromID    int
	Text      string
	CreatedAt time.Time
}
