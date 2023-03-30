CREATE TABLE IF NOT EXISTS t_voyage (
    id serial PRIMARY KEY,
    user_id INTEGER NOT NULL,
    seed INTEGER NOT NULL,
    level INTEGER NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    records TEXT NOT NULL DEFAULT '',
    start_time TIMESTAMP NOT NULL DEFAULT current_timestamp,
    pass_time TIMESTAMP,
    last_try_time TIMESTAMP NOT NULL DEFAULT current_timestamp
);