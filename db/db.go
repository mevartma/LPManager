package db

import (
	"LPManager/model"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUrl  string = "proxymanger:YaABs8NW@tcp(localhost:3306)/proxymanger"
	server string = "mysql"
)

func UpdateUser(u model.InternalUsers, m string) error {
	var err error
	sessionDB, err := sql.Open(server, dbUrl)
	defer sessionDB.Close()

	switch m {
	case "delete":
		stmt, err := sessionDB.Prepare("DELETE FROM users WHERE id=? LIMIT 1")
		_, err = stmt.Exec(u.Id)
		return err
	case "update":
		stmt, err := sessionDB.Prepare("UPDATE users SET username=?, email=?, salt=? where id=?")
		_, err = stmt.Exec(u.UserName, u.Email, u.Salt, u.Id)
		return err
	case "add":
		stmt, err := sessionDB.Prepare("INSERT INTO users(username,email,salt) VALUES (?,?,?)")
		_, err = stmt.Exec(u.UserName, u.Email, u.Salt)
		return err
	default:
		err = errors.New("Command Not Found")
	}

	return err
}

func GetUser(s string) (model.InternalUsers, error) {
	var result model.InternalUsers
	sessionDB, err := sql.Open(server, dbUrl)
	if err != nil {
		return nil, err
	}
	defer sessionDB.Close()

	query := "SELECT * FROM user WHERE username=? LIMIT 1"
	rows, err := sessionDB.Query(query, s)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&result.Id, &result.UserName, &result.Email, &result.Email)
	}
	return result, err
}

func GetAllUsers() (*[]model.User, error) {
	var results []model.User
	sessionDB, err := sql.Open(server, dbUrl)
	if err != nil {
		return nil, err
	}
	defer sessionDB.Close()

	query := "SELECT id,username,email FROM users"
	rows, err := sessionDB.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.User
		err = rows.Scan(&r.Id, &r.UserName, &r.Email)
		results = append(results, r)
	}

	return &results, err
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
