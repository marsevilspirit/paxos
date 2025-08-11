package paxos

import "net"

type learner struct {
	lis net.Listener

	// learner id
	id int

	// int -> acceptor id
	acceptedMsg map[int]MsgArgs
}
