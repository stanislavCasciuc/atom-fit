CREATE TABLE IF NOT EXISTS workouts (
	id SERIAL PRIMARY KEY,
	user_id bigint not null,
	name VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
	tutorial_link TEXT NOT NULL,
	created_at TIMESTAMP(0) with time zone DEFAULT NOW(),

	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);
