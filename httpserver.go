// Author: Harsh Nayyar
// Copyright 2017 - PocketHealth

package main

import (
	"database/sql"
	"net/http"

	"github.com/frankweilin/challenge/controllers"
	_ "github.com/mattn/go-sqlite3"
)

var listeningPort = "8080"

func setDataBase() {
	db, _ := sql.Open("sqlite3", "./pockethealth.db")
	defer db.Close()
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS Subscriber (id INTEGER PRIMARY KEY, name TEXT, email TEXT, phone TEXT, color TEXT)")
	statement.Exec()

}

func main() {
	setDataBase()

	mux := http.NewServeMux()
	mux.HandleFunc("/", controllers.HandleDefault)
	mux.HandleFunc("/subscribe", controllers.HandleSubscribe)
	mux.HandleFunc("/subscribeconfirm", controllers.HandleSubscribeConfirm)
	mux.HandleFunc("/v1/internal/subscriptions/list", controllers.HandleSubscriberAPI)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// specify listening port
	server := &http.Server{Addr: ":" + listeningPort, Handler: mux}
	// start the server, listening to LISTENING_PORT
	server.ListenAndServeTLS("cert.pem", "key.pem")

}
