CREATE TABLE IF NOT EXISTS user_bio (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL,
  fullname TEXT NOT NULL,
  location TEXT NOT NULL,
  bio TEXT NOT NULL,
  web TEXT NOT NULL,
  picture TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,

  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS user_id_idx ON user_bio(user_id);
