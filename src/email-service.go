package main

import (
	"bytes"
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

	InfoLogger.Printf("successfully posted cronMessage %s", strings.TrimSpace(bytes.NewBuffer(bod).String()))
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	InfoLogger.Println("Starting up email-service")

	InfoLogger.Println("trigger a cronmessage function")
	cronPublish("publish my cron message")
}
