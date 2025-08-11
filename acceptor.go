package paxos

import "net"

type Acceptor struct {
	lis net.Listener

	// server id
	id int

	minProposal int

	acceptedNumber int

	acceptedValue any

	learners []int
}
