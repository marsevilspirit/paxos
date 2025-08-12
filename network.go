package paxos

import (
	"log"
	"time"
)

// message 消息接口
type message interface {
	GetFrom() int
	GetTo() int
	GetType() MessageType
}

// network 网络接口
type network interface {
	send(m message)
	recv(timeout time.Duration) (message, bool)
}

// paxosNetwork Paxos网络实现
type paxosNetwork struct {
	recvQueues map[int]chan message
}

// newPaxosNetwork 创建新的Paxos网络
func newPaxosNetwork(agents ...int) *paxosNetwork {
	pn := &paxosNetwork{
		recvQueues: make(map[int]chan message, 0),
	}

	for _, a := range agents {
		pn.recvQueues[a] = make(chan message, 1024)
	}
	return pn
}

// agentNetwork 代理网络
func (pn *paxosNetwork) agentNetwork(id int) *agentNetwork {
	return &agentNetwork{id: id, paxosNetwork: pn}
}

// send 发送消息
func (pn *paxosNetwork) send(m message) {
	log.Printf("nt: send %+v", m)
	if queue, exists := pn.recvQueues[m.GetTo()]; exists {
		queue <- m
	} else {
		// 如果目标代理不存在，忽略消息
		log.Printf("nt: target agent %d not found, ignoring message", m.GetTo())
	}
}

// empty 检查网络是否为空
func (pn *paxosNetwork) empty() bool {
	var n int
	for i, q := range pn.recvQueues {
		log.Printf("nt: %d left %d", i, len(q))
		n += len(q)
	}
	return n == 0
}

// recvFrom 从指定代理接收消息
func (pn *paxosNetwork) recvFrom(from int, timeout time.Duration) (message, bool) {
	select {
	case m := <-pn.recvQueues[from]:
		log.Printf("nt: recv %+v", m)
		return m, true
	case <-time.After(timeout):
		return nil, false
	}
}

// agentNetwork 代理网络
type agentNetwork struct {
	id int
	*paxosNetwork
}

// send 发送消息
func (an *agentNetwork) send(m message) {
	an.paxosNetwork.send(m)
}

// recv 接收消息
func (an *agentNetwork) recv(timeout time.Duration) (message, bool) {
	return an.recvFrom(an.id, timeout)
}
