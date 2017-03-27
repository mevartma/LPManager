package utils

import (
	"LPManager/model"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"errors"
)

var CookieJar []model.Session

func init() {
	go func() {
		for {
			go sessionsMaintance()
			time.Sleep(1 * time.Second)
		}
	}()
}

func CopyHeader(from, to http.Header) {
	for k, v := range from {
		for _, val := range v {
			to.Add(k, val)
		}
	}
	to.Set("Content-Type", from.Get("Content-Type"))
}

func GetSessionById(s string) *model.Session {
	for _, sess := range CookieJar {
		if sess.SessionId == s {
			return &sess
		}
	}
	return nil
}

func SaveSession(sess model.Session) {
	result := false
	for _, session := range CookieJar {
		if sess.SessionId == session.SessionId {
			result = true
		}
	}
	if result == false {
		CookieJar = append(CookieJar, sess)
	}
}

func CheckSession(sessionId string) bool {
	for _, sess := range CookieJar {
		if sess.SessionId == sessionId {
			return true
		}
	}
	return false
}

func sessionsMaintance() {
	var temp []model.Session
	for _, sess := range CookieJar {
		sess.TTL -= 1
		if sess.TTL == 0 || sess.TTL < 0 {
			continue
		} else {
			temp = append(temp, sess)
		}
	}

	CookieJar = nil
	for _, tSess := range temp {
		CookieJar = append(CookieJar, tSess)
	}
}

func CreateSessionCoockie(toHash string, t time.Time) (string,error) {
	h := sha1.New()
	var err error
	stringToHash := fmt.Sprintf("%s%v", toHash, t)
	h.Write([]byte(stringToHash))

	sha512String := hex.EncodeToString(h.Sum(nil))
	result := CheckSession(sha512String)
	if result == true {
		err = errors.New("new sess")
		return "",err
	}
	return sha512String,nil
}
