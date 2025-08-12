package paxos

import (
	"testing"
	"time"
)

func TestPaxosWithSingleProposer(t *testing.T) {
	// 1, 2, 3 are acceptors
	// 1001 is a proposer
	// 2001 is a learner
	pn := newPaxosNetwork(1, 2, 3, 1001, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	p := newProposer(1001, "hello world", pn.agentNetwork(1001), 1, 2, 3)
	go p.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	if value != "hello world" {
		t.Errorf("value = %s, want %s", value, "hello world")
	}
}

func TestPaxosWithTwoProposers(t *testing.T) {
	// 1, 2, 3 are acceptors
	// 1001,1002 is a proposer
	pn := newPaxosNetwork(1, 2, 3, 1001, 1002, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	p1 := newProposer(1001, "hello world", pn.agentNetwork(1001), 1, 2, 3)
	go p1.run()

	time.Sleep(time.Millisecond)
	p2 := newProposer(1002, "bad day", pn.agentNetwork(1002), 1, 2, 3)
	go p2.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	if value != "hello world" {
		t.Errorf("value = %s, want %s", value, "hello world")
	}
	time.Sleep(time.Millisecond)
}

func TestPaxosWithThreeProposers(t *testing.T) {
	// 测试三个提议者并发提案
	pn := newPaxosNetwork(1, 2, 3, 1001, 1002, 1003, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	// 三个提议者几乎同时启动
	p1 := newProposer(1001, "first value", pn.agentNetwork(1001), 1, 2, 3)
	p2 := newProposer(1002, "second value", pn.agentNetwork(1002), 1, 2, 3)
	p3 := newProposer(1003, "third value", pn.agentNetwork(1003), 1, 2, 3)

	go p1.run()
	time.Sleep(10 * time.Microsecond)
	go p2.run()
	time.Sleep(10 * time.Microsecond)
	go p3.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()

	// 应该达成某个共识值
	if value == nil {
		t.Error("expected consensus value, got nil")
	}

	// 验证共识值是三个值中的一个
	validValues := []string{"first value", "second value", "third value"}
	found := false
	for _, v := range validValues {
		if value == v {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("consensus value %s is not one of the expected values", value)
	}
}

func TestPaxosWithFiveAcceptors(t *testing.T) {
	// 测试更多接受者的情况
	pn := newPaxosNetwork(1, 2, 3, 4, 5, 1001, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 5; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	p := newProposer(1001, "five acceptors test", pn.agentNetwork(1001), 1, 2, 3, 4, 5)
	go p.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3, 4, 5)
	value := l.learn()
	if value != "five acceptors test" {
		t.Errorf("value = %s, want %s", value, "five acceptors test")
	}
}

func TestPaxosWithMultipleLearners(t *testing.T) {
	// 测试多个学习者的情况
	pn := newPaxosNetwork(1, 2, 3, 1001, 2001, 2002, 2003)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001, 2002, 2003))
	}

	for _, a := range as {
		go a.run()
	}

	p := newProposer(1001, "multiple learners test", pn.agentNetwork(1001), 1, 2, 3)
	go p.run()

	// 创建三个学习者
	l1 := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	l2 := newLearner(2002, pn.agentNetwork(2002), 1, 2, 3)
	l3 := newLearner(2003, pn.agentNetwork(2003), 1, 2, 3)

	// 所有学习者应该学习到相同的值
	value1 := l1.learn()
	value2 := l2.learn()
	value3 := l3.learn()

	if value1 != value2 || value2 != value3 {
		t.Errorf("learners got different values: %s, %s, %s", value1, value2, value3)
	}

	if value1 != "multiple learners test" {
		t.Errorf("value = %s, want %s", value1, "multiple learners test")
	}
}

func TestPaxosWithDifferentValues(t *testing.T) {
	// 测试不同类型的值
	// 注意：在 Paxos 中，一旦达成共识，后续提案必须使用相同的值
	pn := newPaxosNetwork(1, 2, 3, 1001, 1002, 1003, 2001, 2002, 2003)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001, 2002, 2003))
	}

	for _, a := range as {
		go a.run()
	}

	// 测试字符串值
	p1 := newProposer(1001, "string value", pn.agentNetwork(1001), 1, 2, 3)
	go p1.run()

	l1 := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value1 := l1.learn()
	if value1 != "string value" {
		t.Errorf("string value = %s, want %s", value1, "string value")
	}

	// 测试数字值 - 由于 Paxos 特性，应该学习到第一个共识值
	p2 := newProposer(1002, 42, pn.agentNetwork(1002), 1, 2, 3)
	go p2.run()

	l2 := newLearner(2002, pn.agentNetwork(2002), 1, 2, 3)
	value2 := l2.learn()
	// 在 Paxos 中，一旦达成共识，后续提案必须使用相同的值
	if value2 != "string value" {
		t.Errorf("number value = %v, want %s (Paxos consensus)", value2, "string value")
	}

	// 测试结构体值 - 同样应该学习到第一个共识值
	type TestStruct struct {
		Name  string
		Value int
	}
	testStruct := TestStruct{Name: "test", Value: 100}

	p3 := newProposer(1003, testStruct, pn.agentNetwork(1003), 1, 2, 3)
	go p3.run()

	l3 := newLearner(2003, pn.agentNetwork(2003), 1, 2, 3)
	value3 := l3.learn()
	// 在 Paxos 中，一旦达成共识，后续提案必须使用相同的值
	if value3 != "string value" {
		t.Errorf("struct value = %v, want %s (Paxos consensus)", value3, "string value")
	}
}

