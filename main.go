package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/go-github/v71/github"
)

type GitHubPushPayload struct {
	After string `json:"after"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		http.Error(w, "Could not read request", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var payload GitHubPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// use git client to fetch the latest commit
	// For example, you can use the "github.com/go-git/go-git/v5" package to interact with git repositories.
	// This is a placeholder for the actual git client code.
	client := github.NewClient(nil)
	ref, resp, err := client.Git.GetRef(context.Background(), "Julian-Chu", "github-webhook-test", "main")
	if err != nil {
		log.Printf("Failed to fetch latest commit: %v", err)
	}
	_ = resp
	log.Printf("Latest commit SHA1: %s", ref.Object.SHA)

	log.Printf("Received push event. Newest commit SHA1: %s", body)
	fmt.Fprintf(w, "Received commit SHA1: %s\n", payload.After)
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	port := "8080"
	log.Printf("Server listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
