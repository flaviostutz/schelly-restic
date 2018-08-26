package main

import (
	"flag"
	// "strings"
    // "encoding/json"
	"net/http"
	"regexp"
	"os"
	"fmt"
	"strings"
	"encoding/json"
    "github.com/gorilla/mux"
	"github.com/Sirupsen/logrus"
)

type Options struct {
	listenPort int
	listenIp string
}

type Response struct {
	Id string `json:"id",omitempty`
	Status string `json:"status",omitempty`
	Message string `json:"message",omitempty`
}

var options = new(Options)

func main() {
	listenPort    := flag.Int("listen-port", 8080, "REST API server listen port")
	listenIp      := flag.String("listen-ip", "0.0.0.0", "REST API server listen ip address")
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

	options.listenPort = *listenPort
	options.listenIp = *listenIp

	logrus.Info("====Starting Restic REST server====")

	router := mux.NewRouter()
	router.HandleFunc("/backups", GetBackups).Methods("GET")
	router.HandleFunc("/backups", CreateBackup).Methods("POST")
	router.HandleFunc("/backups/{id}", DeleteBackup).Methods("DELETE")
	listen := fmt.Sprintf("%s:%d", options.listenIp, options.listenPort)
	logrus.Infof("Listening at %s", listen)
	err := http.ListenAndServe(listen, router)
	if err!=nil {
		logrus.Errorf("Error while listening requests: %s", err)
		os.Exit(1)
	}
}

func GetBackups(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("GetBackups r=%s", r)
	result, err := sh("restic", "snapshots", "--json", "--quiet")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if result == "null" {
		result = "{}"
	}
	w.Write([]byte(result))
	logrus.Infof("result: %s", result)
}

func CreateBackup(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("CreateBackup r=%s", r)
	result, err := sh("restic", "backup", "/backup-source")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Debugf("result: %s", result)
	rex, _ := regexp.Compile("snapshot ([0-9a-zA-z]+) saved")
	id := rex.FindStringSubmatch(result)
	if len(id) != 2 {
		http.Error(w, "Couldn't find returned id from response", http.StatusInternalServerError)
		return
	}
	resp := Response {
		Id: id[1],
		Status: "done",
		Message: result,
	}
	w.Header().Set("Content-Type", "application/json")
	err1 := json.NewEncoder(w).Encode(resp)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Debugf("Response sent %s", resp)
}

func DeleteBackup(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("DeleteBackup r=%s", r)
	params := mux.Vars(r)

	res0, err0 := sh("restic", "snapshots", params["id"])
	if err0 != nil {
		http.Error(w, err0.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Debugf("Query snapshots result: %s", res0)
	if strings.Contains(res0, "it is not a snapshot id") || strings.Contains(res0, "0 snapshots") {
		http.Error(w, fmt.Sprintf("Snapshot %s not found",params["id"]), http.StatusNotFound)
		return
	}

	logrus.Debugf("Snapshot exists")

	result, err := sh("restic", "forget", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Debugf("result: %s", result)

	rex, _ := regexp.Compile("removed snapshot ([0-9a-zA-z]+)")
	id := rex.FindStringSubmatch(result)
	if len(id) != 2 {
		http.Error(w, "Couldn't find returned id from response", http.StatusInternalServerError)
		return
	}

	if id[1] != params["id"] {
		logrus.Errorf("Returned id from forget is different from requested. %s != %s", id[1], params["id"])
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	response := Response {
		Id: id[1],
		Status: "deleted",
		Message: result,
	}
	w.Header().Set("Content-Type", "application/json")
	err1 := json.NewEncoder(w).Encode(response)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Debugf("Response sent %s", response)
}
