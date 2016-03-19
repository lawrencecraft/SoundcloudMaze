package main

import "testing"

type TestNotifier struct {
	Messages []Message
}

func (t *TestNotifier) AddClient(u UserClient) {}
func (t *TestNotifier) Notify(client string, m Message) {
	t.Messages = append(t.Messages, m)
}

func TestForwarderControllerNotifiesAMessage(t *testing.T) {
	f := NewForwarderController()
	notifier := TestNotifier{}
	f.Notifier = &notifier

	f.ProcessMessage(Message{Type: "F", FromID: "001", ToID: "003", Timestamp: 1})

	if len(notifier.Messages) != 1 {
		t.Error("Expected 1 message but got", len(notifier.Messages))
	}
}

func TestForwarderControllerDoesNotDeliverAMessageOutOfOrder(t *testing.T) {
	f := NewForwarderController()
	notifier := TestNotifier{}
	f.Notifier = &notifier

	f.ProcessMessage(Message{Type: "F", FromID: "001", ToID: "003", Timestamp: 2})

	if len(notifier.Messages) != 0 {
		t.Error("Should not have sent the message yet")
	}

	f.ProcessMessage(Message{Type: "F", FromID: "003", ToID: "004", Timestamp: 1})
	if len(notifier.Messages) != 2 {

		t.Error("Expected 2 messages but got", len(notifier.Messages))
	}
}
