CREATE TABLE sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  access_token text,
  refresh_token_hash text NOT NULL,
  jwt_id text NOT NULL,

  ip_address inet,
  user_agent text,
  device_name text,
  location text,

  created_at timestamptz NOT NULL DEFAULT now(),
  last_used_at timestamptz,
  expires_at timestamptz NOT NULL,
  revoked boolean NOT NULL DEFAULT false,
  revoked_at timestamptz,
  revoked_reason revoked_reason,

  fingerprint text,
  is_current boolean DEFAULT false
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_revoked ON sessions(revoked);
CREATE INDEX idx_sessions_jwt_id ON sessions(jwt_id);
