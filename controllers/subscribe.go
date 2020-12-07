package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/frankweilin/challenge/models"
)

func HandleSubscribeConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		subscriber := models.Subscriber{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Tel:      r.FormValue("tel"),
			FavColor: r.FormValue("color"),
		}
		subscriber.UpdateSubscriber()
		t, _ := template.ParseFiles("tmpl/subscribeconfirm.html")
		t.Execute(w, subscriber)
	} else {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func HandleDefault(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/default.html")
	t.Execute(w, "")
}
func HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/subscribe.html")
	t.Execute(w, "")
}

func HandleSubscriberAPI(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetAllSubscriber(), w)
}
func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
