create schema if not exists club;

create table if not exists club.addresses (
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

create table if not exists club.clubs (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name text not null unique,
	website text,
	address_id bigint,
	created_at timestamp not null default now(),
	updated_at timestamp not null default now(),
	constraint fk_club_address foreign key (address_id) references club.addresses(id)
);

create table if not exists club.members (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email text not null,
    first_name text not null,
    middle_name text not null,
    last_name text not null,
    height integer not null,
    weight integer not null,
    phone text,
    club_id bigint not null,
	address_id bigint,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

	constraint uq_member_email unique (club_id, email),
	constraint fk_member_address foreign key (address_id) references club.addresses(id),
	constraint fk_member_club foreign key (club_id) references club.clubs(id)
);
---- create above / drop below ----

drop table if exists club.members;
drop table if exists club.clubs;
drop table if exists club.addresses;
drop schema if exists club;