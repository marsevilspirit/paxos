package paxos

type Proposer struct {
	// server id
	id int

	round int

	// propose number
	number int

	acceptors []int
}
