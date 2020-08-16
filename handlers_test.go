// /handler_test.go
package main

import (
	"github.com/julienschmidt/httprouter"
	"gowebapp/driver"
	"gowebapp/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestViewMessage(t *testing.T) {

	//cfg := &gocql.ClusterConfig{}
	//handler := &gocql.Session{}
	handler := driver.InitCluster()
	defer handler.Close()
	router := httprouter.New()
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll pass 'nil' as the third parameter.
	//postBody := map[string]interface{}{
	//	"email":"sz.sawicki1@gmail.com", "title":"testTitle", "content":"testContent", "magic_number":25,
	//}
	//body , _ := json.Marshal(postBody)
	//fmt.Println(body)
	//request, err := http.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(body))
	router.GET("/api/message/:email", handlers.ViewMessage(handler))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//router.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	request, _ := http.NewRequest("GET", "/api/message/sz.sawicki1@gmail.com", nil)
	response := httptest.NewRecorder()
	querry := request.URL.Query()
	querry.Add("page", "1")
	request.URL.RawQuery = querry.Encode()
	//request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// returned := ViewMessage()
	//handlerTemp = ViewMessage()
	//httprouter.
	//handlery := http.HandlerFunc(httprouter.Handle(ViewMessage(handler))
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(response, request)

	// Check the status code is what we expect.
	if status := response.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
