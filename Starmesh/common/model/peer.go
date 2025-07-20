package model

import "time"

type Peer struct {
	IP       string    `json:"ip"`
	Port     string    `json:"port"`
	Lat      float64   `json:"lat"`
	Lon      float64   `json:"lon"`
	LastSeen time.Time `json:"lastSeen"`
}
