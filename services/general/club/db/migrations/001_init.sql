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
	description text,
	website text,
	email text,
	phone text,
	address_id bigint,
	created_at timestamp not null default now(),
	updated_at timestamp not null default now(),
	constraint fk_club_address foreign key (address_id) references clubs.addresses(id)
);

create table if not exists clubs.contacts (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	club_id bigint not null,
	identity_id uuid,
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

	constraint fk_contact_club foreign key (club_id) references clubs.clubs(id) on delete cascade
);

create table if not exists clubs.roles (
	name text not null unique PRIMARY KEY,
	description text,
	system_role boolean not null default false, -- system roles are managed by the application itself and cannot be removed or altered

	created_at timestamp not null default now(),
	updated_at timestamp not null default now()
);

create table if not exists clubs.contact_roles (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	contact_id bigint not null,
	role_name text not null,

	created_at timestamp not null default now(),
	updated_at timestamp not null default now(),

	constraint fk_contact_role_contact foreign key (contact_id) references clubs.contacts(id) on delete cascade,
	constraint fk_contact_role_role foreign key (role_name) references clubs.roles(name) on delete cascade
);

create table if not exists clubs.grants (
	entity text not null,
	action text not null,
	
	PRIMARY KEY (entity, action)
);

create table if not exists clubs.role_grants(
	role_name text not null,
	grant_entity text not null,
	grant_action text not null,
	
	PRIMARY KEY (role_name, grant_entity, grant_action),
	constraint fk_role_grants_role foreign key (role_name) references clubs.roles(name) on delete cascade,
	constraint fk_role_grants_grant foreign key (grant_entity, grant_action) references clubs.grants(entity, action) on delete cascade
);

insert into clubs.grants (entity, action) values ('club', 'view') on conflict (entity, action) do nothing;
insert into clubs.grants (entity, action) values ('club', 'edit') on conflict (entity, action) do nothing;
insert into clubs.grants (entity, action) values ('club', 'delete') on conflict (entity, action) do nothing;

insert into clubs.grants (entity, action) values ('club:member', 'view') on conflict (entity, action) do nothing;
insert into clubs.grants (entity, action) values ('club:member', 'edit') on conflict (entity, action) do nothing;
insert into clubs.grants (entity, action) values ('club:member', 'delete') on conflict (entity, action) do nothing;
insert into clubs.grants (entity, action) values ('club:member', 'invite') on conflict (entity, action) do nothing;
insert into clubs.grants (entity, action) values ('club:member', 'list') on conflict (entity, action) do nothing;

insert into clubs.roles (name, description, system_role) values ('admin', 'Club admin', true) on conflict (name) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('admin', 'club', 'view') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('admin', 'club', 'edit') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('admin', 'club', 'delete') on conflict (role_name, grant_entity, grant_action) do nothing;

insert into clubs.roles (name, description, system_role) values ('manager', 'Club manager', true) on conflict (name) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('manager', 'club', 'view') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('manager', 'club', 'edit') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('manager', 'club:member', 'view') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('manager', 'club:member', 'edit') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('manager', 'club:member', 'delete') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('manager', 'club:member', 'invite') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('manager', 'club:member', 'list') on conflict (role_name, grant_entity, grant_action) do nothing;

insert into clubs.roles (name, description, system_role) values ('trainer', 'Club trainer', true) on conflict (name) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('trainer', 'club', 'view') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('trainer', 'club:member', 'view') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('trainer', 'club:member', 'edit') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('trainer', 'club:member', 'delete') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('trainer', 'club:member', 'invite') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('trainer', 'club:member', 'list') on conflict (role_name, grant_entity, grant_action) do nothing;

insert into clubs.roles (name, description, system_role) values ('member', 'Club member', true) on conflict (name) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('member', 'club', 'view') on conflict (role_name, grant_entity, grant_action) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('member', 'club:member', 'view') on conflict (role_name, grant_entity, grant_action) do nothing;

insert into clubs.roles (name, description, system_role) values ('guest', 'Club guest', true) on conflict (name) do nothing;
insert into clubs.role_grants (role_name, grant_entity, grant_action) values ('guest', 'club', 'view') on conflict (role_name, grant_entity, grant_action) do nothing;
---- create above / drop below ----

drop table if exists clubs.contact_roles;
drop table if exists clubs.role_grants;
drop table if exists clubs.grants;
drop table if exists clubs.contacts;
drop table if exists clubs.clubs;
drop table if exists clubs.addresses;
drop table if exists clubs.roles;
drop schema if exists clubs;