package repo

import (
	"github.com/badoux/checkmail"
	"github.com/gocql/gocql"
	"github.com/sawickiszymon/gowebapp/models"
	"github.com/sawickiszymon/gowebapp/repo"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"reflect"
)

// Read only queries
const (
	INSERT               = `INSERT INTO Email (email, title, content, magic_number) VALUES (?, ?, ?, ?) USING TTL 300`
	SELECT_EMAIL_TO_SEND = `SELECT email, title, content FROM Email WHERE magic_number = ?`
	SELECT_EMAIL         = `SELECT email, title, content, magic_number FROM Email WHERE email = ?`
	SELECT_COUNT         = `SELECT Count(*) FROM Email WHERE email = ?`
	DELETE_MESSAGE       = `DELETE FROM Email WHERE email = ? AND magic_number = ?`
)

var pageState []byte

// Constructor of PostRepo initializing Cassandra instance
func NewRepo(s *gocql.Session) repo.PostRepo {
	return &cassandraPostRepo{
		session: s,
	}
}

// Struct of cassandra session
type cassandraPostRepo struct {
	session *gocql.Session
}

// Create saves message with valid request body to database
// Takes pointer to the email struct as argument
// If fails returns ErrBodyNotAllowed
func (s *cassandraPostRepo) Create(e *models.Email) error {

	if isValid := PostRequestValidation(e); !isValid {
		return http.ErrBodyNotAllowed
	}

	err := checkmail.ValidateFormat(e.Email)
	if err != nil {
		return http.ErrBodyNotAllowed
	}

	PostEmail(e, s.session)

	return nil
}

// SendEmails sends emails with provided magicNumber and then deletes them from database
// Returns array containing all sent emails and any error encountered.
func (s *cassandraPostRepo) SendEmails(magicNumber int) (error, []string){
	var emailsToSend []models.Email
	e := new(models.Email)
	var emails []string

	iter := s.session.Query(SELECT_EMAIL_TO_SEND,
		magicNumber).Iter()

	for iter.Scan(&e.Email, &e.Title, &e.Content) {
		emailsToSend = append(emailsToSend, *e)
	}
	SendEmail(emailsToSend)
	for el := range emailsToSend {
		emails = append(emails, emailsToSend[el].Email)
		if err := s.session.Query(DELETE_MESSAGE,
			emailsToSend[el].Email, magicNumber).Exec(); err != nil {
			log.Fatal(err)
		}
	}
	emailsToSend = nil
	return nil, emails
}

// ViewMessages displays pageNumber of email string saved in database in batches of pageLimit=4
// Returns array of email struct and any error it encounters
func (s *cassandraPostRepo) ViewMessages(pageNumber int, email string) ([]models.Email, error) {

	var emailToDisplay []models.Email
	pageLimit := 4
	e := new(models.Email)

	var numberOfEmails = GetEmailCount(email, s.session)
	var firstRowEmail = (pageNumber * pageLimit) - pageLimit


	if err := s.session.Query(SELECT_EMAIL, email).PageState(pageState).Scan(&e.Email, &e.Title, &e.Content, &e.MagicNumber); err != nil {
		return nil, err
	}

	for i := 0; i < pageNumber; i++ {

		if numberOfEmails <= firstRowEmail {
			return nil, http.ErrBodyNotAllowed
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


// GetEmailCount returns count of rows with specified email in the database
func GetEmailCount(email string, s *gocql.Session) int {
	var count int
	iter := s.Query(SELECT_COUNT, email).Iter()
	for iter.Scan(&count) {
	}
	return count
}

// PostRequestValidation returns false if any value of Email struct is missing or if email address is invalid
// Else returns true
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

// PostEmail injects specified message into database
func PostEmail(e *models.Email, session *gocql.Session) {
	if err := session.Query(INSERT, e.Email, e.Title, e.Content, e.MagicNumber).Exec(); err != nil {
		log.Println(err)
	}
}

// SendEmail takes array of Email structs and sends them through SMTP
// Returns array of email addresses
func SendEmail(e []models.Email) []string{

	var emails []string
	s := NewSmtpConfig()
	addr := s.SmtpAddress + s.SmtpPort
	auth := smtp.PlainAuth(" ", s.SmtpEmail, s.SmtpPass, s.SmtpAddress)

	for _, elem := range e {

		msg := []byte("To:" + elem.Email + "\r\n" +
			"Subject:" + elem.Title + "\r\n" +
			"\r\n" +
			elem.Content + "\r\n")
		emails = append(emails, elem.Email)
		to := []string{elem.Email}
		err := smtp.SendMail(addr, auth, s.SmtpEmail, to, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	return emails
}

// NewSmtpConfig initializes SMTP config with environment variables provided in appEnv.env file
func NewSmtpConfig() *models.SmtpConfig {
	return &models.SmtpConfig{
		SmtpAddress: os.Getenv("SMTP_SERV"),
		SmtpPort:    os.Getenv("SMTP_PORT"),
		SmtpEmail:   os.Getenv("FROM"),
		SmtpPass:    os.Getenv("PASS"),
	}
}
