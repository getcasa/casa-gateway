package gateway

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

var emailRegExp = "(?:[a-z0-9!#$%&'*+=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])"

type signupReq struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	Firstname            string `json:"firstname"`
	Lastname             string `json:"lastname"`
	Birthdate            string `json:"birthdate"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SignUp route create an user
func SignUp(c echo.Context) error {
	req := new(signupReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	var missingFields []string
	if req.Email == "" {
		missingFields = append(missingFields, "email")
	}
	if req.Password == "" {
		missingFields = append(missingFields, "password")
	}
	if req.PasswordConfirmation == "" {
		missingFields = append(missingFields, "passwordConfirmation")
	}
	if req.Firstname == "" {
		missingFields = append(missingFields, "firstname")
	}
	if req.Lastname == "" {
		missingFields = append(missingFields, "lastname")
	}
	if req.Birthdate == "" {
		missingFields = append(missingFields, "birthdate")
	}
	if len(missingFields) > 0 {
		return c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Some fields missing: " + strings.Join(missingFields, ", "),
		})
	}
	isValid, _ := regexp.MatchString(emailRegExp, req.Email)
	if isValid == false {
		return c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Invalid email",
		})
	}
	if req.Password != req.PasswordConfirmation {
		return c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Passwords mismatch",
		})
	}
	// statement, _ := DB.Prepare("SELECT * FROM people WHERE (?, ?, ?, ?, ?, ?, ?)")
	// row, _ := DB.QueryRow("SELECT * FROM people WHERE")

	hashedPassword, _ := hashPassword(req.Password)
	statement, _ = DB.Prepare("INSERT INTO user (id, email, password, firstname, lastname, birthdate, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	statement.Exec(NewULID().String(), req.Email, hashedPassword, req.Firstname, req.Lastname, req.Birthdate, time.Now())

	return c.JSON(http.StatusCreated, MessageResponse{
		Message: "Account created",
	})
}