func TestPaxosConcurrentProposals(t *testing.T) {
	// 测试大量并发提案
	pn := newPaxosNetwork(1, 2, 3, 1001, 1002, 1003, 1004, 1005, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	// 启动多个提议者
	proposers := []*proposer{
		newProposer(1001, "value1", pn.agentNetwork(1001), 1, 2, 3),
		newProposer(1002, "value2", pn.agentNetwork(1002), 1, 2, 3),
		newProposer(1003, "value3", pn.agentNetwork(1003), 1, 2, 3),
		newProposer(1004, "value4", pn.agentNetwork(1004), 1, 2, 3),
		newProposer(1005, "value5", pn.agentNetwork(1005), 1, 2, 3),
	}

	// 几乎同时启动所有提议者
	for _, p := range proposers {
		go p.run()
	}

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()

	// 应该达成某个共识值
	if value == nil {
		t.Error("expected consensus value, got nil")
	}

	// 验证共识值是预期值中的一个
	expectedValues := []string{"value1", "value2", "value3", "value4", "value5"}
	found := false
	for _, v := range expectedValues {
		if value == v {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("consensus value %s is not one of the expected values", value)
	}
}

func TestPaxosNetworkEmpty(t *testing.T) {
	// 测试网络为空的情况
	pn := newPaxosNetwork(1001)

	// 尝试发送消息到不存在的代理
	msg := &Message{
		Type:   Prepare,
		From:   1001,
		To:     999, // 不存在的代理
		Number: 1,
	}

	// 这应该不会导致panic
	pn.send(msg)

	// 检查网络是否为空
	if !pn.empty() {
		t.Error("network should be empty after sending to non-existent agent")
	}
}

func TestPaxosMessageInterface(t *testing.T) {
	// 测试消息接口实现
	msg := &Message{
		Type:   Prepare,
		From:   1001,
		To:     2001,
		Number: 42,
		Value:  "test",
		Ok:     true,
	}

	// 测试接口方法
	if msg.GetFrom() != 1001 {
		t.Errorf("GetFrom() = %d, want %d", msg.GetFrom(), 1001)
	}
	if msg.GetTo() != 2001 {
		t.Errorf("GetTo() = %d, want %d", msg.GetTo(), 2001)
	}
	if msg.GetType() != Prepare {
		t.Errorf("GetType() = %d, want %d", msg.GetType(), Prepare)
	}

	// 测试String方法
	str := msg.String()
	if str == "" {
		t.Error("String() method returned empty string")
	}
}
