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
func addTripHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c := appengine.NewContext(r)

	q := datastore.NewQuery("UserData").Filter("Username =", vars["username"])
	t := q.Run(c)
	for {
		logrus.Warn("-----")
		var u User
		k, err := t.Next(&u)
		if err == datastore.Done {
			logrus.Warn("DONE")
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
		u.InvalidDates = append(u.InvalidDates, tr)
		logrus.Warn(u)
		//u.InvalidDates = append(u.InvalidDates, tr)
		_, err = datastore.Put(c, k, &u)
		if err == nil {
			fmt.Fprintf(w, "%s ", "OK")
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}

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

func getTripsHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("UserData").Ancestor(my_datastore_Key(c)).Filter("Username =", mux.Vars(r)["username"])
	t := q.Run(c)
	ret := ""
	for {
		var u User
		_, err := t.Next(&u)
		if err == datastore.Done {
			break
		}
		if err != nil {
			logrus.Errorf("fetching next Person: %v", err)
			break
		}

		ret += u.Username
		ret += "\n"
		for i, tmp := range u.InvalidDates {
			logrus.Warn(i)
			ret += "start: "
			ret += tmp.StartDate.String()
			ret += " end: "
			ret += tmp.EndDate.String()
			ret += "\n"
		}
	}
	fmt.Fprintf(w, "%s", ret)
}

func init() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api", defaultHandler)                                          // set router
	router.HandleFunc("/api/addUser/{username}/{startDate}", addUserHandler)           // set router
	router.HandleFunc("/api/addTrip/{username}/{startDate}/{endDate}", addTripHandler) // set router
	router.HandleFunc("/api/removeTrip", removeTripHandler)                            // set router
	router.HandleFunc("/api/getRemaningDays", getRemaningDaysHandler)                  // set router
	router.HandleFunc("/api/count/{username}", countHandler)                           // set router
	router.HandleFunc("/api/getStartDate/{username}", startDateHandler)                // set router
	router.HandleFunc("/api/getTrips/{username}", getTripsHandler)                     // set router
	http.Handle("/", router)

	logrus.Info("Init completed.")
}
