package paxos

import "fmt"

// 消息类型枚举
type MessageType int

const (
	Prepare MessageType = iota + 1
	Propose
	Promise
	Accept
)

// 基础消息结构
type Message struct {
	Type   MessageType
	From   int
	To     int
	Number int
	Value  any
	Ok     bool
}

// GetFrom 实现 message 接口
func (m *Message) GetFrom() int {
	return m.From
}

// GetTo 实现 message 接口
func (m *Message) GetTo() int {
	return m.To
}

// GetType 实现 message 接口
func (m *Message) GetType() MessageType {
	return m.Type
}

// String 实现 Stringer 接口，让日志更易读
func (m *Message) String() string {
	typeStr := "Unknown"
	switch m.Type {
	case Prepare:
		typeStr = "Prepare"
	case Propose:
		typeStr = "Propose"
	case Promise:
		typeStr = "Promise"
	case Accept:
		typeStr = "Accept"
	}

	return fmt.Sprintf("Message{Type:%s From:%d To:%d Number:%d Value:%v Ok:%t}",
		typeStr, m.From, m.To, m.Number, m.Value, m.Ok)
}

// 提案参数
type Proposal struct {
	Number int
	Value  any
}

// 接受者状态
type AcceptorState struct {
	PromisedNumber int
	AcceptedNumber int
	AcceptedValue  any
}

// 提议者状态
type ProposerState struct {
	Round  int
	Number int
	Value  any
}
