CREATE TABLE IF NOT EXISTS forgot_passwords (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  old_password TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
)