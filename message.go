package paxos

import "net/rpc"

type MsgArgs struct {
	// proposal id
	Number int

	// proposal value
	Value any

	// sender id
	From int

	// recver id
	To int
}

type MsgReply struct {
	Ok bool

	Number int

	Value any
}

func call(srv string, name string, args any, reply any) bool {
	c, err := rpc.Dial("tcp", srv)
	if err != nil {
		return false
	}
	defer c.Close()

	err = c.Call(name, args, reply)
	if err != nil {
		return false
	}
	return true
}
