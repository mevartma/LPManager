package db

import (
	"LPManager/model"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUrl         string = "proxymanger:YaABs8NW@tcp(localhost:3306)/proxymanger"
	server        string = "mysql"
)

func UpdateProxy(p model.ProxySetting, m string) error {
	var err error
	db, err := sql.Open(server, dbUrl)
	defer db.Close()

	switch m {
	case "update":
		stmt, err := db.Prepare("UPDATE proxy SET local_path=? ,full_url=? ,remote_host=? ,remote_path=? WHERE id=?")
		_, err = stmt.Exec(p.LocalPath, p.FullURL, p.RemoteHost, p.RemotePath, p.Id)
		return err
	case "delete":
		stmt, err := db.Prepare("DELETE FROM proxy WHERE id=? LIMIT 1")
		_, err = stmt.Exec(p.Id)
		return err
	case "add":
		stmt, err := db.Prepare("INSERT INTO proxy(local_path,full_url,remote_host,remote_path) VALUES(?,?,?,?)")
		_, err = stmt.Exec(p.LocalPath, p.FullURL, p.RemoteHost, p.RemotePath)
		return err
	default:
		err = errors.New("Command Not Found")
	}

	return err
}

func GetAllProxies() (*[]model.ProxySetting, error) {
	var results []model.ProxySetting
	db, err := sql.Open(server, dbUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM proxy"
	rows, err := db.Query(query)
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