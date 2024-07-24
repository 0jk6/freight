package models

type OutputRequest struct {
	JobID  string `json:"job_id"`
	Output string `json:"output"`
}
