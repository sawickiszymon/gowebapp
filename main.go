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

	handler := handlers.NewPostHandler(s)

	router := httprouter.New()
	router.POST("/api/message", handler.PostMessage)
	router.GET("/api/message/:email", handler.ViewMessages)
	router.POST("/api/send", handlers.SendMessages(s))
	router.POST("/api/sendx", handler.Test)
	log.Fatal(http.ListenAndServe(":8080", router))

}
