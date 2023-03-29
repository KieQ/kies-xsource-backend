CREATE TABLE IF NOT EXISTS user (
    id serial PRIMARY KEY,
    account VARCHAR(255) NOT NULL UNIQUE DEFAULT '',
    password VARCHAR(255) NOT NULL DEFAULT '',
    nick_name VARCHAR(100) NOT NULL DEFAULT '',
    profile VARCHAR(400) NOT NULL DEFAULT '',
    phone VARCHAR(20) NOT NULL DEFAULT '',
    email VARCHAR(30) NOT NULL DEFAULT '',
    gender INTEGER NOT NULL DEFAULT -1,
    self_introduction VARCHAR(500) NOT NULL DEFAULT '',
    create_time TIMESTAMP NOT NULL DEFAULT current_timestamp,
    update_time TIMESTAMP NOT NULL DEFAULT current_timestamp
);