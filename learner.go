package paxos

import "time"

// learner 学习者状态机
type learner struct {
	id int

	// 记录每个接受者接受的消息
	acceptedMessages map[int]*Message

	// 网络接口
	net network

	// 接受者ID列表
	acceptorIds []int
}

// newLearner 创建新的学习者
func newLearner(id int, net network, acceptorIds ...int) *learner {
	learner := &learner{
		id:               id,
		net:              net,
		acceptorIds:      acceptorIds,
		acceptedMessages: make(map[int]*Message),
	}

	// 初始化每个接受者的消息记录
	for _, aid := range acceptorIds {
		learner.acceptedMessages[aid] = &Message{
			Number: 0,
			Value:  nil,
		}
	}

	return learner
}

// learn 学习并返回共识值
func (l *learner) learn() any {
	// 等待并接收来自接受者的消息
	timeout := 100 * time.Millisecond
	for {
		if msg, ok := l.net.recv(timeout); ok {
			if message, ok := msg.(*Message); ok {
				l.Learn(message)

				// 检查是否达成共识
				if chosen := l.Chosen(); chosen != nil {
					return chosen
				}
			}
		}
	}
}

// Learn 学习接受者接受的消息
func (l *learner) Learn(msg *Message) bool {
	// 处理来自接受者的消息
	current, exists := l.acceptedMessages[msg.From]
	if !exists {
		// 如果来源不存在，初始化一个空消息
		current = &Message{
			Number: 0,
			Value:  nil,
		}
	}

	if current.Number < msg.Number {
		// 更新为更高编号的消息
		l.acceptedMessages[msg.From] = &Message{
			Number: msg.Number,
			Value:  msg.Value,
		}
		return true
	}
	return false
}

// Chosen 判断是否已经选择了值
func (l *learner) Chosen() any {
	// 统计每个编号被接受的次数
	acceptCounts := make(map[int]int)
	acceptMessages := make(map[int]*Message)

	for _, msg := range l.acceptedMessages {
		if msg.Number != 0 {
			acceptCounts[msg.Number]++
			acceptMessages[msg.Number] = msg
		}
	}

	// 检查是否有编号获得了多数派支持
	for number, count := range acceptCounts {
		if count >= l.majority() {
			return acceptMessages[number].Value
		}
	}
	return nil
}

// majority 计算多数派数量
func (l *learner) majority() int {
	return len(l.acceptorIds)/2 + 1
}
