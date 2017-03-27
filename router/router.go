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
	h.Handle("/", http.HandlerFunc(home))
	return h
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

			/*--------------------------------------------------------------------------*/

			//cookie, _ := req.Cookie("SessionID")

			//reqBody,_ := req.GetBody()
			logRequest, err := model.NewRequest(req.Method,req.URL.String(),bytes.NewReader(b))
			utils.CopyHeader(req.Header,logRequest.Header)

			SessionID := model.Session{
				SessionId: "dfgsdfgsdfgd",
				TTL: 111111111111,
			}

			respBosy,_ := ioutil.ReadAll(rs.Body)
			logResponse := http.Response{}

			temp := bytes.NewReader(respBosy)
			logResponse.Body = ioutil.NopCloser(temp)


			HttpRequest := model.HttpRequestNode{
				SessionId: SessionID,
				AfterNodeId: "11111111111",
				BeforeNodeId: "11111111111",
				HttpResp: logResponse,
				HttpReq: logRequest,
				NodeId: "11111111111",
			}

			err = db.SaveToElasticSearch(HttpRequest)
			if err != nil {
				log.Println(err)
			}

			/*--------------------------------------------------------------------------*/

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
		/*var clIP string
		if r.Header.Get("X-Forwarded-For") == "" {
			clIP = r.RemoteAddr
		} else {
			clIP = r.Header.Get("X-Forwarded-For")
		}

		uAgent := r.Header.Get("User-Agent")
		log.Printf("\"Method\": \"%s\", \"User-Agent\": \"%s\", \"URL\": \"%s\", \"Host\": \"[%s]\", \"Client-IP\": \"%v\"", r.Method, uAgent, r.URL, r.Host, clIP)*/
		next.ServeHTTP(w, r)
	})
}

func sessionMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("SessionID")
		if err != nil {
			next.ServeHTTP(resp,req)
		}

		if cookie == nil {
			t1 := time.Now()
			unixtime := t1.Unix()
			uAgent := req.Header.Get("User-Agent")
			cookie, err := utils.CreateSessionCoockie(uAgent,t1)
			if err.Error() == "new sess" {
				t1 = time.Now()
				unixtime = t1.Unix()
				cookie, err = utils.CreateSessionCoockie(uAgent,t1)
				if err.Error() == "new sess" {
					next.ServeHTTP(resp,req)
				}

				cookieMonster := &http.Cookie{
					Name: "SessionID",
					Expires: t1,
					Value: cookie,
					HttpOnly: true,
					MaxAge: int(unixtime),
					Path: "/",
				}

				var sessionItem model.Session
				sessionItem.SessionId = cookie
				sessionItem.TTL = unixtime
				utils.SaveSession(sessionItem)

				http.SetCookie(resp,cookieMonster)
				next.ServeHTTP(resp,req)
			}

			cookieMonster := &http.Cookie{
				Name: "SessionID",
				Expires: t1,
				Value: cookie,
				HttpOnly: true,
				MaxAge: int(unixtime),
				Path: "/",
			}

			var sessionItem model.Session
			sessionItem.SessionId = cookie
			sessionItem.TTL = unixtime
			utils.SaveSession(sessionItem)

			http.SetCookie(resp,cookieMonster)
			next.ServeHTTP(resp,req)
		}

		re := utils.CheckSession(cookie.Value)

		if re == true {
			t1 := time.Now()
			unixtime := t1.Unix()
			uAgent := req.Header.Get("User-Agent")
			cookie, err := utils.CreateSessionCoockie(uAgent,t1)
			if err.Error() == "new sess" {
				t1 = time.Now()
				unixtime = t1.Unix()
				cookie, err = utils.CreateSessionCoockie(uAgent,t1)
				if err.Error() == "new sess" {
					next.ServeHTTP(resp,req)
				}

				cookieMonster := &http.Cookie{
					Name: "SessionID",
					Expires: t1,
					Value: cookie,
					HttpOnly: true,
					MaxAge: int(unixtime),
					Path: "/",
				}

				var sessionItem model.Session
				sessionItem.SessionId = cookie
				sessionItem.TTL = unixtime
				utils.SaveSession(sessionItem)

				http.SetCookie(resp,cookieMonster)
				next.ServeHTTP(resp,req)
			}
			cookieMonster := &http.Cookie{
				Name: "SessionID",
				Expires: t1,
				Value: cookie,
				HttpOnly: true,
				MaxAge: int(unixtime),
				Path: "/",
			}

			var sessionItem model.Session
			sessionItem.SessionId = cookie
			sessionItem.TTL = unixtime
			utils.SaveSession(sessionItem)

			http.SetCookie(resp,cookieMonster)
			next.ServeHTTP(resp,req)
		}


		next.ServeHTTP(resp,req)
	})
}
