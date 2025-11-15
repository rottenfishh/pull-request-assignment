CREATE TABLE pr_reviewers(
    pull_request_id VARCHAR(255) NOT NULL REFERENCES pull_requests(pull_request_id),
    reviewer_id VARCHAR(255) REFERENCES users(user_id),
    PRIMARY KEY (pull_request_id, reviewer_id)
)