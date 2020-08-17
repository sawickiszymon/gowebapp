// /handler_test.go
package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sawickiszymon/gowebapp/driver"
	"github.com/sawickiszymon/gowebapp/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestViewMessage(t *testing.T) {

	fmt.Println("yes")
	os.Setenv("CASSANDRA_URL", "cassandra")
	os.Setenv("CASSANDRA_KEYSPACE", "cass")
	//cfg := &gocql.ClusterConfig{}
	fmt.Println("Docker test started somehow")
	//handler := &gocql.Session{}
	s := driver.InitCluster()
	handler := handlers.NewPostHandler(s)
	//defer s.Close()
	router := httprouter.New()
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll pass 'nil' as the third parameter.
	//postBody := map[string]interface{}{
	//	"email":"sz.sawicki1@gmail.com", "title":"testTitle", "content":"testContent", "magic_number":25,
	//}
	//body , _ := json.Marshal(postBody)
	//fmt.Println(body)
	//request, err := http.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(body))
	router.GET("/api/message/:email", handler.ViewMessages)
	//if err != nil {
	//	t.Fatal(err)
	//}"github.com/jwilder/dockerize"
	//router.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	request, _ := http.NewRequest("GET", "/api/message/sz.sawicki1@gmail.com", nil)
	response := httptest.NewRecorder()

	querry := request.URL.Query()
	querry.Add("page", "1")
	request.URL.RawQuery = querry.Encode()
	//request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(response, request)
	fmt.Println(response.Body)

	// Check the status code is what we expect.
	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}
}