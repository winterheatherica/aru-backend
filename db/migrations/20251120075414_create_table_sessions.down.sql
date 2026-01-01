DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_revoked;
DROP INDEX IF EXISTS idx_sessions_jwt_id;

DROP TABLE IF EXISTS sessions;
