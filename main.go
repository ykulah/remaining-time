package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	firego "gopkg.in/zabawaba99/firego.v1"
)

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	f := firego.New("https://remaining-time-c9dd7.firebaseio.com/", nil)

	data := map[string]string{}
	data["startDate"] = vars["startDate"]

	v := make(map[string]interface{})
	v[vars["username"]] = data
	if err := f.Update(v); err != nil {
		log.Fatal(err)
	}

}
func addTripHandler(w http.ResponseWriter, r *http.Request) {

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

func dump(w http.ResponseWriter, r *http.Request) {
	f := firego.New("https://remaining-time-c9dd7.firebaseio.com/", nil)
	var v map[string]interface{}
	if err := f.Value(&v); err != nil {
		logrus.Fatal(err)
	}
	if out, err := json.Marshal(v); err == nil {
		w.Write(out)
	} else {
		logrus.Error(err)
	}
}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api", defaultHandler) // set router
	router.HandleFunc("/api/dump", dump)
	router.HandleFunc("/api/addUser/{username}/{startDate}", addUserHandler)           // set router
	router.HandleFunc("/api/addTrip/{username}/{startDate}/{endDate}", addTripHandler) // set router
	router.HandleFunc("/api/removeTrip", removeTripHandler)                            // set router
	router.HandleFunc("/api/getRemaningDays", getRemaningDaysHandler)                  // set router
	http.Handle("/", router)

	port := ":9090"
	logrus.Info("Server started @ http://localhost" + port)
	err := http.ListenAndServe(port, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
