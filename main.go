package main

import (
	_ "encoding/json"
	_ "flag"
	_ "fmt"
	_ "io/ioutil"
	"log"
	_ "net/http"
	_ "os"
	_ "regexp"
	_ "strings"

	_ "github.com/Sirupsen/logrus"
	_ "github.com/flaviostutz/schelly-webhook/schellyhook"
	_ "github.com/gorilla/mux"
)

func main() {
	log.Print("Should not start this class.")
}
