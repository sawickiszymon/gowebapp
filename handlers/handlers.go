package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	models "github.com/sawickiszymon/gowebapp/models"
	repository "github.com/sawickiszymon/gowebapp/repo"
	post "github.com/sawickiszymon/gowebapp/repo/post"
)

const (
	SELECT_EMAIL_TO_SEND = `SELECT email, title, content FROM Email WHERE magic_number = ?`
	DELETE_MESSAGE       = `DELETE FROM Email WHERE email = ? AND magic_number = ?`
)

func NewPostHandler(s *gocql.Session) *Post {
	return &Post{
		repo: post.NewRepo(s),
	}
}

type Post struct {
	repo repository.PostRepo
}

func (p *Post) SendMessages(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	e := DecodeRequest(writer, request)
	err, emails:= p.repo.SendEmails(e.MagicNumber)
	if err != nil  {
		json.NewEncoder(writer).Encode(err)
		return
	}
	fmt.Println(e.Email)

	json.NewEncoder(writer).Encode("Emails were sent: " + e.Email)
}


func (p *Post) PostMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	e := DecodeRequest(writer, request)
	if err := p.repo.Create(&e); err != nil {
		json.NewEncoder(writer).Encode(err.Error())
		return
	}
	json.NewEncoder(writer).Encode("Email was saved: " + e.Email)
}

func (p *Post) ViewMessages(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {

	var pageNumber int
	pages, _ := request.URL.Query()["page"]
	//If page not specified return first page, else return specified page
	if len(pages) < 1 {
		pageNumber = 1
	} else {
		key := pages[0]
		pageNumber, _ = strconv.Atoi(key)
	}

	emailToDisplay, err := p.repo.ViewMessages(pageNumber, ps.ByName("email"))
	if err != nil {
		json.NewEncoder(writer).Encode(err)
		return
	}
	json.NewEncoder(writer).Encode(&emailToDisplay)
	emailToDisplay = nil
}

func DecodeRequest(w http.ResponseWriter, r *http.Request) models.Email {
	var e models.Email
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return e
}
