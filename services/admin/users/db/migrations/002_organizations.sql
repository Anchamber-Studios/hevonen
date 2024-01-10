create table if not exists users.organizations (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid uuid not null unique default uuid_generate_v4(),
	name text not null unique,
	description text,
	updated_at timestamp with time zone not null default now(),
	created_at timestamp with time zone not null default now()
);
alter table users.users add column organization_id bigint references users.organizations(id);

---- create above / drop below ----
alter table users.users drop column organization_id;
drop table if exists users.organizations;