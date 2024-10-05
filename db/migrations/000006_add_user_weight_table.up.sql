CREATE TABLE IF NOT EXISTS user_weight(
	user_id bigint,
  weigh_date DATE DEFAULT CURRENT_DATE,
  weight REAL NOT NULL,
  PRIMARY KEY (user_id, weigh_date),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);
