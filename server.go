package main

import (
	"bufio"
	"log"
	"net"
)

// SoundcloudServer holds information and behavior for the forwarder
type SoundcloudServer struct {
	EventHoseBindingAddr  string
	UserClientBindingAddr string
	Controller            Controller
}

// ListenAndServe starts listening
func (scs *SoundcloudServer) ListenAndServe() error {
	eventListener, err := scs.listenEventHose()
	if err != nil {
		return err
	}

	userConnection, err := scs.listenUserChannel()
	if err != nil {
		return err
	}

	for {
		select {
		case conn := <-userConnection:
			scs.Controller.AddUserClient(conn)
		case eventString := <-eventListener:
			err := scs.Controller.HandleEvent(eventString)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (scs *SoundcloudServer) listenUserChannel() (chan UserClient, error) {
	listener, err := net.Listen("tcp", scs.UserClientBindingAddr)
	ch := make(chan UserClient)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			defer listener.Close()
			conn, err := listener.Accept()

			if err != nil {
				log.Println("Unable to accept user client connection due to error:", err)
				continue
			}
			go func(c net.Conn) {
				scanner := bufio.NewScanner(c)
				scanner.Scan()

				ch <- UserClient{ID: scanner.Text(), Conn: c}
			}(conn)
		}
	}()

	return ch, nil
}

func (scs *SoundcloudServer) listenEventHose() (chan string, error) {
	ch := make(chan string)
	listener, err := net.Listen("tcp", scs.EventHoseBindingAddr)
	if err != nil {
		return nil, err
	}
	go func() {
		defer listener.Close()

		// Loop back up to the top if the event source disconnects and listen for another
		for {
			// Accept the first connection...
			c, err := listener.Accept()
			if err != nil {
				log.Println("Unable to accept event hose connection due to error:", err)
				continue
			}

			// Read lines out of it until the socket closes.
			scanner := bufio.NewScanner(c)
			for scanner.Scan() {
				ch <- scanner.Text()
			}

			log.Println("Event hose connection failed. Restarting")
		}
	}()
	return ch, nil
}
