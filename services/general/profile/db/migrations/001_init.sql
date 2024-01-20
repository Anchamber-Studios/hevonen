create schema if not exists profile;

create table if not exists profile.profiles (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	identity_id uuid unique,
	first_name text not null,
	middle_name text,
	last_name text not null,
	height smallint,
	weight smallint,
	birthday date,
	created_at timestamp with time zone not null default now(),
	updated_at timestamp with time zone not null default now()
);

create table if not exists profile.contact_informations (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid uuid not null unique default uuid_generate_v4(),
	profile_id bigint not null references profile.profiles(id),
	contact_type text not null,
	contact_value text not null,
	created_at timestamp with time zone not null default now(),
	updated_at timestamp with time zone not null default now()
);

create table if not exists profile.addresses (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	profile_id bigint not null references profile.profiles(id),
	address_line_1 text not null,
	address_line_2 text,
	address_line_3 text,
	city text not null,
	state text not null,
	zip text not null,
	country text not null,
	created_at timestamp with time zone not null default now(),
	updated_at timestamp with time zone not null default now()
);
---- create above / drop below ----
drop table if exists profile.addresses;
drop table if exists profile.contact_informations;
drop table if exists profile.profiles;
drop schema if exists profile;