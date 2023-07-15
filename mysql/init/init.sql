CREATE DATABASE IF NOT EXISTS naro;
USE naro;

CREATE TABLE IF NOT EXISTS users(
    username varchar(255) not null primary key,
    `password` varchar(255) not null
);

CREATE TABLE IF NOT EXISTS todos(
    id int auto_increment primary key,
    name varchar(255) not null,
    creater varchar(255),
    is_finished boolean default 0,
    foreign key (creater) references users(username)
);