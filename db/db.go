package db

import (
	"LPManager/model"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
	"log"
)

const (
	dbUrl  string = "root:@tcp(localhost:3306)/test"
	server string = "mysql"
)

var (
	baseURL = "http://35.157.18.149:9200/"
)

func SaveToElasticSearch(d model.HttpRequestNode) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s/%s",baseURL,"q","http")
	log.Println(url)
	jsData, err := json.Marshal(d)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST",url,bytes.NewBuffer(jsData))
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

/*
func UpdateHttpNode(h model.HttpRequestNode, m string) error {
	var err error
	sessionDB, err := sql.Open(server, dbUrl)
	defer sessionDB.Close()

	switch m {
	case "delete":
	case "update":
	case "add":
	}
	return err
}

func GetNodeById(nodeid, sessionid string) (*[]model.HttpRequestNode, error) {
	var results []model.HttpRequestNode
	sessionDB, err := sql.Open(server, dbUrl)
	if err != nil {
		return nil, err
	}
	defer sessionDB.Close()

	stmt, err := sessionDB.Prepare("SELECT * FROM HttpNodes WHERE sessionid=? AND nodeid=?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(sessionid,nodeid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.HttpRequestNode
		err = rows.Scan(&r.NodeId,&r.BeforeNodeId,&r.AfterNodeId,&r.HttpReq,&r.HttpResp,&r.SessionId)
		results = append(results,r)
	}

	return &results,err
}
*/
func GetRedirects() (*[]model.RedirectType, error) {
	var results []model.RedirectType
	sessionDB, err := sql.Open(server, dbUrl)
	if err != nil {
		return nil, err
	}
	defer sessionDB.Close()

	query := "select * from redirects"
	rows, err := sessionDB.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.RedirectType
		err = rows.Scan(&r.Id, &r.From, &r.To, &r.Domain)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return &results, err
}

func UpdateRedirect(r model.RedirectType, m string) error {
	var err error
	sessionDB, err := sql.Open(server, dbUrl)
	defer sessionDB.Close()

	switch m {
	case "delete":

		stmt, err := sessionDB.Prepare("DELETE FROM redirects WHERE id=? LIMIT 1")
		_, err = stmt.Exec(r.Id)
		return err
	case "update":
		stmt, err := sessionDB.Prepare("UPDATE redirects SET urlfrom=?, urlto=?, urldomain=? where id=?")
		_, err = stmt.Exec(r.From, r.To, r.Domain, r.Id)
		return err
	case "add":
		stmt, err := sessionDB.Prepare("INSERT INTO redirects(urlfrom,urlto,urldomain) VALUES (?,?,?)")
		_, err = stmt.Exec(r.From, r.To, r.Domain)
		return err
	default:
		err = errors.New("Command Not Found")
	}
	return err
}

func GetRedirect(s string) (*model.RedirectType, error) {
	var result model.RedirectType
	sessionDB, err := sql.Open(server, dbUrl)
	if err != nil {
		return nil, err
	}
	defer sessionDB.Close()

	query := "SELECT * FROM redirects WHERE urlfrom=? LIMIT 1"
	rows, err := sessionDB.Query(query, s)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&result.Id, &result.From, &result.To, &result.Domain)
	}
	return &result, err
}

func UpdateProxy(p model.ProxySetting, m string) error {
	var err error
	sessionDB, err := sql.Open(server, dbUrl)
	defer sessionDB.Close()

	switch m {
	case "update":
		stmt, err := sessionDB.Prepare("UPDATE proxy SET local_path=? ,full_url=? ,remote_host=? ,remote_path=? WHERE id=?")
		_, err = stmt.Exec(p.LocalPath, p.FullURL, p.RemoteHost, p.RemotePath, p.Id)
		return err
	case "delete":
		stmt, err := sessionDB.Prepare("DELETE FROM proxy WHERE id=? LIMIT 1")
		_, err = stmt.Exec(p.Id)
		return err
	case "add":
		stmt, err := sessionDB.Prepare("INSERT INTO proxy(local_path,full_url,remote_host,remote_path) VALUES(?,?,?,?)")
		_, err = stmt.Exec(p.LocalPath, p.FullURL, p.RemoteHost, p.RemotePath)
		return err
	default:
		err = errors.New("Command Not Found")
	}

	return err
}

func GetAllProxies() (*[]model.ProxySetting, error) {
	var results []model.ProxySetting
	sessionDB, err := sql.Open(server, dbUrl)
	if err != nil {
		return nil, err
	}
	defer sessionDB.Close()

	query := "SELECT * FROM proxy"
	rows, err := sessionDB.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.ProxySetting
		err = rows.Scan(&r.Id, &r.LocalPath, &r.FullURL, &r.RemoteHost, &r.RemotePath)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return &results, err
}
