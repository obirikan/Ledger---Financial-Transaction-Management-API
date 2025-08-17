CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(50) NOT NULL DEFAULT 'user'
);

CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  balance NUMERIC NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  from_account_id INT REFERENCES accounts(id),
  to_account_id INT REFERENCES accounts(id),
  amount NUMERIC NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);
