CREATE TABLE IF NOT EXISTS user_attributes(
	user_id bigint PRIMARY KEY,
	is_male BOOLEAN NOT NULL DEFAULT TRUE,
	height int NOT NULL DEFAULT 175,
	goal varchar(100) NOT NULL DEFAULT 'lose',
	weight_goal REAL NOT NULL DEFAULT 65,
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id),
	CONSTRAINT goal_check CHECK (goal IN ('lose', 'gain', 'maintain'))
);
