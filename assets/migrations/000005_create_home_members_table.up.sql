CREATE TABLE home_members (
    id BIGSERIAL PRIMARY KEY,
    home_id BIGINT NOT NULL REFERENCES homes(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(home_id, user_id)
);

CREATE INDEX idx_home_members_home_id ON home_members(home_id);
CREATE INDEX idx_home_members_user_id ON home_members(user_id);