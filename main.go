package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/MohamedDenta/Chat-App/chat"
	"github.com/MohamedDenta/Chat-App/config"
)

var configuration config.Configuration
var serverHostName string

// Init init info
func Init() {
	configuration = config.LoadConfigAndSetUpLogging()

	port := os.Getenv("PORT")
	if port == "" {
		serverHostName = fmt.Sprintf("%s:%s", configuration.Hostname, strconv.Itoa(configuration.Port))
	} else {
		serverHostName = fmt.Sprintf(":%s", port)
	}
	log.Println("The serverHost ", serverHostName)

}

func main() {
	Init()
	// websoket server
	server := chat.NewServer()
	//fmt.Println(len(server.Messages))
	go server.Listen()
	// http.HandleFunc("/messages", handleHomePage)
	http.HandleFunc("/", handleHomePage)

	http.ListenAndServe(serverHostName, nil)
}

func handleHomePage(responseWriter http.ResponseWriter, request *http.Request) {
	http.ServeFile(responseWriter, request, "clientlogin.htm")
}
