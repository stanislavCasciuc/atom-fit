CREATE TABLE IF NOT EXISTS workout_reviews(
	id bigserial PRIMARY KEY,
	user_id bigint,
  workout_id bigint,
	rating int NOT NULL DEFAULT 5,
	title varchar(225) NOT NULL,
	content text NOT NULL,
	created_at TIMESTAMP(0) with time zone DEFAULT NOW(),

	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id),
  CONSTRAINT fk_workout FOREIGN KEY(workout_id) REFERENCES workouts(id),
	CONSTRAINT chk_rating CHECK (rating >= 1 AND rating <= 5)
)
