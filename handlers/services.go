package handlers

import (
	"effective-status/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api Env) processServices(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		var newService model.Service
		err := json.NewDecoder(r.Body).Decode(&newService)
		if err != nil {
			return err
		}
		id, err := api.services.Create(newService)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		resp := fmt.Sprintf(`{"id": "%d"}`, id)
		fmt.Fprint(w, resp)
		w.WriteHeader(http.StatusCreated)
	}
	if r.Method == http.MethodGet {
		services, err := api.services.All()
		if err != nil {
			httpErr := fmt.Sprint("Database connection error:", err.Error())
			http.Error(w, httpErr, http.StatusInternalServerError)
			return err
		}
		resp, err := json.Marshal(services)
		if err != nil {
			errMessage := fmt.Sprint("Failed to parse the list of services:", err)
			http.Error(w, errMessage, http.StatusInternalServerError)
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		w.WriteHeader(http.StatusOK)
	}
	return nil
}
