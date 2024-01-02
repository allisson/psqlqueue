CREATE TABLE IF NOT EXISTS topics(
    id VARCHAR PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS subscriptions(
    id VARCHAR PRIMARY KEY NOT NULL,
    topic_id VARCHAR NOT NULL,
    queue_id VARCHAR NOT NULL,
    message_filters JSONB,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE,
    FOREIGN KEY (queue_id) REFERENCES queues (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS subscriptions_topic_id_queue_id_idx ON subscriptions (topic_id, queue_id);
