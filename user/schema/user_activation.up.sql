-- Add email activation columns to user_auth.
--
-- activated_at IS NULL ⇒ user cannot log in while EmailActivationEnabled is on.
-- activation_token is opaque random emitted to the user via email; cleared on activate.
-- activation_sent_at lets future rate-limiters reject rapid resends.
--
-- The backfill marks every pre-existing row as activated so the migration is
-- safe to apply on a populated database (users created before this feature
-- existed should not be locked out).
ALTER TABLE user_auth ADD COLUMN activated_at DATETIME;
ALTER TABLE user_auth ADD COLUMN activation_token TEXT NOT NULL DEFAULT '';
ALTER TABLE user_auth ADD COLUMN activation_sent_at DATETIME;

UPDATE user_auth SET activated_at = CURRENT_TIMESTAMP WHERE activated_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_auth_activation_token ON user_auth(activation_token);
