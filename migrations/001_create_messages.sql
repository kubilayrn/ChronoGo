CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    "to" VARCHAR(20) NOT NULL,
    content VARCHAR(320) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'unsent' CHECK (status IN ('unsent', 'sent')),
    sent_at TIMESTAMP,
    message_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status);

CREATE INDEX IF NOT EXISTS idx_messages_sent_at ON messages(sent_at);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_messages_updated_at BEFORE UPDATE ON messages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
