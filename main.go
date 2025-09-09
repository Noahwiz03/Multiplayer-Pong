package main

import (
	"fmt"
	"log"
	"net/http"
	"pong/server"
)

func main() {

	http.HandleFunc("/ws", server.HandleWS)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	port := ":8080"
	fmt.Println("server listening on port " + port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
		log.Fatal(err)
	}
}
