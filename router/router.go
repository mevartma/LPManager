package router

import (
	"LPManager/db"
	"LPManager/model"
	"LPManager/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "strconv"
	"strings"
	"time"
)

var pages []model.ProxySetting

func init() {
	err := updatePages()
	if err != nil {
		log.Fatal(err)
	}
}

//NewMux return Handler by URL path
func NewMux() http.Handler {
	h := http.NewServeMux()
	fs := http.FileServer(http.Dir("templates/"))
	h.Handle("/app/", loggerMid(http.StripPrefix("/app", fs)))
	h.Handle("/api/v1/proxy", loggerMid(http.HandlerFunc(proxy)))
	h.Handle("/", loggerMid(http.HandlerFunc(home)))
	return h
}

func logoutUser(resp http.ResponseWriter, req *http.Request) {
}

func loginUser(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var user model.User

	user.Password = fmt.Sprintf("%v", req.Form["password"])
	user.UserName = fmt.Sprintf("%v", req.Form["username"])

	ps, err := utils.CreateSalt(&user)
	if err != nil {
		log.Fatal(err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	tUser, err := db.GetUser(user.UserName)
	if err != nil {
		log.Fatal(err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	if tUser.Salt != ps {
		//status not auth
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookieMonster := &http.Cookie{
		Name:    "SessionID",
		Expires: time.Now().AddDate(0, 0, 1),
		Value:   ps,
	}

	http.SetCookie(resp, cookieMonster)
	http.Redirect(resp,req,"/app",http.StatusOK)
}

func registerUser(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var user model.User
	var inUser model.InternalUsers
	var endStatus model.Status

	user.Email = fmt.Sprintf("%v", req.Form["email"])
	user.Password = fmt.Sprintf("%v", req.Form["password"])
	user.UserName = fmt.Sprintf("%v", req.Form["username"])

	ps, err := utils.CreateSalt(&user)
	if err != nil {
		log.Fatal(err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	inUser.UserName = user.UserName
	inUser.Email = user.Email
	inUser.Salt = ps

	err = db.UpdateUser(inUser, "add")
	if err != nil {
		log.Fatal(err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	endStatus.Message = fmt.Sprintf("User %s with email %s has created", inUser.UserName, inUser.Email)
	endStatus.Action = "Create User"
	js, err := json.Marshal(endStatus)
	if err != nil {
		log.Fatal(err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}
	resp.Header().Set("Content-type", "application/json")
	resp.Write(js)
	return
}

func proxy(resp http.ResponseWriter, req *http.Request) {
	var p model.ProxySetting
	var err error
	if req.Method != "GET" {
		err = json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	switch req.Method {
	case "POST":
		exist := false
		for _, pr := range pages {
			if pr.FullURL == p.FullURL {
				exist = true
			}
		}
		if exist == false {
			err = db.UpdateProxy(p, "add")
			err = updatePages()
		}
	case "GET":
		err = updatePages()
	case "PUT":
		err = db.UpdateProxy(p, "update")
		err = updatePages()
	case "DELETE":
		err = db.UpdateProxy(p, "delete")
		err = updatePages()
	default:
		err = errors.New("Method Not Allow")
	}
	if err != nil {
		fmt.Println(err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(pages)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-type", "application/json")
	resp.Write(js)
	return
}

func home(resp http.ResponseWriter, req *http.Request) {
	for _, page := range pages {
		startWith := strings.HasPrefix(strings.ToLower(req.RequestURI), strings.ToLower(page.LocalPath))
		ref := strings.Contains(strings.ToLower(req.Referer()), strings.ToLower(page.FullURL))
		if startWith == true || ref == true {
			defer req.Body.Close()
			var newURI, newURL string
			if page.RemotePath == "/" {
				if strings.Contains(req.RequestURI, ".css") {
					newURI = strings.Replace(req.RequestURI, page.LocalPath, page.RemotePath, -1)
				} else {
					newURI = strings.Replace(strings.ToLower(req.RequestURI), strings.ToLower(page.LocalPath), "", -1)
				}
			} else {
				newURI = strings.Replace(req.RequestURI, page.LocalPath, page.RemotePath, -1)
			}
			newURL = fmt.Sprintf("http://%s%s", page.RemoteHost, newURI)
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			r, err := http.NewRequest(req.Method, newURL, bytes.NewReader(b))
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			utils.CopyHeader(req.Header, r.Header)
			if r.TLS != nil {
				r.Header.Set("X-Forwarded-Proto", "https")
			}

			client := http.Client{}
			rs, err := client.Do(r)
			if err != nil {
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}

			utils.CopyHeader(rs.Header, resp.Header())
			resp.WriteHeader(rs.StatusCode)
			io.Copy(resp, rs.Body)
		}
	}
}

func updatePages() error {
	pps, err := db.GetAllProxies()
	pages = nil
	for _, ps := range *pps {
		pages = append(pages, ps)
	}
	return err
}

func loggerMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clIP string
		if r.Header.Get("X-Forwarded-For") == "" {
			clIP = r.RemoteAddr
		} else {
			clIP = r.Header.Get("X-Forwarded-For")
		}

		uAgent := r.Header.Get("User-Agent")
		log.Printf("\"Method\": \"%s\", \"User-Agent\": \"%s\", \"URL\": \"%s\", \"Host\": \"[%s]\", \"Client-IP\": \"%v\"", r.Method, uAgent, r.URL, r.Host, clIP)
		next.ServeHTTP(w, r)
	})
}

/*func userMid(next http.Handler) http.Handler {
	return http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		coocki, _ := r.Cookie("SessionID")
	})
}*/
