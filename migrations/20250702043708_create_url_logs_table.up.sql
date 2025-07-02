CREATE TABLE url_logs (
    id UUID PRIMARY KEY,
    url_id UUID REFERENCES monitored_urls(id) ON DELETE CASCADE,
    status_code INT NOT NULL,
    response_time_ms INT NOT NULL,
    checked_at TIMESTAMP DEFAULT NOW(),
    is_up BOOLEAN NOT NULL
);