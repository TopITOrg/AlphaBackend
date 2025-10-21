-- +goose Up
-- +goose StatementBegin
create table if not exists attachments
(
    id              bigint generated always as identity primary key,
    filename        varchar(255)    not null,
    type_id         bigint          not null,
    created_at      timestamptz     not null default now(),
    updated_at      timestamptz     not null default now(),
    is_deleted      bool            not null default false,
    constraint fk_attachment_type foreign key (type_id) references attachment_types (id) on delete restrict,
);
create unique index attachments_uk
    on attachments (filename)
    where attachments.is_deleted = false;

create table if not exists attachment_types
(
    id              bigint generated always as identity primary key,
    name            varchar(128)    not null,
    created_at      timestamptz     not null default now(),
    is_deleted      bool            not null default false,
);

create unique index attachment_types_uk
    on attachment_types (name)
    where attachment_types.is_deleted = false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists attachments;
drop table if exists attachment_types;
-- +goose StatementEnd
