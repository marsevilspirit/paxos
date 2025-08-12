package paxos

import "time"

// proposer 提议者状态机
type proposer struct {
	id int

	// 提案值
	value any

	// 网络接口
	net network

	// 接受者ID列表
	acceptorIds []int

	// 当前轮次
	round int
}

// newProposer 创建新的提议者
func newProposer(id int, value any, net network, acceptorIds ...int) *proposer {
	return &proposer{
		id:          id,
		value:       value,
		net:         net,
		acceptorIds: acceptorIds,
		round:       0,
	}
}

// run 启动提议者
func (p *proposer) run() {
	p.round++
	number := p.proposalNumber()

	// Phase 1: Prepare
	prepareResponses := p.prepare(number)

	// 检查是否获得多数派承诺
	if len(prepareResponses) < p.majority() {
		return // 未获得多数派支持
	}

	// 找到最高编号的已接受值
	maxNumber := 0
	for _, resp := range prepareResponses {
		if resp.Ok && resp.Number > maxNumber {
			maxNumber = resp.Number
			p.value = resp.Value
		}
	}

	// Phase 2: Accept
	acceptResponses := p.accept(number)

	// 检查是否获得多数派接受
	if len(acceptResponses) < p.majority() {
		return // 未获得多数派接受
	}
}

// prepare 准备阶段
func (p *proposer) prepare(number int) []*Message {
	var responses []*Message

	// 发送Prepare消息给所有接受者
	for _, acceptorId := range p.acceptorIds {
		msg := &Message{
			Type:   Prepare,
			From:   p.id,
			To:     acceptorId,
			Number: number,
		}

		p.net.send(msg)
	}

	// 接收响应
	timeout := 100 * time.Millisecond
	for i := 0; i < len(p.acceptorIds); i++ {
		if response, ok := p.net.recv(timeout); ok {
			if msg, ok := response.(*Message); ok && msg.Type == Promise && msg.Ok {
				responses = append(responses, msg)
			}
		}
	}

	return responses
}

// accept 接受阶段
func (p *proposer) accept(number int) []*Message {
	var responses []*Message

	// 发送Accept消息给所有接受者
	for _, acceptorId := range p.acceptorIds {
		msg := &Message{
			Type:   Accept,
			From:   p.id,
			To:     acceptorId,
			Number: number,
			Value:  p.value,
		}

		p.net.send(msg)
	}

	// 接收响应
	timeout := 100 * time.Millisecond
	for i := 0; i < len(p.acceptorIds); i++ {
		if response, ok := p.net.recv(timeout); ok {
			if msg, ok := response.(*Message); ok && msg.Type == Accept && msg.Ok {
				responses = append(responses, msg)
			}
		}
	}

	return responses
}

// majority 计算多数派数量
func (p *proposer) majority() int {
	return len(p.acceptorIds)/2 + 1
}

// proposalNumber 生成提案编号
func (p *proposer) proposalNumber() int {
	return p.round<<16 | p.id
}
