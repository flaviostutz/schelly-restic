//This is a hard fork from the great job done by 
//http://github.com/yp-engineering/rbd-docker-plugin
package main

import (
	"flag"
	"strings"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
	"github.com/Sirupsen/logrus"
)

const VERSION = "1.0.0-beta"

type Options struct {
	listenPort int
	listenIp string
}

var options = new(Options)

func main() {
	listenPort    := flag.Int("listen-port", 8080, "REST API server listen port")
	listenIp      := flag.Int("listen-ip", "", "REST API server listen ip address")
	logLevel      := flag.String("log-level", "info", "debug, info, warning or error")
	flag.Parse()

	switch *logLevel {
		case "debug":
			logrus.SetLevel(logrus.DebugLevel)
			break;
		case "warning":
			logrus.SetLevel(logrus.WarnLevel)
			break;
		case "error":
			logrus.SetLevel(logrus.ErrorLevel)
			break;
		default:
			logrus.SetLevel(logrus.InfoLevel)
	}

	options.listenPort = retentionParams(*listenPort)
	options.listenIp = retentionParams(*listenIp)

	logrus.Infof("====Starting Backy2 REST server %s====", VERSION)

	router := mux.NewRouter()
	router.HandleFunc("/backups", GetBackups).Methods("GET")
	router.HandleFunc("/backups/{id}", CreateBackup).Methods("POST")
	router.HandleFunc("/backups/{id}", DeleteBackup).Methods("DELETE")
    http.ListenAndServe(options.listenIp + ":" + options.listenPort, router)

}

func GetBackups(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range people {
    if item.ID == params["id"]
}

func CreateBackup(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var person Person
    _ = json.NewDecoder(r.Body).Decode(&person)
    person.ID = params["id"]
    people = append(people, person)
    json.NewEncoder(w).Encode(people)
}
func DeleteBackup(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, item := range people {
        if item.ID == params["id"] {
            people = append(people[:index], people[index+1]...)
            break
    }
    json.NewEncoder(w).Encode(people)
}
