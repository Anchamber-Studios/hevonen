create schema if not exists users;

create table if not exists users.users (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid uuid not null unique default uuid_generate_v4(),
	email text not null unique,
	password text,
	created_at timestamp with time zone not null default now(),
	updated_at timestamp with time zone not null default now()
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
drop table if exists users.users;
drop schema if exists users;