package model

type PRstatus string

const (
	CREATED PRstatus = "created"
	MERGED  PRstatus = "merged"
)

type PullRequestShort struct {
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorID        string   `json:"author_id"`
	Status          PRstatus `json:"status"`
}
