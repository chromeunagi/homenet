package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var port = flag.Int("port", -1, "Port to listen on")

type PageVariables struct {
	Date string
	Time string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	homePageVars := PageVariables{
		Date: now.Format("02-01-2006"),
		Time: now.Format("15:04:05"),
	}

	t, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		log.Printf("template parsing error: %s", err)
	}
	if err := t.Execute(w, homePageVars); err != nil {
		log.Printf("template executing error: %s", err)
	}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	flag.Parse()

	// Set up HTTP server and necessary handlers
	server := &http.Server{Addr: fmt.Sprintf(":%d", *port)}
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/homepage", homePage)

	// Start listening on HTTP server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf(err.Error())
		}
	}()

	// Set up kernel signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf(err.Error())
	}
}
