package models

type CheckRequest struct {
	JobID string `json:"job_id"`
}

type CheckResponse struct {
	State  string `json:"state"`
	Output string `json:"output"`
}
