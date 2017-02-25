package proxyserver

import (
	"LPManager/model"
	"net/http"
)

type ProxyServer struct {
	Client http.Client
}

func (p *ProxyServer) ProxyRequest(setting model.ProxySetting, req *http.Request) () {

}

func NewProxyServer() *ProxyServer {
	var p ProxyServer
	return &p
}
