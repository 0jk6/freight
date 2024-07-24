package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/0jk6/freight/internal/db"
	"github.com/0jk6/freight/internal/models"
)

// container will post to this
func OutputHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//pull data from request body
	defer r.Body.Close()
	outputRequest := models.OutputRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&outputRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//we got the output request, push this output into the database
	pool := db.GetConnectionPool()
	_, err = pool.Exec(context.Background(), "UPDATE submissions SET output = $1 WHERE job_id = $2", outputRequest.Output, outputRequest.JobID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonEncoder := json.NewEncoder(w)
	resp["msg"] = "success"
	err = jsonEncoder.Encode(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
