-- +migrate Up
create schema if not exists main;
set search_path to main;

CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS users CASCADE;

create unlogged table users
(
    nickname citext not null unique primary key,
    email citext not null unique,
    fullname varchar(128),
    about text
);

