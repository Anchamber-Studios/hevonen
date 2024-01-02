create schema if not exists users;

create table if not exists users.users (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid uuid not null unique default uuid_generate_v4(),
	email text not null unique,
	password text,
	email_confirmed boolean not null default false,
	active boolean not null default true,
	created_at timestamp with time zone not null default now(),
	updated_at timestamp with time zone not null default now()
);

create table if not exists users.email_conformation_keys (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id bigint not null,
	key text not null,
	created_at timestamp with time zone not null default now(),
	updated_at timestamp with time zone not null default now(),

	constraint fk_email_conformation_keys_user foreign key (user_id) references users.users(id)
);

create table if not exists users.password_reset_keys (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id bigint not null,
	key text not null,
	created_at timestamp with time zone not null default now(),
	updated_at timestamp with time zone not null default now(),

	constraint fk_password_reset_key_user foreign key (user_id) references users.users(id)
);

create table if not exists users.login_attempts (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id bigint not null,
	ip text not null,
	success boolean not null default false,
	created_at timestamp with time zone not null default now(),
	updated_at timestamp with time zone not null default now(),

	constraint fk_login_attempt_user foreign key (user_id) references users.users(id)
);

create table if not exists users.application_tokens (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  	user_id bigint not null,
	app text not null,
  	token text not null unique,
  	created_at timestamp with time zone not null default now(),
  	updated_at timestamp with time zone not null default now(),

	constraint fk_application_token_user foreign key (user_id) references users.users(id)
);
---- create above / drop below ----
drop table if exists users.application_tokens;
drop table if exists users.email_conformation_keys;
drop table if exists users.password_reset_keys;
drop table if exists users.login_attempts;
drop table if exists users.users;
drop schema if exists users;