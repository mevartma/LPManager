package utils

import (
	"LPManager/model"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func CopyHeader(from, to http.Header) {
	for k, v := range from {
		for _, val := range v {
			to.Add(k, val)
		}
	}
	to.Set("Content-Type", from.Get("Content-Type"))
}

func CreateSalt(u *model.User) (string, error) {
	stringBeforeENC := fmt.Sprintf("%s%s", u.Email, u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(stringBeforeENC), 10)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return string(hashedPassword), nil
}
