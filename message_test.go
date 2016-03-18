package main

import "testing"

func TestMessageProperlyConvertsABroadcastString(t *testing.T) {
	msg, err := Parse("542532|B")
	if err != nil {
		t.Error("Got error:", err)
	}

	if msg.Timestamp != 542532 || msg.Type != "B" {
		t.Error("Broadcast message mismatch, got", msg)
	}
}

func TestMessageProperlyConvertsAFullMessageString(t *testing.T) {
	msg, err := Parse("43|P|32|56")
	if err != nil {
		t.Error("Got error:", err)
	}

	if msg.Timestamp != 43 || msg.Type != "P" || msg.ToID != "56" || msg.FromID != "32" {
		t.Error("Full message mismatch, got", msg)
	}
}

func TestMessageErrorsWhenInputIsTooShort(t *testing.T) {
	_, err := Parse("1234")
	if err == nil {
		t.Error("Did not get error")
	}
}

func TestMessageErrorsWhenInputIsTooLong(t *testing.T) {
	_, err := Parse("1234|23|23|23|23|23|23|23")
	if err == nil {
		t.Error("Did not get error")
	}
}

func TestMessageErrorsWhenTimestampIsNotAnInt(t *testing.T) {
	_, err := Parse("43A|P|32|56")
	if err == nil {
		t.Error("Did not get error")
	}
}

func TestMessageQueueSwapActuallySwaps(t *testing.T) {
	mq := MessageQueue([]Message{Message{Timestamp: 1234}, Message{Timestamp: 2222}})
	mq.Swap(0, 1)
	if mq[0].Timestamp != 2222 {
		t.Error("Expected 2222, got", mq[0].Timestamp)
	}
}

func TestMessageQueuePutsLowerNumbersFirst(t *testing.T) {
	q := MessageQueue{}

	q.Enqueue(Message{Timestamp: 222})
	q.Enqueue(Message{Timestamp: 111})
	q.Enqueue(Message{Timestamp: 333})

	tip := q.Peek()

	if tip.Timestamp != 111 {
		t.Error("Expected 111, got", tip.Timestamp)
	}
}

func TestMessageQueueDequeuesThingsCorrectly(t *testing.T) {
	q := MessageQueue{}

	q.Enqueue(Message{Timestamp: 222})
	q.Enqueue(Message{Timestamp: 111})
	q.Enqueue(Message{Timestamp: 444})
	q.Enqueue(Message{Timestamp: 333})

	values := []int{111, 222, 333, 444}
	for _, v := range values {
		tip := q.Dequeue()

		if tip.Timestamp != v {
			t.Error("Expected", v, "but got", tip.Timestamp)
		}
	}
}
