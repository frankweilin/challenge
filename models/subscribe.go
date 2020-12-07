package models

import "database/sql"

type Subscriber struct {
	Name     string
	Email    string
	Tel      string
	FavColor string
}

func (subscriber *Subscriber) UpdateSubscriber() {
	db, _ := sql.Open("sqlite3", "./pockethealth.db")
	defer db.Close()
	// check previous subscription by email
	rows, _ := db.Query("SELECT id FROM Subscriber WHERE email = ?", subscriber.Email)
	defer rows.Close()
	id := 0
	for rows.Next() {
		rows.Scan(&id)
	}
	if id == 0 {
		//New Subscriber
		statement, _ := db.Prepare("INSERT INTO Subscriber (name, email, phone,color) VALUES (?, ? ,?,?)")
		statement.Exec(subscriber.Name, subscriber.Email, subscriber.Tel, subscriber.FavColor)
	} else {
		// Existing Subscriber
		statement, _ := db.Prepare("UPDATE Subscriber set name=?,  phone=?,color =? WHERE email =?")
		statement.Exec(subscriber.Name, subscriber.Tel, subscriber.FavColor, subscriber.Email)
	}

}

func GetAllSubscriber() []*Subscriber {
	var users = []*Subscriber{}
	database, _ := sql.Open("sqlite3", "./pockethealth.db")
	defer database.Close()
	rows, _ := database.Query("SELECT name, email,phone,color FROM Subscriber")
	for rows.Next() {
		var subscribe Subscriber
		rows.Scan(&subscribe.Name, &subscribe.Email, &subscribe.Tel, &subscribe.FavColor)
		users = append(users, &subscribe)
	}
	return users
}
