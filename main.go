package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sawickiszymon/gowebapp/driver"
	"github.com/sawickiszymon/gowebapp/handlers"
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
