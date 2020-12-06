// Author: Harsh Nayyar
// Copyright 2017 - PocketHealth

package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var LISTENING_PORT = "8080"

type SubscribeConfirmPage struct {
	Name     string
	Email    string
	Tel      string
	FavColor string
}

func handleSubscribeConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		subscribeConfirmPage := SubscribeConfirmPage{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Tel:      r.FormValue("tel"),
			FavColor: r.FormValue("color"),
		}
		updateSubscriber(subscribeConfirmPage)
		t, _ := template.ParseFiles("tmpl/subscribeconfirm.html")
		t.Execute(w, subscribeConfirmPage)
	} else {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/default.html")
	t.Execute(w, "")
}
func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/subscribe.html")
	t.Execute(w, "")
}

func handleSubscriberAPI(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(getAll(), w)
}
func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

func setDataBase() {
	db, _ := sql.Open("sqlite3", "./pockethealth.db")
	defer db.Close()
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS Subscriber (id INTEGER PRIMARY KEY, name TEXT, email TEXT, phone TEXT, color TEXT)")
	statement.Exec()

}
func getAll() []*SubscribeConfirmPage {
	var users = []*SubscribeConfirmPage{}
	database, _ := sql.Open("sqlite3", "./pockethealth.db")
	defer database.Close()
	rows, _ := database.Query("SELECT name, email,phone,color FROM Subscriber")
	for rows.Next() {
		var subscribe SubscribeConfirmPage
		rows.Scan(&subscribe.Name, &subscribe.Email, &subscribe.Tel, &subscribe.FavColor)
		users = append(users, &subscribe)
	}
	return users
}

func updateSubscriber(subscribeConfirmPage SubscribeConfirmPage) {
	db, _ := sql.Open("sqlite3", "./pockethealth.db")

	// check previous subscription by email
	rows, _ := db.Query("SELECT id FROM Subscriber WHERE email = ?", subscribeConfirmPage.Email)
	defer rows.Close()
	id := 0
	for rows.Next() {
		rows.Scan(&id)
	}
	if id == 0 {
		//New Subscriber
		statement, _ := db.Prepare("INSERT INTO Subscriber (name, email, phone,color) VALUES (?, ? ,?,?)")
		statement.Exec(subscribeConfirmPage.Name, subscribeConfirmPage.Email, subscribeConfirmPage.Tel, subscribeConfirmPage.FavColor)
	} else {
		// Existing Subscriber
		statement, _ := db.Prepare("UPDATE Subscriber set name=?,  phone=?,color =? WHERE email =?")
		statement.Exec(subscribeConfirmPage.Name, subscribeConfirmPage.Tel, subscribeConfirmPage.FavColor, subscribeConfirmPage.Email)
	}

}

func main() {
	setDataBase()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleDefault)
	mux.HandleFunc("/subscribe", handleSubscribe)
	mux.HandleFunc("/subscribeconfirm", handleSubscribeConfirm)
	mux.HandleFunc("/v1/internal/subscriptions/list", handleSubscriberAPI)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// specify listening port
	server := &http.Server{Addr: ":" + LISTENING_PORT, Handler: mux}
	// start the server, listening to LISTENING_PORT
	server.ListenAndServeTLS("cert.pem", "key.pem")

}
