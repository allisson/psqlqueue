CREATE TABLE IF NOT EXISTS queues(
    id VARCHAR PRIMARY KEY NOT NULL,
    ack_deadline_seconds INT NOT NULL,
    message_retention_seconds INT NOT NULL,
    delivery_delay_seconds INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS messages(
    id VARCHAR PRIMARY KEY NOT NULL,
    queue_id VARCHAR NOT NULL,
    body VARCHAR NOT NULL,
    label VARCHAR,
    attributes JSONB,
    delivery_attempts INT NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (queue_id) REFERENCES queues (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS messages_queue_id_idx ON messages (queue_id);
CREATE INDEX IF NOT EXISTS messages_label_idx ON messages (label);
CREATE INDEX IF NOT EXISTS messages_expired_at_idx ON messages USING BRIN (expired_at);
CREATE INDEX IF NOT EXISTS messages_scheduled_at_idx ON messages USING BRIN (scheduled_at);
CREATE INDEX IF NOT EXISTS messages_expired_at_scheduled_at_idx ON messages USING BRIN (expired_at, scheduled_at);
