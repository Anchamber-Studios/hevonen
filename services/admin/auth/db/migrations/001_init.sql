create schema if not exists hauth;

create table if not exists hauth.authorizations (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid uuid not null unique default uuid_generate_v4(),
	name text not null unique,
	service_id uuid not null,
	service_name text not null,
	description text
);

create table if not exists hauth.services (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid uuid not null unique default uuid_generate_v4(),
	name text not null unique,
	description text,
	authorizion_endpoint text not null,
	updated_at timestamp with time zone not null default now(),
	created_at timestamp with time zone not null default now()
);

create table if not exists hauth.service_keys (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid uuid not null unique default uuid_generate_v4(),
	service_id bigint not null references hauth.services(id),
	key text not null unique,
	created_at timestamp with time zone not null default now()
);

create table if not exists hauth.groups (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid uuid not null unique default uuid_generate_v4(),
	name text not null unique,
	description text,
	parent bigint references hauth.groups(id),
	updated_at timestamp with time zone not null default now(),
	created_at timestamp with time zone not null default now()
);

create table if not exists hauth.users (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid uuid not null unique default uuid_generate_v4(),
	email text not null unique
);

create table if not exists hauth.authorizations_groups (
	authorization_id bigint not null references hauth.authorizations(id),
	group_id bigint not null references hauth.groups(id),
	updated_at timestamp with time zone not null default now(),
	created_at timestamp with time zone not null default now(),
	primary key (authorization_id, group_id)
);

create table if not exists hauth.groups_users (
	group_id bigint not null references hauth.groups(id),
	user_id bigint not null references users.users(id),
	updated_at timestamp with time zone not null default now(),
	created_at timestamp with time zone not null default now(),
	primary key (group_id, user_id)
);

---- create above / drop below ----
drop table if exists hauth.authorizations;
drop table if exists hauth.services;
drop table if exists hauth.groups;
drop table if exists hauth.users;
drop table if exists hauth.authorizations_groups;
drop table if exists hauth.groups_users;
drop schema if exists hauth;