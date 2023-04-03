CREATE TABLE IF NOT EXISTS t_user (
    id serial PRIMARY KEY,
    email VARCHAR(50) NOT NULL UNIQUE DEFAULT '',
    password VARCHAR(255) NOT NULL DEFAULT '',
    nick_name VARCHAR(100) NOT NULL DEFAULT '',
    profile TEXT NOT NULL DEFAULT '',
    gender SMALLINT NOT NULL DEFAULT -1,
    phone VARCHAR(20) NOT NULL DEFAULT '',
    self_introduction TEXT NOT NULL DEFAULT '',
    create_time TIMESTAMP NOT NULL DEFAULT current_timestamp,
    update_time TIMESTAMP NOT NULL DEFAULT current_timestamp
);