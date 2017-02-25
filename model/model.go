package model

type ProxySetting struct {
	Id         int64  `json:"id"`
	LocalPath  string `json:"localpath"`
	FullURL    string `json:"fullurl"`
	RemoteHost string `json:"remotehost"`
	RemotePath string `json:"remotepath"`
}

type User struct {
	Id       int64
	UserName string
	Password string
	Email    string
}

type InternalUsers struct {
	Id       int64
	UserName string
	Email    string
	Salt     string
}

type Status struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type SessionStore struct {
	Id      int64
	Session string
}
