package main

import (
	"LPManager/router"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080",router.NewMux()))
}