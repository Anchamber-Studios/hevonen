create schema if not exists members;
create table if not exists members.members (
    id serial primary key,
    club_id integer not null,
    email text not null unique,
    first_name text not null,
    middle_name text not null,
    last_name text not null,
    height integer not null,
    weight integer not null,
    phone text,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

---- create above / drop below ----

drop table if exists members.members;
drop schema if exists members;