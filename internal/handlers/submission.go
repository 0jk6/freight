package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0jk6/freight/internal/jobs"
	"github.com/0jk6/freight/internal/models"
)

// submission handler
// users submit the code, language
// this will push those to the queue
func SubmissionHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//pull data from request body
	defer r.Body.Close()
	submissionRequest := models.SubmissionRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&submissionRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//we got the request, now push it to the queue
	//we will return the job_id to the user
	resp["job_id"] = jobs.Push(submissionRequest)

	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
