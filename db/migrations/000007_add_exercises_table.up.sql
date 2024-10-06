CREATE TABLE IF NOT EXISTS exercises (
  id SERIAL PRIMARY KEY,
  user_id bigint NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT not null,
  is_duration BOOLEAN DEFAULT false NOT NULL,
  duration INTEGER DEFAULT 0 NOT NULL,
  tutorial_link TEXT NOT NULL,
  created_at TIMESTAMP(0) with time zone DEFAULT now() NOT NULL,
  muscles TEXT[] NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);
