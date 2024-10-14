CREATE TABLE IF NOT EXISTS exercise_likes(
  user_id bigint,
  exercise_id bigint,

  PRIMARY KEY(user_id, exercise_id),
  CONSTRAINT fk_exercise_likes_exercises FOREIGN KEY (exercise_id) REFERENCES exercises(id),
  CONSTRAINT fk_exercise_likes_users FOREIGN KEY (user_id) REFERENCES users(id)
); 
