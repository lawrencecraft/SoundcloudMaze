package main

import (
	"container/heap"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Message is an individual event
type Message struct {
	Timestamp int
	Type      string
	FromID    string
	ToID      string
    
    originalMessage string
}

const messageDelimiter = "|"

// Parse converts a string into a message
func Parse(str string) (Message, error) {
	var fields [4]string
	split := strings.Split(str, messageDelimiter)

	if len(split) > 4 || len(split) < 2 {
		return Message{}, errors.New(fmt.Sprint("Unable to parse message: got", len(split), "fields"))
	}

	for i := range split {
		fields[i] = split[i]
	}

	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return Message{}, err
	}
	return Message{Timestamp: id, Type: fields[1], FromID: fields[2], ToID: fields[3], originalMessage:str}, nil
}

// String converts a message to a string
func (m *Message) String() string {
	return m.originalMessage
}

// MessageQueue is a PQ of messages
type MessageQueue []Message

func (mq MessageQueue) Less(i, j int) bool {
	return mq[i].Timestamp < mq[j].Timestamp
}

func (mq MessageQueue) Swap(i, j int) {
	mq[i], mq[j] = mq[j], mq[i]
}

func (mq MessageQueue) Len() int {
	return len(mq)
}

func (mq *MessageQueue) Push(x interface{}) {
	*mq = append(*mq, x.(Message))
}

func (mq *MessageQueue) Pop() interface{} {
	index := len(*mq) - 1
	item := (*mq)[index]
	*mq = (*mq)[0:index]
	return item
}

func (mq *MessageQueue) Peek() Message {
	return (*mq)[0]
}

func (mq *MessageQueue) Enqueue(msg Message) {
	heap.Push(mq, msg)
}

func (mq *MessageQueue) Dequeue() Message {
	return heap.Pop(mq).(Message)
}

func (mq *MessageQueue) Any() bool {
	return len(*mq) > 0
}
