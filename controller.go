package main

import (
	"log"
	"net"
)

// Controller is an interface representing a controller that can handle incoming events
// and add new clients
type Controller interface {
	HandleEvent(string) error
	Reset()
	AddUserClient(UserClient)
}

// ForwarderController connects the domain model and the notification logic.
// This is the class that receives the message, translates it into the domain,
// and retrieves the response from the domain and notifies the clients
type ForwarderController struct {
	Users    UserCollection
	Notifier Notifier

	pendingMessages MessageQueue
	lastDeliveredID int
}

// NewForwarderController creates a controller with the default
func NewForwarderController() *ForwarderController {
	return &ForwarderController{
		Users:           NewUserCollection(),
		Notifier:        NewTCPNotifier(),
		pendingMessages: MessageQueue{},
		lastDeliveredID: 0}
}

// UserClient represents an incoming UserClient
type UserClient struct {
	Conn net.Conn
	ID   string
}

// HandleEvent handles an incoming message from the pipe
func (fw *ForwarderController) HandleEvent(evt string) error {
	msg, err := Parse(evt)
	if err != nil {
		return err
	}
	fw.ProcessMessage(msg)

	return nil
}

// Reset resets the queue and message id
func (fw *ForwarderController) Reset() {
	fw.lastDeliveredID = 0
	fw.pendingMessages = MessageQueue{}
}

// AddUserClient creates the user in our data store and passes the ID to the notifier
func (fw *ForwarderController) AddUserClient(uc UserClient) {
	id := uc.ID
	fw.Users.GetOrCreateUser(id)
	fw.Notifier.AddClient(uc)
}

// ProcessMessage handles a new message
func (fw *ForwarderController) ProcessMessage(msg Message) {
	fw.pendingMessages.Enqueue(msg)
	fw.flushPendingQueue()
}

// IsNext determines whether a message is the next message to be handled in order
func (fw *ForwarderController) isNext(msg Message) bool {
	return msg.Timestamp == fw.lastDeliveredID+1
}

// HandleMessage handles an individual message. It passes it to the domain, then
// asks the notifier to notify all client connections about the
func (fw *ForwarderController) dispatchMessage(msg Message) {
	fw.lastDeliveredID = msg.Timestamp
	users, err := fw.Users.UpdateAndGetNotifiees(msg)
	if err != nil {
		log.Println(err)
		return
	}

	for _, u := range users {
		fw.Notifier.Notify(u.ID, msg)
	}
}

// FlushPendingQueue handles any messages in the queue it possibly can
func (fw *ForwarderController) flushPendingQueue() {
	for fw.pendingMessages.Any() && fw.isNext(fw.pendingMessages.Peek()) {
		msg := fw.pendingMessages.Dequeue()
		fw.dispatchMessage(msg)
	}
}
