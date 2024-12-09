CREATE TABLE IF NOT EXISTS finished_workouts (
    date timestamp(0) with time zone default NOW() PRIMARY KEY,
    workout_id INT NOT NULL,
    user_id INT NOT NULL,
    duration INT NOT NULL DEFAULT 0,
    FOREIGN KEY (workout_id) REFERENCES workouts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
