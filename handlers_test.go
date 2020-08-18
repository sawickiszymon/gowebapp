// /handler_test.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sawickiszymon/gowebapp/driver"
	"github.com/sawickiszymon/gowebapp/handlers"
	"github.com/sawickiszymon/gowebapp/models"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)
func prepareEnvVar(){
	os.Setenv("CASSANDRA_URL", "cassandra")
	os.Setenv("CASSANDRA_KEYSPACE", "cass")
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

	body , _ := json.Marshal(postBody)
	request, err := http.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(response, request)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Email was saved: sz.sawicki1@gmail.com"
	if !cmp.Equal(response.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Body.String(), expected)
	}
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

	request, _ := http.NewRequest("GET", "/api/message/sz.sawicki1@gmail.com", nil)
	response := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("Wrong status")
	}

	var e []models.Email


	if err := json.Unmarshal([]byte(response.Body.String()), &e); err != nil {
		fmt.Println("ugh: ", err)
	}
	fmt.Println("E resposne: ", e)

	//expected := map[string]interface{}{
	//	"email":"sz.sawicki1@gmail.com", "title":"testTitle", "content":"testContent", "magic_number":25,
	//}
	//can be more emails
	expected := "sz.sawicki1@gmail.com"
	fmt.Println(cmp.Equal(e, expected))
	for el := range e {
		if  e[el].Email != expected{
			t.Errorf("handler returned unexpected body: got %v want %+v",
				response.Body.String(), expected)
		}
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