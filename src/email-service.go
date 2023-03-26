package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Message struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Email   string `json:"email"`
	Content string `json:"content"`
}

func sending(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	msg := Message{
		Id:      1,
		Title:   "Very Important Message!!!",
		Email:   "test@socratic.dev",
		Content: "Lorem ipsum dipsum more",
	}

	json.NewEncoder(res).Encode(msg)
}

func main() {
	http.HandleFunc("/sending", sending)
	err := http.ListenAndServe(":4000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
