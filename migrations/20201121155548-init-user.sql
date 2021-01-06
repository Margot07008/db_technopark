-- +migrate Up
create schema if not exists main;
set search_path to main;

CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums CASCADE;

create unlogged table users
(
    nickname citext not null unique primary key,
    email citext not null unique,
    fullname varchar(128),
    about text
);

create unlogged table forums
(
    slug citext not null unique primary key,
    "user" citext not null references users (nickname),
    title text not null ,
    posts integer default 0,
    threads integer default 0
);

