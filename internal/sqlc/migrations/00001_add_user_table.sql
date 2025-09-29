-- +goose Up
-- +goose StatementBegin
create table if not exists group_types
(
    id         bigint generated always as identity primary key,
    name       varchar(128) not null,
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now(),
    is_deleted bool         not null default false
);

create unique index if not exists group_types_name_uk
    on group_types (name)
    where group_types.is_deleted = false;

create table if not exists groups
(
    id              bigint generated always as identity primary key,
    institute       varchar(64) not null,
    enrollment_year integer     not null,
    prefix          varchar(32),
    group_type_id   bigint      not null,
    group_number    smallint,
    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now(),
    is_deleted      bool        not null default false,

    constraint fk_group_type foreign key (group_type_id) references group_types (id) on delete restrict
);

create unique index if not exists groups_uk
    on groups (institute, enrollment_year, prefix, group_type_id, group_number)
    where groups.is_deleted = false;

create table if not exists users
(
    id                  bigint generated always as identity primary key,
    full_name           varchar(255) not null,
    social_network_link varchar(255) not null,
    phone_number        varchar(32)  not null,
    email               varchar(255) not null,
    birth_date          timestamptz  not null,
    role                varchar(128) not null,
    password            bytea        not null,
    group_id            bigint,
    created_at          timestamptz  not null default now(),
    updated_at          timestamptz  not null default now(),
    is_deleted          bool         not null default false,

    constraint fk_group foreign key (group_id) references groups (id) on delete restrict
);

create unique index if not exists users_email_uk
    on users (email)
    where users.is_deleted = false;

create unique index if not exists users_social_network_link_uk
    on users (social_network_link)
    where users.is_deleted = false;

create unique index if not exists users_phone_number_uk
    on users (phone_number)
    where users.is_deleted = false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists groups;
drop table if exists group_types;
drop table if exists users;
-- +goose StatementEnd
