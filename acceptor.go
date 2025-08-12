package paxos

import "time"

// acceptor 接受者状态机
type acceptor struct {
	id int

	// 状态
	state AcceptorState

	// 网络接口
	net network

	// 学习者ID列表
	learnerIds []int
}

// newAcceptor 创建新的接受者
func newAcceptor(id int, net network, learnerIds ...int) *acceptor {
	return &acceptor{
		id:         id,
		net:        net,
		learnerIds: learnerIds,
		state: AcceptorState{
			PromisedNumber: 0,
			AcceptedNumber: 0,
			AcceptedValue:  nil,
		},
	}
}

// HandlePrepare 处理准备请求
func (a *acceptor) HandlePrepare(msg *Message) *Message {
	if msg.Number > a.state.PromisedNumber {
		// 承诺接受更高编号的提案
		a.state.PromisedNumber = msg.Number

		return &Message{
			Type:   Promise,
			From:   a.id,
			To:     msg.From,
			Number: a.state.AcceptedNumber,
			Value:  a.state.AcceptedValue,
			Ok:     true,
		}
	}

	// 拒绝低编号的提案
	return &Message{
		Type: Promise,
		From: a.id,
		To:   msg.From,
		Ok:   false,
	}
}

// HandleAccept 处理接受请求
func (a *acceptor) HandleAccept(msg *Message) *Message {
	if msg.Number >= a.state.PromisedNumber {
		// 接受提案
		a.state.PromisedNumber = msg.Number
		a.state.AcceptedNumber = msg.Number
		a.state.AcceptedValue = msg.Value

		return &Message{
			Type: Accept,
			From: a.id,
			To:   msg.From,
			Ok:   true,
		}
	}

	// 拒绝提案
	return &Message{
		Type: Accept,
		From: a.id,
		To:   msg.From,
		Ok:   false,
	}
}

// GetState 获取当前状态
func (a *acceptor) GetState() AcceptorState {
	return a.state
}

// GetID 获取接受者ID
func (a *acceptor) GetID() int {
	return a.id
}

// run 启动接受者消息处理循环
func (a *acceptor) run() {
	for {
		if msg, ok := a.net.recv(100 * time.Millisecond); ok {
			if message, ok := msg.(*Message); ok {
				switch message.GetType() {
				case Prepare:
					response := a.HandlePrepare(message)
					a.net.send(response)
				case Accept:
					response := a.HandleAccept(message)
					// 通知所有学习者
					for _, learnerId := range a.learnerIds {
						learnMsg := &Message{
							Type:   Accept,
							From:   a.id,
							To:     learnerId,
							Number: message.Number,
							Value:  message.Value,
							Ok:     true,
						}
						a.net.send(learnMsg)
					}
					a.net.send(response)
				}
			}
		}
	}
}
