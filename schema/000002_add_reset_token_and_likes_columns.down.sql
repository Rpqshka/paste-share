ALTER TABLE users
    DROP COLUMN recovery_code,
    DROP COLUMN expired_at;

ALTER TABLE pastes
    DROP COLUMN paste_date,
    DROP COLUMN likes;
