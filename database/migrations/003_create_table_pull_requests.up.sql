CREATE TABLE pull_requests(
    pull_request_id VARCHAR(255) PRIMARY KEY NOT NULL,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users(user_id),
    status VARCHAR(255) NOT NULL,
    created_at TIME NOT NULL,
    merged_at TIME
);
