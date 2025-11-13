package model

type PRstatus string

const (
	CREATED PRstatus = "created"
	MERGED  PRstatus = "merged"
)

type PullRequestShort struct {
	PullRequestId   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorId        string   `json:"author_id"`
	Status          PRstatus `json:"status"`
}
