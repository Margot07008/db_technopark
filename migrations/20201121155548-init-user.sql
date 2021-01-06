-- +migrate Up
create schema if not exists main;
set search_path to main;

create extension if not exists citext;

drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists threads cascade;

create unlogged table users
(
    nickname citext not null unique primary key,
    email citext not null unique,
    fullname varchar(128),
    about text
);

create index if not exists idx_users_nickname on users (nickname);
create index if not exists idx_users_email on users (email);

create unlogged table forums
(
    slug citext not null unique primary key,
    "user" citext not null references users (nickname),
    title text not null ,
    posts integer default 0,
    threads integer default 0
);

create index if not exists idx_forum_user on forums ("user");

create unlogged table threads
(
    id serial primary key,
    title text not null,
    author citext references users (nickname) not null,
    forum citext references forums (slug) not null,
    message text not null ,
    votes integer not null default 0,
    slug citext default null unique,
    created timestamptz default current_timestamp
);

create index if not exists idx_threads_slug on threads (slug);
create index if not exists idx_threads_forum_created on threads (forum, created);
create index if not exists idx_threads_author_forum on threads (author, forum);