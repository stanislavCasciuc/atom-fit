-- Create the extension and indexes for full-text search
-- Check article: https://niallburkley.com/blog/index-columns-for-like-in-postgres/
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_exercise_description ON exercises USING gin (description gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_workout_description ON workouts USING gin (description gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_workout_name ON workouts USING gin(name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_exercise_name ON exercises USING gin(name gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_exercise_muscles ON exercises USING gin(muscles);

CREATE INDEX IF NOT EXISTS idx_user_username ON users (username);
