ALTER TABLE users
    ADD COLUMN recovery_code VARCHAR(6),
    ADD COLUMN expired_at TIMESTAMP;

ALTER TABLE pastes
    ADD COLUMN paste_date TIMESTAMP,
    ADD COLUMN likes INT;
