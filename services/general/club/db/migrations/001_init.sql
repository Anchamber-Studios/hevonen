CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create schema if not exists clubs;

create table if not exists clubs.addresses (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    address_line_1 text not null,
    address_line_2 text,
	postal_code text not null,
	city text not null,
	state text not null,
	country text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists clubs.clubs (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name text not null unique,
	website text,
	email text,
	phone text,
	address_id bigint,
	created_at timestamp not null default now(),
	updated_at timestamp not null default now(),
	constraint fk_club_address foreign key (address_id) references clubs.addresses(id)
);

create table if not exists clubs.members (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	club_id bigint not null,
	identity_id uuid unique,
    email text not null,
	first_name text,
	middle_name text,
	last_name text,
	birth_date date,
	phone text,
	weight int,
	height int,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

	constraint fk_member_club foreign key (club_id) references clubs.clubs(id)
);

create table if not exists clubs.roles (
	name text not null unique PRIMARY KEY,
	description text,
	system_role boolean not null default false, -- system roles are managed by the application itself and cannot be removed or altered

	created_at timestamp not null default now(),
	updated_at timestamp not null default now()
);

create table if not exists clubs.member_roles (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	member_id bigint not null,
	role_name text not null,

	created_at timestamp not null default now(),
	updated_at timestamp not null default now(),

	constraint fk_member_role_member foreign key (member_id) references clubs.members(id),
	constraint fk_member_role_role foreign key (role_name) references clubs.roles(name)
);

insert into clubs.roles (name, description, system_role) values ('admin', 'Club administrator', true);
insert into clubs.roles (name, description, system_role) values ('manager', 'Club manager', true);
insert into clubs.roles (name, description, system_role) values ('trainer', 'Club trainer', true);
insert into clubs.roles (name, description, system_role) values ('member', 'Club member', true);
insert into clubs.roles (name, description, system_role) values ('guest', 'Club guest', true);

---- create above / drop below ----

drop table if exists clubs.member_roles;
drop table if exists clubs.members;
drop table if exists clubs.clubs;
drop table if exists clubs.addresses;
drop table if exists clubs.roles;
drop schema if exists clubs;