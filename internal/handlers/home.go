package handlers

import (
	"encoding/json"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"msg": "Hello from freight!"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonEncoder := json.NewEncoder(w)

	err := jsonEncoder.Encode(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
