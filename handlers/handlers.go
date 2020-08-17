package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	//"reflect"
	"strconv"
	//
	//"github.com/badoux/checkmail"
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	models "github.com/sawickiszymon/gowebapp/models"
	repository "github.com/sawickiszymon/gowebapp/repo"
	post "github.com/sawickiszymon/gowebapp/repo/post"
)

const (
	INSERT               = `INSERT INTO Email (email, title, content, magic_number) VALUES (?, ?, ?, ?) USING TTL 300`
	SELECT_EMAIL_TO_SEND = `SELECT email, title, content FROM Email WHERE magic_number = ?`
	SELECT_EMAIL         = `SELECT email, title, content, magic_number FROM Email WHERE email = ?`
	SELECT_COUNT         = `SELECT Count(*) FROM Email WHERE email = ?`
	DELETE_MESSAGE       = `DELETE FROM Email WHERE email = ? AND magic_number = ?`
)

var pageState []byte

func NewPostHandler(s *gocql.Session) *Post {
	return &Post{
		repo: post.NewRepo(s),
	}
}


type Post struct {
	repo repository.PostRepo
}

func (p *Post) PostMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params){

	e := DecodeRequest(writer, request)
	if err := p.repo.Create(&e); err != nil {
			json.NewEncoder(writer).Encode(err)
			return
	}
}

//func PostMessage(s *gocql.Session) func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//		e := DecodeRequest(writer, request)
//		fmt.Println(e)
//
//
//		//if isValid := repo.PostRequestValidation(e); !isValid {
//		//	json.NewEncoder(writer).Encode("Fill all fields!")
//		//	return
//		//}
//
//		//err := checkmail.ValidateFormat(e.Email)
//		//if err != nil {
//		//	log.Fatal(err)
//		//}
//		//PostEmail(&e, s)
//	}
//}

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

func (p *Post) Test(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	var e *models.Email
	var pageNumber int
	pages, _ := request.URL.Query()["page"]
	//If page not specified return first page else return specified page
	if len(pages) < 1 {
		pageNumber = 1
	} else {
		key := pages[0]
		pageNumber, _ = strconv.Atoi(key)
	}
	fmt.Println(ps.ByName("email"))
	emailToDisplay, err := p.repo.View(e, pageNumber, ps.ByName("email"))
	if err != nil {
		json.NewEncoder(writer).Encode(err)
	}
	json.NewEncoder(writer).Encode(&emailToDisplay)
	emailToDisplay = nil

}

func ViewMessage(s *gocql.Session) func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	return func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		var e models.Email
		var emailToDisplay []models.Email
		var pageNumber int
		pages, _ := request.URL.Query()["page"]
		pageLimit := 4
		//If page not specified return first page else return specified page
		if len(pages) < 1 {
			pageNumber = 1
		} else {
			key := pages[0]
			pageNumber, _ = strconv.Atoi(key)
		}

		var numberOfEmails = GetEmailCount(ps.ByName("email"), s)
		var firstRowEmail = (pageNumber * pageLimit) - pageLimit

		if err := s.Query(SELECT_EMAIL, ps.ByName( "email")).PageState(pageState).Scan(&e.Email, &e.Title, &e.Content, &e.MagicNumber); err != nil {
			log.Println(err)
		}

		for i := 0; i < pageNumber; i++ {

			if numberOfEmails <= firstRowEmail {
				json.NewEncoder(writer).Encode("There is no emails to display")
				return
			}

			iter := s.Query(SELECT_EMAIL, e.Email).PageState(pageState).PageSize(pageLimit).Iter()

			for iter.Scan(&e.Email, &e.Title, &e.Content, &e.MagicNumber) {
				if pageNumber%2 == 1 && i+1 == pageNumber {
					emailToDisplay = append(emailToDisplay, e)
				} else if pageNumber%1 == 0 && i+1 == pageNumber {
					emailToDisplay = append(emailToDisplay, e)
				}
				pageState = iter.PageState()
			}
		}
		json.NewEncoder(writer).Encode(&emailToDisplay)
		emailToDisplay = nil
		pageState = nil
	}
}
func GetEmailCount(email string, s *gocql.Session) int {
	var count int
	iter := s.Query(SELECT_COUNT, email).Iter()
	for iter.Scan(&count) {
	}
	return count
}

//func PostEmail(e *models.Email, session *gocql.Session) {
//	if err := session.Query(INSERT, e.Email, e.Title, e.Content, e.MagicNumber).Exec(); err != nil {
//		log.Println(err)
//	}
//}


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

//func PostRequestValidation(e models.Email) bool {
//	isValid := true
//	v := reflect.ValueOf(e)
//	for i := 0; i < v.NumField(); i++ {
//		value := v.Field(i)
//		if value.IsZero() {
//			isValid = false
//		}
//	}
//	return isValid
//}

func NewSmtpConfig() *models.SmtpConfig {
	return &models.SmtpConfig{
		SmtpAddress: os.Getenv("SMTP_SERV"),
		SmtpPort:    os.Getenv("SMTP_PORT"),
		SmtpEmail:   os.Getenv("FROM"),
		SmtpPass:    os.Getenv("PASS"),
	}
}
