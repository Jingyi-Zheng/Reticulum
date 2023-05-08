package main

import "time"

type Node struct {
	ID            string
	Addr          string
	PSID          int
	CSID          int
	PSchain       PSchain
	CSchain       CSchain
	PSPeerID      []string
	CSPeerID      []string
	PSPeerAddr    []string
	CSPeerAddr    []string
	PSLeader      int
	CSLeader      int
	globalTimer2  *time.Timer
	timerExpired2 bool
	globalTimer1  *time.Timer
	timerExpired1 bool
	addps         bool
	T2            int
	Advepoch      []int
}
