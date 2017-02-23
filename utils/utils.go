package utils

import (
	"net/http"
)

//CopyHeader copy headers
func CopyHeader(from, to http.Header) {
	for k, v := range from {
		for _, val := range v {
			to.Add(k, val)
		}
	}
	to.Set("Content-Type", from.Get("Content-Type"))
}