CREATE TABLE users(
	id serial PRIMARY KEY,
	name VARCHAR(50) UNIQUE NOT NULL,
	password VARCHAR(256) NOT NULL,
	salt VARCHAR(32) NOT NULL
);
