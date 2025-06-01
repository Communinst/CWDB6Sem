create table users
(
    user_id serial primary key,
    name varchar(64) not null unique
);
