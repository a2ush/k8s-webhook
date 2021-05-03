package main

import (
	"net/http"
	"log"

	"github.com/a2ush/k8s-webhook/src/webhook"
)

func main() {
	log.Print("Server is running...")

	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", webhook.Mutate_handler)
	mux.HandleFunc("/validate", webhook.Validate_handler)

	err := http.ListenAndServeTLS(":8080", "/tls/tls.crt", "/tls/tls.key", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

