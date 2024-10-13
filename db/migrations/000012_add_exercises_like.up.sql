CREATE TABLE IF NOT EXISTS exercise_likes(
  user_id bigint,
  exercise_id bigint,

  PRIMARY KEY(user_id, exercise_id)
); 
