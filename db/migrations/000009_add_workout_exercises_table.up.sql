CREATE TABLE IF NOT EXISTS workout_exercises(
  workout_id bigint,
  exercise_id bigint,
  duration int NOT NULL,
  PRIMARY KEY (workout_id, exercise_id),
  CONSTRAINT workout_fk FOREIGN KEY(workout_id) REFERENCES workouts(id),
  CONSTRAINT exercise_fk FOREIGN KEY(exercise_id) REFERENCES exercises(id)
);
