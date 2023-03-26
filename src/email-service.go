package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func Init(
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	DebugLogger = log.New(debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	WarningLogger = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ErrorLogger = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

type Message struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Email   string `json:"email"`
	Content string `json:"content"`
}

func sending(res http.ResponseWriter, req *http.Request) {

	// todo: future me: remove this logging calls from function body
	InfoLogger.Println("sending function called with route")

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

func defaultRoute(res http.ResponseWriter, req *http.Request) {

	// todo: remove logging from function body
	WarningLogger.Printf("defaultRoute is hit with route: %s", req.URL.Path)

	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	msg := Message{
		Id:      0,
		Title:   "root endpoint",
		Email:   "void@socratic.dev",
		Content: "nothing here",
	}

	json.NewEncoder(res).Encode(msg)
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	InfoLogger.Println("Starting up email-service")

	http.HandleFunc("/sending", sending)

	http.HandleFunc("/", defaultRoute)

	err := http.ListenAndServe(":4000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		WarningLogger.Println("server closed")
	} else if err != nil {
		ErrorLogger.Println("error starting server:", err)
		os.Exit(1)
	}

}
