package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Trip struct {
	StartDate time.Time
	EndDate   time.Time
}

type User struct {
	Username     string
	StartDate    time.Time
	InvalidDates []Trip
}

// guestbookKey returns the key used for all guestbook entries.
func my_datastore_Key(c context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "my_datastore", "default_dataStore", 0, nil)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c := appengine.NewContext(r)

	template := "02-01-2006"
	timeStr, _ := time.Parse(template, vars["startDate"])

	g := User{
		Username:     vars["username"],
		StartDate:    timeStr,
		InvalidDates: make([]Trip, 1),
	}

	key := datastore.NewIncompleteKey(c, "UserData", my_datastore_Key(c))
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "User added.")
	}
}
func addTripHandler2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c := appengine.NewContext(r)

	q := datastore.NewQuery("UserData").Filter("Username =", vars["username"])
	t := q.Run(c)
	for {
		var u User
		k, err := t.Next(&u)
		if err == datastore.Done {
			break
		}
		if err != nil {
			logrus.Errorf("fetching next Person: %v", err)
			break
		}
		template := "02-01-2006"
		start, _ := time.Parse(template, vars["startDate"])
		end, _ := time.Parse(template, vars["endDate"])

		tr := Trip{
			StartDate: start,
			EndDate:   end,
		}
		u.InvalidDates[0] = tr
		//u.InvalidDates = append(u.InvalidDates, tr)
		_, err = datastore.Put(c, k, u)
		fmt.Fprintf(w, "OK")
	}
}
func addTripHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//username := vars["username"]

	c := appengine.NewContext(r)
	q := datastore.NewQuery("UserData").Ancestor(my_datastore_Key(c))

	var user []User
	if _, err := q.GetAll(c, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	template := "02-01-2006"
	start, _ := time.Parse(template, vars["startDate"])
	end, _ := time.Parse(template, vars["endDate"])

	tr := Trip{
		StartDate: start,
		EndDate:   end,
	}

	user[0].InvalidDates = append(user[0].InvalidDates, tr)
	key := datastore.NewIncompleteKey(c, "UserData", my_datastore_Key(c))
	logrus.Warn(user)
	_, err := datastore.Put(c, key, &user[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("error here")
	} else {
		fmt.Fprintf(w, "trip added")
	}

}
func removeTripHandler(w http.ResponseWriter, r *http.Request)      {}
func getRemaningDaysHandler(w http.ResponseWriter, r *http.Request) {}
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]string)
	m["version"] = "1.0"
	if out, err := json.Marshal(m); err == nil {
		w.Write(out)
	}
}

func startDateHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("UserData").Filter("Username =", mux.Vars(r)["username"])

	var user []User
	if _, err := q.GetAll(c, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%v", user[0].StartDate)
}

func countHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	q, err := datastore.NewQuery("UserData").Ancestor(my_datastore_Key(c)).Filter("Username =", mux.Vars(r)["username"]).Count(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "%d", q)
	}
}

func init() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api", defaultHandler)                                           // set router
	router.HandleFunc("/api/addUser/{username}/{startDate}", addUserHandler)            // set router
	router.HandleFunc("/api/addTrip/{username}/{startDate}/{endDate}", addTripHandler2) // set router
	router.HandleFunc("/api/removeTrip", removeTripHandler)                             // set router
	router.HandleFunc("/api/getRemaningDays", getRemaningDaysHandler)                   // set router
	router.HandleFunc("/api/count/{username}", countHandler)                            // set router
	router.HandleFunc("/api/getStartDate/{username}", startDateHandler)                 // set router
	http.Handle("/", router)

	logrus.Info("Init completed.")
}
