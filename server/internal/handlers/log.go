package handlers

import (
	"log"
	"net/http"
)

func logError(err error, r *http.Request) {
	log.Printf("%v: %v", r.URL.Path, err)
}
