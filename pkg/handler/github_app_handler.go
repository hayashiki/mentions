package handler

import (
	"log"
	"net/http"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Webhook!!")
}
