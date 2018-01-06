package main

import "time"

type message struct {
	Name      string
	Message   string
	Time      time.Time
	AvatarURL string
}
