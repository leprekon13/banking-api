CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(100) NOT NULL UNIQUE,
                       email VARCHAR(100) NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
                          id SERIAL PRIMARY KEY,
                          user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                          balance DECIMAL(12, 2) DEFAULT 0.00,
                          created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cards (
                       id SERIAL PRIMARY KEY,
                       account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
                       card_number VARCHAR(16) NOT NULL UNIQUE,
                       expiration_date DATE,
                       cvv TEXT NOT NULL,
                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              from_account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
                              to_account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
                              amount DECIMAL(12, 2) NOT NULL,
                              transaction_type VARCHAR(50),
                              created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE credits (
                         id SERIAL PRIMARY KEY,
                         user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                         amount DECIMAL(12, 2) NOT NULL,
                         interest_rate DECIMAL(5, 2) NOT NULL,
                         start_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         months INTEGER NOT NULL,
                         created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payment_schedules (
                                   id SERIAL PRIMARY KEY,
                                   credit_id INTEGER REFERENCES credits(id) ON DELETE CASCADE,
                                   amount DECIMAL(12, 2) NOT NULL,
                                   due_date TIMESTAMPTZ NOT NULL,
                                   paid BOOLEAN NOT NULL DEFAULT FALSE,
                                   created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
