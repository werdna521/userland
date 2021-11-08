CREATE TABLE IF NOT EXISTS "audit_logs" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL,
  remote_ip TEXT NOT NULL,

  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
