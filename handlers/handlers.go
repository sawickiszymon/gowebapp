package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"os"
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

func (p *Post) Test(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	e := DecodeRequest(writer, request)
	if err := p.repo.SendEmails(e.Email); err != nil {
		json.NewEncoder(writer).Encode(err)
		return
	}
}


func (p *Post) PostMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	e := DecodeRequest(writer, request)
	if err := p.repo.Create(&e); err != nil {
		json.NewEncoder(writer).Encode(err)
		return
	}
}

func SendMessages(s *gocql.Session) func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	var emails []models.Email
	return func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		e := DecodeRequest(writer, request)
		iter := s.Query(SELECT_EMAIL_TO_SEND,
			e.MagicNumber).Iter()

		for iter.Scan(&e.Email, &e.Title, &e.Content) {
			emails = append(emails, e)
		}
		SendEmails(emails)
		for el := range emails {
			if err := s.Query(DELETE_MESSAGE,
				emails[el].Email, e.MagicNumber).Exec(); err != nil {
				log.Fatal(err)
			}
		}
		emails = nil
	}
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

func SendEmails(e []models.Email) {

	s := NewSmtpConfig()
	addr := s.SmtpAddress + s.SmtpPort
	auth := smtp.PlainAuth(" ", s.SmtpEmail, s.SmtpPass, s.SmtpAddress)

	for _, elem := range e {

		msg := []byte("To:" + elem.Email + "\r\n" +
			"Subject:" + elem.Title + "\r\n" +
			"\r\n" +
			elem.Content + "\r\n")
		to := []string{elem.Email}
		err := smtp.SendMail(addr, auth, s.SmtpEmail, to, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DecodeRequest(w http.ResponseWriter, r *http.Request) models.Email {
	var e models.Email
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return e
}

func NewSmtpConfig() *models.SmtpConfig {
	return &models.SmtpConfig{
		SmtpAddress: os.Getenv("SMTP_SERV"),
		SmtpPort:    os.Getenv("SMTP_PORT"),
		SmtpEmail:   os.Getenv("FROM"),
		SmtpPass:    os.Getenv("PASS"),
	}
}
