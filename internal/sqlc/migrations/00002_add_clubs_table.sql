-- +goose Up
-- +goose StatementBegin
create table if not exists sport_types
(
    id         bigint generated always as identity primary key, 
    name       varchar(255) not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    is_deleted bool not null default false
);

create unique index if not exists sport_types_uk
    on sport_types (name)
    where sport_types.is_deleted = false;

create table if not exists education_levels
(
    id          bigint generated always as identity primary key,
    name        varchar(255) not null,
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now(),
    is_deleted bool         not null default false
);

create unique index if not exists education_levels_uk
    on education_levels (name)
    where education_levels.is_deleted = false;

create table if not exists clubs
(
    id                          bigint generated always as identity primary key,
    name                        varchar(255) not null,
    description                 text not null,
    sport_type_id               bigint not null,
    teacher_id                  bigint not null,
    total_places                int,
    place                       text not null,
    education_level_id          bigint not null,
    required_workout_per_week   int not null,
    created_at                  timestamptz not null default now(),
    updated_at                  timestamptz not null default now(),
    is_deleted                  bool not null default false,

    constraint fk_sport_type foreign key (sport_type_id) references sport_types (id) on delete restrict,
    constraint fk_teacher_id foreign key (teacher_id) references users (id) on delete restrict,
    constraint fk_education_level foreign key (education_level_id) references education_levels (id) on delete restrict
);

create unique index if not exists clubs_uk
    on clubs (name, sport_type_id, teacher_id, education_level_id)
    where clubs.is_deleted = false;

create table if not exists reviews
(
    id              bigint generated always as identity primary key,
    rating          int not null,
    content         text not null,
    creator_id      bigint not null,
    club_id         bigint not null,
    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now(),
    is_deleted      bool not null default false,

    constraint fk_creator_id foreign key (creator_id) references users (id) on delete restrict
);

create unique index if not exists reviews_uk
    on reviews (creator_id, club_id)
    where reviews.is_deleted = false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists sport_types;
drop table if exists education_levels;
drop table if exists clubs;
drop table if exists reviews;
-- +goose StatementEnd
