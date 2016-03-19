package main

import "net"

// Notifier stores and notifies client connections about stuff
type Notifier interface {
	AddClient(UserClient)
	Notify(string, Message)
}

// TCPConnectionNotifier handles TCP connections
type TCPConnectionNotifier struct {
	clients map[string][]chan string
}

// NewTCPNotifier returns a new initialized notifier
func NewTCPNotifier() *TCPConnectionNotifier {
	return &TCPConnectionNotifier{clients: map[string][]chan string{}}
}

// Notify actually sends notifications
func (notifier *TCPConnectionNotifier) Notify(id string, msg Message) {
	// Iterate over a group of channels
	// If a channel's buffer is full, close it
	// Otherwise, send (non-blocking)
	channels, ok := notifier.clients[id]
	if ok {
		for index, ch := range channels {
			select {
			case ch <- msg.String():
				// sent the notification, loop
			default:
				// It's filled up its buffer. Crush, kill, destroy. Remove it from the list
				close(ch)
				notifier.clients[id] = append(notifier.clients[id][:index], notifier.clients[id][index+1:]...)
			}
		}
	}
}

// AddClient creates a new client and channel pair
func (notifier *TCPConnectionNotifier) AddClient(uc UserClient) {
	ch := make(chan string, 30)
	go processNotifications(uc.Conn, ch)
	notifier.clients[uc.ID] = append(notifier.clients[uc.ID], ch)
}

func processNotifications(conn net.Conn, notificationChannel <-chan string) {
	defer conn.Close()
	for notification := range notificationChannel {
		_, err := conn.Write([]byte(notification + "\n"))
		if err != nil {
			return
		}
	}
}
