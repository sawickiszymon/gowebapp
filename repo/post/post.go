package repo

import (
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/gocql/gocql"
	models "github.com/sawickiszymon/gowebapp/models"
	repo "github.com/sawickiszymon/gowebapp/repo"
	"log"
	"net/http"
	"reflect"
)

const (
	INSERT               = `INSERT INTO Email (email, title, content, magic_number) VALUES (?, ?, ?, ?) USING TTL 300`
	SELECT_EMAIL_TO_SEND = `SELECT email, title, content FROM Email WHERE magic_number = ?`
	SELECT_EMAIL         = `SELECT email, title, content, magic_number FROM Email WHERE email = ?`
	SELECT_COUNT         = `SELECT Count(*) FROM Email WHERE email = ?`
	DELETE_MESSAGE       = `DELETE FROM Email WHERE email = ? AND magic_number = ?`
)

var pageState []byte

func NewRepo(s *gocql.Session) repo.PostRepo {
	return &cassandraPostRepo{
		session: s,
	}
}

type cassandraPostRepo struct {
	session *gocql.Session
}

func (s *cassandraPostRepo) Create(e *models.Email) error {

	if isValid := PostRequestValidation(e); !isValid {
		return http.ErrBodyNotAllowed
	}

	err := checkmail.ValidateFormat(e.Email)
	if err != nil {
		log.Fatal(err)
	}

	PostEmail(e, s.session)

	return nil
}

func (s *cassandraPostRepo) View(e *models.Email, pageNumber int, email string) ([]models.Email, error) {


	var emailToDisplay []models.Email
	pageLimit := 4
	fmt.Println(email)


	var numberOfEmails = GetEmailCount(email, s.session)
	fmt.Println(numberOfEmails)
	var firstRowEmail = (pageNumber * pageLimit) - pageLimit
	fmt.Println(pageState)

	if err := s.session.Query(SELECT_EMAIL, email).PageState(pageState).Scan(&e.Email, &e.Title, &e.Content, &e.MagicNumber); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < pageNumber; i++ {

		if numberOfEmails <= firstRowEmail {
			return nil, http.ErrContentLength
		}

		iter := s.session.Query(SELECT_EMAIL, e.Email).PageState(pageState).PageSize(pageLimit).Iter()

		for iter.Scan(&e.Email, &e.Title, &e.Content, &e.MagicNumber) {
			if pageNumber%2 == 1 && i+1 == pageNumber {
				emailToDisplay = append(emailToDisplay, *e)
			} else if pageNumber%1 == 0 && i+1 == pageNumber {
				emailToDisplay = append(emailToDisplay, *e)
			}
			pageState = iter.PageState()
		}
	}
	pageState = nil

	return emailToDisplay, nil
}

func GetEmailCount(email string, s *gocql.Session) int {
	var count int
	iter := s.Query(SELECT_COUNT, email).Iter()
	for iter.Scan(&count) {
	}
	return count
}

func PostRequestValidation(e *models.Email) bool {
	isValid := true
	v := reflect.Indirect(reflect.ValueOf(e))

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		if value.IsZero() {
			isValid = false
		}
	}
	return isValid
}

func PostEmail(e *models.Email, session *gocql.Session) {
	if err := session.Query(INSERT, e.Email, e.Title, e.Content, e.MagicNumber).Exec(); err != nil {
		log.Println(err)
	}
}
