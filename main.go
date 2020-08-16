package main

import (
	"fmt"
	"gowebapp/driver"
	"gowebapp/handlers"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	s := driver.InitCluster()

	fmt.Println("Cassandra init done")
	fmt.Println("RestAPI server")

	router := httprouter.New()
	router.POST("/api/message", handlers.PostMessage(s))
	router.POST("/api/send", handlers.SendMessages(s))
	router.GET("/api/message/:email", handlers.ViewMessage(s))
	log.Fatal(http.ListenAndServe(":8080", router))

}
