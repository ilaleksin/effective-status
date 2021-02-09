package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/mux"
)

var services = []Service{
	{
		Name: "Auth service",
		Checks: []HealthCheck{
			{Status: 0, Title: "Ping https", Details: "OK", Priority: 1},
			{Status: 1, Title: "Error Response Rate", Details: "5xx reponse rate is 10% higher", Priority: 2},
			{Status: 0, Title: "NFS Endpoint", Details: "OK", Priority: 1},
		},
		Tags: []string{"Prod", "Auth"},
	},
	{
		Name: "NFS share",
		Checks: []HealthCheck{
			{Status: 0, Title: "Ping", Details: "OK", Priority: 1},
			{Status: 0, Title: "Error Log Rate", Details: "OK", Priority: 3},
		},
		Tags: []string{"Prod", "Infra", "Storage"},
	},
	{
		Name: "CI\\CD framework",
		Checks: []HealthCheck{
			{Status: 0, Title: "Ping", Details: "OK", Priority: 1},
			{Status: 1, Title: "VCS avg response time", Details: "Average response time is 12% higher", Priority: 2},
			{Status: 2, Title: "Service Account Denied logins", Details: "OK", Priority: 3},
		},
		Tags: []string{"Prod", "Build"},
	},
}

func getServiceBoard(w http.ResponseWriter, r *http.Request) {
	var serviceSummary []Service
	for _, service := range services {
		serviceSummary = append(serviceSummary, service.GetShortService())
	}
	resp, err := json.Marshal(serviceSummary)
	if err != nil {
		log.Panic("Failed to parse health checks of services", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getServiceHealth(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	title := variables["service"]
	log.Println("Service name", title)
	service := findService(title)
	resp, err := json.Marshal(service)
	if err != nil {
		err := fmt.Errorf("Failed to parse service object %v", service)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func findService(serviceTitle string) *Service {
	for _, value := range services {
		if value.Name == serviceTitle {
			return &value
		}
	}
	return &Service{}
}

func updateService(w http.ResponseWriter, r *http.Request) {
	var patchService Service
	err := json.NewDecoder(r.Body).Decode(&patchService)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var failedUpdates []string
	var validCheck bool
	service := findService(patchService.Name)
	if reflect.DeepEqual(Service{}, *service) {
		err := fmt.Errorf("Service %v is not supported, please create this service first", patchService.Name)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range services {
		if services[i].Name == patchService.Name {
			for _, patchCheck := range patchService.Checks {
				checks := services[i].Checks
				for k := range checks {
					if checks[k].Title == patchCheck.Title {
						err := checks[k].SetStatus(patchCheck.Status)
						if err != nil {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						validCheck = true
					}
				}
				if !validCheck {
					failMessage := fmt.Sprintf("%v: %v", patchService.Name, patchCheck.Details)
					failedUpdates = append(failedUpdates, failMessage)
				}
				services[i].Checks = checks
			}
		}
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Updated service check successfully"))
	return
}

func initAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/template/status.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Everything is fine")
	w.WriteHeader(http.StatusOK)
	t.Execute(w, services)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintln(w, "hello")
}

func main() {
	logFile, err := os.Create("tmp.log")
	if err != nil {
		fmt.Println("Failed to open a file")
	}
	mw := io.MultiWriter(os.Stdout, logFile)

	logger := log.New(mw, "Logger bruh: ", log.Ldate|log.Lshortfile)
	logMdlw := LoggingMiddleware(logger)
	router := mux.NewRouter()

	router.HandleFunc("/", initAction).Methods("GET")
	router.HandleFunc("/board", getServiceBoard).Methods("GET")
	router.HandleFunc("/check", updateService).Methods("PATCH")
	router.HandleFunc("/services/{service}", getServiceHealth).Methods("GET")

	finalMux := logMdlw(router)
	http.ListenAndServe(":9091", finalMux)
}
