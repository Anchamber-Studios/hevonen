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
	profile_id bigint unique not null,
	identity_id uuid unique,
    email text not null,
    first_name text not null,
    middle_name text not null,
    last_name text not null,
    height integer not null,
    weight integer not null,
    phone text,
    club_id bigint not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

	constraint fk_member_club foreign key (club_id) references clubs.clubs(id)
);
---- create above / drop below ----

drop table if exists clubs.members;
drop table if exists clubs.clubs;
drop table if exists clubs.addresses;
drop schema if exists clubs;