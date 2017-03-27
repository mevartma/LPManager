package main

import (
	"LPManager/router"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe("192.168.1.86:8080", router.NewMux()))
}
