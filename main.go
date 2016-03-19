package main

import "log"

func main() {
	server := NewSoundcloudServer(":9090", ":9099")
	log.Fatal(server.ListenAndServe())
}
