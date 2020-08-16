package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)



func main() {
	s := InitCluster()

	fmt.Println("Cassandra init done")
	fmt.Println("RestAPI server")

	router := httprouter.New()
	router.POST("/api/message", PostMessage(s))
	router.POST("/api/send", SendMessages(s))
	router.GET("/api/message/:email", ViewMessage(s))
	log.Fatal(http.ListenAndServe(":8080", router))

}

