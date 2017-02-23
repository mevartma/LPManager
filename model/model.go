package model

type ProxySetting struct {
	Id         int64  `json:"id"`
	LocalPath  string `json:"localpath"`
	FullURL    string `json:"fullurl"`
	RemoteHost string `json:"remotehost"`
	RemotePath string `json:"remotepath"`
}