package models

type SubmissionRequest struct {
	Code     string `json:"code"`
	Language string `json:"lang"`
	JobID    string `json:"job_id"`
}

type SubmissionResponse struct {
	JobID  string `json:"job_id"`
	Output string `json:"output"`
}
