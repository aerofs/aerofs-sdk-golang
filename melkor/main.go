package main

// The entrypoint for the Melkor webapp demonstrating the AeroFS Golang SDK

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

// Global logger
var logger *log.Logger

const (
	hostName = "localhost:1337"
)

func main() {
	err := initLogger()
	if err != nil {
		fmt.Println("Unable to initialize log file")
		os.Exit(1)
	}
	logger.Print("Melkor beginning startup...")

	// Set Handlers
	router := mux.NewRouter()
	router.HandleFunc("/", defaultHandler).Methods("GET")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/{miscellaneous}", MiscHandler).Methods("GET")

	//r.GET("/test", test_1)
	//r.GET("/", arrive)
	//r.GET("/redirect", redirect)
	//r.GET("/tokenization", tokenization)
	http.ListenAndServe("localhost:1337", router)
}

// Initialize the Global server logger
func initLogger() error {
	t := time.Now()
	logTime := fmt.Sprintf("%d-%d-%d_%d-%d-%d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	logName := fmt.Sprintf("logs/Melkor_Logs_%s", logTime)
	logFile, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// Log the file location, time and date
	logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	return err
}
