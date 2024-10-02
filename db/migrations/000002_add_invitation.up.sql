CREATE TABLE IF NOT EXISTS invitation (
  token bytea PRIMARY KEY,
  user_id bigint NOT NULL
  );
