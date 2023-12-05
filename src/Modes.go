package src

import "time"

type Message struct {
	Sender string
	Data   []byte
	Time   time.Time
}

const MStore_TYPE_Manifet = 1

type MStore struct {
	Type int    `json:"type"`
	Data []byte `json:"data"`
}
