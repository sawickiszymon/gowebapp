package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	models "github.com/sawickiszymon/gowebapp/models"
	repository "github.com/sawickiszymon/gowebapp/repo"
	post "github.com/sawickiszymon/gowebapp/repo/post"
)

// Constructor initializing Cassandra instance
func NewPostHandler(s *gocql.Session) *Post {
	return &Post{
		repo: post.NewRepo(s),
	}
}

// Repository struct
type Post struct {
	repo repository.PostRepo
}

// SendMessages handles POST /api/send endpoint
// Sends messages with provided Email.magic_number value from request body and then deletes them from database
func (p *Post) SendMessages(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	e := DecodeRequest(writer, request)
	err, emails := p.repo.SendEmails(e.MagicNumber)
	if err != nil  {
		json.NewEncoder(writer).Encode(err)
		return
	}

	json.NewEncoder(writer).Encode("Emails were sent: " + strings.Join(emails, ","))
}

// PostMessage handles POST /api/message endpoint
// Saves messages specified in request body to database
// Database records lifes period is 5 minutes
func (p *Post) PostMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	e := DecodeRequest(writer, request)
	if err := p.repo.Create(&e); err != nil {
		json.NewEncoder(writer).Encode(err.Error())
		return
	}
	json.NewEncoder(writer).Encode("Email was saved: " + e.Email)
}

// ViewMessages handles GET /api/message/{emailValue} endpoint
// Returns emails with given email value in batches of 4
// If page parameter not specified return first page, else return provided page
func (p *Post) ViewMessages(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {

	var pageNumber int
	pages, _ := request.URL.Query()["page"]

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

// DecodeRequest decodes request body, and returns it in form of Email struct
func DecodeRequest(w http.ResponseWriter, r *http.Request) models.Email {
	var e models.Email
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return e
}
