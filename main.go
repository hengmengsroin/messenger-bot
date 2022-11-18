package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// VERIFY_TOKEN use to verify the incoming request
	VERIFY_TOKEN = "12345"
)

// webhook is a handler for Webhook server
func webhook(w http.ResponseWriter, r *http.Request) {
	// return all with status code 200
	w.WriteHeader(http.StatusOK)

	// method that allowed are GET & POST
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		log.Printf("invalid method: not get or post")
		return
	}

	// if the method of request is GET
	if r.Method == http.MethodGet {
		// read token from query parameter
		verifyToken := r.URL.Query().Get("hub.verify_token")

		// verify the token included in the incoming request
		if verifyToken != VERIFY_TOKEN {
			log.Printf("invalid verification token: %s", verifyToken)
			return
		}

		// write string from challenge query parameter
		if _, err := w.Write([]byte(r.URL.Query().Get("hub.challenge"))); err != nil {
			log.Printf("failed to write response body: %v", err)
		}

		return
	}

	// read body in the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
	}

	// print the body
	log.Printf(string(body))

	return
}

func main() {
	// create the handler
	handler := http.NewServeMux()
	handler.HandleFunc("/", webhook)

	// configure http server
	srv := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("localhost:%d", 3000),
	}

	// start http server
	log.Printf("http server listening at %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
