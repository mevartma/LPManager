package replay

import (
	"net/http"
	"LPManager/model"
)



func NewHttpNode(resp http.Response, req *http.Request, sess model.Session) (*model.HttpRequestNode, error) {
	return &model.HttpRequestNode{"", "", "", req, resp, sess}, nil
}

func SaveNode(node *model.HttpRequestNode) error {
	return nil
}
