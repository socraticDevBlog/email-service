package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Email     string `json:"email"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

func cronPublish(msg string) {
	now := time.Now().UTC().Format(time.RFC3339)
	InfoLogger.Printf("now string is: %s", now)

	params := url.Values{}
	params.Add("datetime", now)
	params.Add("message", msg)

	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://paste.c-net.org/", body)
	if err != nil {
		ErrorLogger.Println("Error building a message content:", err)
	}

	respThirdParty, err := http.DefaultClient.Do(req)
	if err != nil {
		ErrorLogger.Println("Error publishing content:", err)
	}
	bod, err := ioutil.ReadAll(respThirdParty.Body)
	if err != nil {
		ErrorLogger.Println("Could not read response from third party: ", err)
	}
	defer respThirdParty.Body.Close()

	InfoLogger.Printf("successfully posted cronmessage %s", strings.TrimSpace(bytes.NewBuffer(bod).String()))
}

func publish(res http.ResponseWriter, req *http.Request) {

	// todo: future me: remove this logging calls from function body
	InfoLogger.Printf("publish function called with route: %s", req.URL.Path)

	params := url.Values{}
	inMsg := req.PostFormValue("msg")
	params.Add("message", inMsg)

	now := time.Now().UTC().Format(time.RFC3339)
	InfoLogger.Printf("now string is: %s", now)
	params.Add("datetime", now)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://paste.c-net.org/", body)
	if err != nil {
		res.WriteHeader(http.StatusNotAcceptable)
		ErrorLogger.Println("Error building a message content:", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	respThirdParty, err := http.DefaultClient.Do(req)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		ErrorLogger.Println("Error publishing content:", err)
	}

	bod, err := ioutil.ReadAll(respThirdParty.Body)
	if err != nil {
		ErrorLogger.Println("Could not read response from third party: ", err)
	}
	defer respThirdParty.Body.Close()

	res.WriteHeader(http.StatusOK)
	msg := Message{
		Id:        1,
		Title:     "PUBLISHED",
		Email:     "test@socratic.dev",
		Content:   strings.TrimSpace(bytes.NewBuffer(bod).String()),
		Timestamp: now,
	}
	json.NewEncoder(res).Encode(msg)

	InfoLogger.Printf("successfully publish a message '%s' at %s", inMsg, now)

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

	InfoLogger.Println("trigger a cronmessage function")
	cronPublish("publish my cron message")

	http.HandleFunc("/sending", sending)

	http.HandleFunc("/publish", publish)

	http.HandleFunc("/", defaultRoute)

	err := http.ListenAndServe(":4000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		WarningLogger.Println("server closed")
	} else if err != nil {
		ErrorLogger.Println("error starting server:", err)
		os.Exit(1)
	}

}
