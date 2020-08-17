// /handler_test.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sawickiszymon/gowebapp/driver"
	"github.com/sawickiszymon/gowebapp/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)
func prepareEnvVar(){
	os.Setenv("CASSANDRA_URL", "cassandra")
	os.Setenv("CASSANDRA_KEYSPACE", "cass")
}


func TestViewMessage(t *testing.T) {

	prepareEnvVar()
	s := driver.InitCluster()
	handler := handlers.NewPostHandler(s)
	defer s.Close()
	//router := httprouter.New()
	//// Create a request to pass to our handler. We don't have any query parameters for now, so we'll pass 'nil' as the third parameter.
	////postBody := map[string]interface{}{
	////	"email":"sz.sawicki1@gmail.com", "title":"testTitle", "content":"testContent", "magic_number":25,
	////}
	//req, err := http.NewRequest("GET", "/health-check", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	////body , _ := json.Marshal(postBody)
	////fmt.Println(body)
	////request, err := http.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(body))
	//router.GET("/api/message/:email", handler.ViewMessages)
	////if err != nil {
	////	t.Fatal(err)
	////}
	////router.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	//request, _ := http.NewRequest("GET", "/api/message/sz.sawicki1@gmail.com", nil)
	//response := httptest.NewRecorder()
	////request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//router.ServeHTTP(response, request)
	//fmt.Println("Request body: ", request.Body)
	//fmt.Println("Response body: ", response.Body.String())
	//
	//// Check the status code is what we expect.
	//if status := response.Code; status != http.StatusOK {
	//	t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	//}

	router := httprouter.New()
	router.GET("/api/message/:email", handler.ViewMessages)

	req, _ := http.NewRequest("GET", "/api/message/sz.sawicki1@gmail.com", nil)
	rr := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status")
	}

	expected := "Email was saved: sz.sawicki1@gmail.com"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPostValidMessage(t *testing.T) {
	prepareEnvVar()
	s := driver.InitCluster()
	handler := handlers.NewPostHandler(s)
	defer s.Close()


	router := httprouter.New()

	router.POST("/api/message", handler.PostMessage)

	postBody := map[string]interface{}{
		"email":"sz.sawicki1@gmail.com", "title":"testTitle", "content":"testContent", "magic_number":25,
	}
	fmt.Println(postBody)
	body , _ := json.Marshal(postBody)
	request, err := http.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(response, request)
	fmt.Println("Request body: ", request.Body)
	fmt.Println("Response body: ", response.Body.String())

	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestPostMessageWithoutMagicNumber(t *testing.T) {
	prepareEnvVar()
	s := driver.InitCluster()
	handler := handlers.NewPostHandler(s)
	defer s.Close()

	router := httprouter.New()

	router.POST("/api/message", handler.PostMessage)

	postBody := map[string]interface{}{
		"email":"sz.sawicki1@gmail.com", "title":"testTitle", "content":"testContent",
	}

	body , _ := json.Marshal(postBody)
	request, err := http.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(response, request)
	fmt.Println("Request body: ", request.Body)
	fmt.Println("Response body: ", response.Body)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}
}