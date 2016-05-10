package main

import (
	"github.com/rs/cors"
	"github.com/stinkyfingers/AuthApi/database"
	"github.com/stinkyfingers/AuthApi/routes"

	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	listenAddr = flag.String("port", "8082", "Assign http port")
)

func main() {
	var err error
	flag.Parse()

	if err = database.Init(); err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	// r := router.New()
	r := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	}).Handler(router.New())

	osPort := os.Getenv("PORT")
	if osPort != "" {
		*listenAddr = osPort
	}

	srv := &http.Server{
		Addr:         ":" + *listenAddr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on 127.0.0.1:%s\n", *listenAddr)
	log.Fatal(srv.ListenAndServe())

}
