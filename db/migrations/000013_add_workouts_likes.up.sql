CREATE TABLE IF NOT EXISTS workout_likes (
  workout_id bigint,
  user_id bigint,

  PRIMARY KEY(user_id, workout_id),
	CONSTRAINT fk_workouts_likes_workouts FOREIGN KEY (workout_id) REFERENCES workouts(id),
	CONSTRAINT fk_workouts_likes_users FOREIGN KEY (user_id) REFERENCES users(id)
);
