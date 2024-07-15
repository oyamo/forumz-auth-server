CREATE TABLE IF NOT EXISTS person (
    id UUID not null primary key,
    first_name varchar(32) not null,
    last_name varchar(32) not null,
    email_address varchar(64) unique not null,
    password_hash varchar(256) not null ,
    username varchar(32) unique not null,
    status varchar(12) not null default 'Active',
    dob timestamp,
    datetime_created timestamp not null default current_timestamp,
    last_modified timestamp with time zone  not null default current_timestamp
);