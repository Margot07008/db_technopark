-- +migrate Up
create schema if not exists main;
set search_path to main;

create extension if not exists citext;

drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;

create unlogged table users
(
    nickname citext not null unique primary key,
    email citext not null unique,
    fullname varchar(100),
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

create unlogged table posts
(
    id       serial primary key ,
    forum    citext references forums (slug),
    parent   integer default 0,
    author   citext references users (nickname) not null,
    created  timestamptz default current_timestamp,
    isEdited boolean default false,
    message  text not null,
    thread   integer references threads (id) not null,
    path     integer[] default array []::int[]
);

create index if not exists idx_posts_path_id on posts (id, (path [1]));
create index if not exists idx_posts_path on posts (path);
create index if not exists idx_posts_path_1 on posts ((path [1]));
create index if not exists idx_posts_thread_id on posts (thread, id);
create index if not exists idx_posts_thread on posts (thread);
create index if not exists idx_posts_thread_path_id on posts (thread, path, id);
create index if not exists idx_posts_thread_id_path_parent on posts (thread, id, (path[1]), parent);
create index if not exists idx_posts_author_forum on posts (author, forum);