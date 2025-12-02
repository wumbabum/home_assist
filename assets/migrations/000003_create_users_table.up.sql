CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    auth0_sub TEXT NOT NULL UNIQUE,  -- Auth0 subject ID
    email TEXT NOT NULL,
    name TEXT NOT NULL,
    picture TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_auth0_sub ON users(auth0_sub);
