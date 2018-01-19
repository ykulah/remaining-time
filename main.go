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

const (
	timeFormat = "02-01-2006"
	daysNeeded = 1095
)

// guestbookKey returns the key used for all guestbook entries.
func my_datastore_Key(c context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "my_datastore", "default_dataStore", 0, nil)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c := appengine.NewContext(r)

	timeStr, _ := time.Parse(timeFormat, vars["startDate"])

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
		var u User
		k, err := t.Next(&u)
		if err == datastore.Done {
			break
		}
		if err != nil {
			logrus.Errorf("fetching next Person: %v", err)
			break
		}
		start, _ := time.Parse(timeFormat, vars["startDate"])
		end, _ := time.Parse(timeFormat, vars["endDate"])

		tr := Trip{
			StartDate: start,
			EndDate:   end,
		}

		dummy, _ := time.Parse(timeFormat, "0001-01-01")
		if u.InvalidDates[0].StartDate == dummy {
			u.InvalidDates[0] = tr
		} else {
			u.InvalidDates = append(u.InvalidDates, tr)
		}
		_, err = datastore.Put(c, k, &u)
		if err == nil {
			fmt.Fprintf(w, "%s ", "OK")
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}

	}
}

func removeTripHandler(w http.ResponseWriter, r *http.Request) {
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
		start, _ := time.Parse(timeFormat, vars["startDate"])
		end, _ := time.Parse(timeFormat, vars["endDate"])

		for idx, tr := range u.InvalidDates {
			if tr.StartDate == start && tr.EndDate == end {
				u.InvalidDates = append(u.InvalidDates[:idx], u.InvalidDates[idx+1:]...)
				break
			}
		}
		_, err = datastore.Put(c, k, &u)
		if err == nil {
			fmt.Fprintf(w, "%s ", "Trip deleted.")
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}

	}
}
func getRemaningDaysHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c := appengine.NewContext(r)

	q := datastore.NewQuery("UserData").Filter("Username =", vars["username"])
	t := q.Run(c)
	var invalidDays int
	var requiredDays int
	var duration int
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

		for _, tr := range u.InvalidDates {
			tmp := tr.EndDate.Sub(tr.StartDate)
			invalidDays += int(tmp.Hours() / 24)
		}

		duration = int(time.Since(u.StartDate).Hours() / 24)
		requiredDays = daysNeeded - duration + invalidDays
	}

	fmt.Fprintf(w, "%d", requiredDays)
}
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]string)
	m["version"] = "1.0"
	if out, err := json.Marshal(m); err == nil {
		w.Write(out)
	}
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
		for _, tmp := range u.InvalidDates {
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

	router.HandleFunc("/api", defaultHandler)                                                // DONE
	router.HandleFunc("/api/addUser/{username}/{startDate}", addUserHandler)                 // DONE
	router.HandleFunc("/api/addTrip/{username}/{startDate}/{endDate}", addTripHandler)       // DONE
	router.HandleFunc("/api/removeTrip/{username}/{startDate}/{endDate}", removeTripHandler) // DONE
	router.HandleFunc("/api/getRemainingDays/{username}", getRemaningDaysHandler)            // DONE
	router.HandleFunc("/api/count/{username}", countHandler)                                 // DONE
	router.HandleFunc("/api/getTrips/{username}", getTripsHandler)                           // DONE
	http.Handle("/", router)

	logrus.Info("Init completed.")
}
