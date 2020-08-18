// /handler_test.go
package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/julienschmidt/httprouter"
	"github.com/sawickiszymon/gowebapp/driver"
	"github.com/sawickiszymon/gowebapp/handlers"
	"github.com/sawickiszymon/gowebapp/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestViewZeroMessages(t *testing.T) {

	prepareEnvVar()
	s := driver.InitCluster()
	handler := handlers.NewPostHandler(s)
	defer s.Close()

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
		t.Error("Error when decoding response: ", err)
	}

	//Check every email
	expected := "sz.sawicki1@gmail.com"
	for el := range e {
		if  e[el].Email != expected{
			t.Errorf("handler returned unexpected body: got %v want %+v",
				response.Body.String(), expected)
		}
	}

}

func TestPostValidMessage(t *testing.T) {

	var postBody = models.Email{
		Email:       "sz.sawicki1@gmail.com",
		Title:       "test",
		Content:     "test",
		MagicNumber: 11,
	}

	response := ExecutePostQuery(postBody)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	readLine := response.Body.String()[:len(response.Body.String())-1]
	readLine = strings.TrimSuffix(readLine, "\n")

	expected := "\"" + "Email was saved: " + "sz.sawicki1@gmail.com" + "\""
	if !cmp.Equal(readLine, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			readLine, expected)
	}
}

func TestViewMessage(t *testing.T) {

	var postBody = models.Email{
		Email:       "sz.sawicki1@gmail.com",
		Title:       "testTitle",
		Content:     "testContent",
		MagicNumber: 11,
	}

	_ = ExecutePostQuery(postBody)

	prepareEnvVar()
	s := driver.InitCluster()
	handler := handlers.NewPostHandler(s)
	defer s.Close()

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
		t.Error("Error when decoding response: ", err)
	}

	//Check every email
	expected := "sz.sawicki1@gmail.com"
	for el := range e {
		if  e[el].Email != expected{
			t.Errorf("handler returned unexpected body: got %v want %+v",
				response.Body.String(), expected)
		}
	}

}


func TestPostMessageWithoutMagicNumber(t *testing.T) {

	var postBody = models.Email{
		Email:       "sz.sawicki1@gmail.com",
		Title:       "test",
		Content:     "test",
	}

	response := ExecutePostQuery(postBody)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	readLine := response.Body.String()[:len(response.Body.String())-1]
	readLine = strings.TrimSuffix(readLine, "\n")

	expected := "\""+http.ErrBodyNotAllowed.Error()+"\""

	if !cmp.Equal(readLine, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			readLine, expected)
	}
}

func TestPostMessageWithoutContent(t *testing.T) {

	var postBody = models.Email{
		Email:       "sz.sawicki1@gmail.com",
		Title:       "test",
		MagicNumber: 11,
	}

	response := ExecutePostQuery(postBody)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	readLine := response.Body.String()[:len(response.Body.String())-1]
	readLine = strings.TrimSuffix(readLine, "\n")

	expected := "\""+http.ErrBodyNotAllowed.Error()+"\""

	if readLine != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			readLine, expected)
	}
}

func TestPostMessageWithWrongEmailAddress(t *testing.T) {

	var postBody = models.Email{
		Email:       "sz.sawicki1gmail.com",
		Title:       "test",
		Content:     "test",
		MagicNumber: 11,
	}

	response := ExecutePostQuery(postBody)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	readLine := response.Body.String()[:len(response.Body.String())-1]
	readLine = strings.TrimSuffix(readLine, "\n")

	expected := "\""+http.ErrBodyNotAllowed.Error()+"\""

	if readLine != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			readLine, expected)
	}
}

func ExecutePostQuery(postBody models.Email) *httptest.ResponseRecorder {
	prepareEnvVar()
	s := driver.InitCluster()
	handler := handlers.NewPostHandler(s)
	defer s.Close()

	router := httprouter.New()
	router.POST("/api/message", handler.PostMessage)

	body , _ := json.Marshal(postBody)
	request, _ := http.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(body))
	response := httptest.NewRecorder()
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(response, request)

	return response
}

func prepareEnvVar(){
	os.Setenv("CASSANDRA_URL", "cassandra")
	os.Setenv("CASSANDRA_KEYSPACE", "cass")
}

