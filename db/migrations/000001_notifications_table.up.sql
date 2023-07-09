CREATE TYPE STATUS AS ENUM ('pending', 'sent', 'failed');
CREATE TYPE CHANNEL AS ENUM ('email', 'sms', 'slack');

CREATE TABLE notifications (
    id         SERIAL PRIMARY KEY,
    channel    CHANNEL,
    status     STATUS,
    subject    TEXT,
    body       TEXT,
    sender     TEXT,
    recipient  TEXT
);