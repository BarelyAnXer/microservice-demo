-- Create user and database
CREATE USER docker;
CREATE DATABASE docker;
GRANT ALL PRIVILEGES ON DATABASE docker TO docker;

-- Connect to the docker database
\c docker

-- Create user table
CREATE TABLE IF NOT EXISTS "user" (
  user_id SERIAL PRIMARY KEY,
  username VARCHAR(255),
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  gender VARCHAR(10),
  password VARCHAR(50),
  status INT
);