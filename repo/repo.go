package repo

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"reflect"

	"github.com/sawickiszymon/gowebapp/models"
)
const (
	INSERT               = `INSERT INTO Email (email, title, content, magic_number) VALUES (?, ?, ?, ?) USING TTL 300`
	SELECT_EMAIL_TO_SEND = `SELECT email, title, content FROM Email WHERE magic_number = ?`
	SELECT_EMAIL         = `SELECT email, title, content, magic_number FROM Email WHERE email = ?`
	SELECT_COUNT         = `SELECT Count(*) FROM Email WHERE email = ?`
	DELETE_MESSAGE       = `DELETE FROM Email WHERE email = ? AND magic_number = ?`
)

//func tsett(pageNumber int, email string) []models.Email {
//	var pageLimit = 4
//	var s gocql.Session
//
//	var numberOfEmails = GetEmailCount(email, s)
//	var firstRowEmail = (pageNumber * pageLimit) - pageLimit
//
//	if err := s.Query(SELECT_EMAIL, ps.ByName("email")).PageState(pageState).Scan(&e.Email, &e.Title, &e.Content, &e.MagicNumber); err != nil {
//		log.Println(err)
//	}
//
//	for i := 0; i < pageNumber; i++ {
//
//		if numberOfEmails <= firstRowEmail {
//			json.NewEncoder(writer).Encode("There is no emails to display")
//			return
//		}
//
//		iter := s.Query(SELECT_EMAIL, e.Email).PageState(pageState).PageSize(pageLimit).Iter()
//
//		for iter.Scan(&e.Email, &e.Title, &e.Content, &e.MagicNumber) {
//			if pageNumber%2 == 1 && i+1 == pageNumber {
//				emailToDisplay = append(emailToDisplay, e)
//			} else if pageNumber%1 == 0 && i+1 == pageNumber {
//				emailToDisplay = append(emailToDisplay, e)
//			}
//			pageState = iter.PageState()
//		}
//	}
//	return nil
//}
//
//func GetEmailCount(email string, s *gocql.Session) int {
//	var count int
//	iter := s.Query(SELECT_COUNT, email).Iter()
//	for iter.Scan(&count) {
//	}
//	return count
//}


func PostRequestValidation(e models.Email) bool {
	isValid := true
	v := reflect.ValueOf(e)
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

func DecodeRequest(w http.ResponseWriter, r *http.Request) models.Email {
	var e models.Email
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return e
}

