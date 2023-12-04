-- Your SQL goes here
CREATE TABLE users (
    id integer primary key not null,
    login text not null unique,
    email text not null unique,
    hash text not null,
    isadmin INTEGER
)