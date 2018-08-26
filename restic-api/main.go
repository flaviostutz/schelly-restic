package main

import (
	"flag"
	// "strings"
    // "encoding/json"
	"net/http"
	"os"
	"fmt"
	"encoding/json"
    "github.com/gorilla/mux"
	"github.com/Sirupsen/logrus"
)

type Options struct {
	listenPort int
	listenIp string
}

type Response struct {
	id string
	status string
	message string
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
	router.HandleFunc("/backups/{id}", CreateBackup).Methods("POST")
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
	logrus.Infof("result: %s", result)
	response := Response {
		id: result,
		status: "done",
		message: "Backup done",
	}
	err1 := json.NewEncoder(w).Encode(&response)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Debugf("Response sent %s", response)
}

func DeleteBackup(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("DeleteBackup r=%s", r)
	params := mux.Vars(r)
	result, err := sh("restic", "forget", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Debugf("result: %s", result)
	response := Response {
		id: result,
		status: "done",
		message: "Backup done",
	}
	err1 := json.NewEncoder(w).Encode(&response)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Debugf("Response sent %s", response)
}
