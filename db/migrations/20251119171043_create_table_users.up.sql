CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

  email TEXT NOT NULL UNIQUE,
  username TEXT UNIQUE,
  full_name TEXT,
  avatar_url TEXT,

  password TEXT NOT NULL,
  role user_role NOT NULL DEFAULT 'MEMBER',
  is_active BOOLEAN NOT NULL DEFAULT true,

  email_verified BOOLEAN NOT NULL DEFAULT false,
  email_verification_token TEXT,
  email_verified_at timestamptz,

  reset_password_token TEXT,
  reset_password_expires timestamptz,
  last_password_change timestamptz,

  last_login_at timestamptz,
  last_login_ip inet,
  last_login_agent TEXT,

  login_count INTEGER NOT NULL DEFAULT 0,
  failed_attempts INTEGER NOT NULL DEFAULT 0,
  locked_until timestamptz,

  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
