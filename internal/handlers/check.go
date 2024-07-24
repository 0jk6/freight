package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/0jk6/freight/internal/db"
	"github.com/0jk6/freight/internal/models"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//pull data from request body
	jobID := r.URL.Query().Get("job_id")
	if jobID == "" {
		http.Error(w, "job_id is required", http.StatusBadRequest)
		return
	}

	//we got the output request, push this output into the database
	pool := db.GetConnectionPool()
	row := pool.QueryRow(context.Background(), "SELECT output FROM submissions WHERE job_id = $1", jobID)
	var output string
	row.Scan(&output)

	checkResponse := models.CheckResponse{}

	checkResponse.Output = output
	if output == "" {
		checkResponse.State = "PENDING"
	} else {
		checkResponse.State = "SUCCESS"
	}

	jsonEncoder := json.NewEncoder(w)
	err := jsonEncoder.Encode(checkResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
