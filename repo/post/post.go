package repo

import (
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"reflect"

	models "github.com/sawickiszymon/gowebapp/models"
	repo "github.com/sawickiszymon/gowebapp/repo"
)
const (
	INSERT               = `INSERT INTO Email (email, title, content, magic_number) VALUES (?, ?, ?, ?) USING TTL 300`
	SELECT_EMAIL_TO_SEND = `SELECT email, title, content FROM Email WHERE magic_number = ?`
	SELECT_EMAIL         = `SELECT email, title, content, magic_number FROM Email WHERE email = ?`
	SELECT_COUNT         = `SELECT Count(*) FROM Email WHERE email = ?`
	DELETE_MESSAGE       = `DELETE FROM Email WHERE email = ? AND magic_number = ?`
)
func NewRepo(s *gocql.Session) repo.PostRepo {
	return &cassandraPostRepo{
		session: s,
	}
}

type cassandraPostRepo struct {
	session *gocql.Session
}


func (s *cassandraPostRepo) Create(e *models.Email) error{

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

func PostRequestValidation(e *models.Email) bool {
	isValid := true
	v := reflect.Indirect(reflect.ValueOf(e))
	fmt.Println(reflect.Indirect(reflect.ValueOf(e)))
	fmt.Println(reflect.ValueOf(e))
	fmt.Println(v.NumField())
	fmt.Println(reflect.ValueOf(e).NumField())
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

