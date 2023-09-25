create extension if not exists "uuid-ossp";

create table if not exists company_types
(
    id serial not null constraint company_types_pk primary key,
    name      varchar(50)               not null
);

INSERT INTO company_types(id, name) VALUES (1, 'Corporations');
INSERT INTO company_types(id, name) VALUES (2, 'NonProfit');
INSERT INTO company_types(id, name) VALUES (3, 'Cooperative');
INSERT INTO company_types(id, name) VALUES (4, 'Sole Proprietorship');

create table if not exists companies
(
    id            uuid default   uuid_generate_v4() not null constraint companies_pk primary key,
    name          varchar(15)    unique             not null,
    description   varchar(3000)                     not null,
    employees     int8                              not null,
    registered    boolean        default false      not null,
    type_id       int4           references    company_types(id) on update cascade on delete restrict not null,
    created_at    timestamp      default now()      not null,
    updated_at    timestamp      default now()      not null
);

create table if not exists users
(
    id            uuid default   uuid_generate_v4() not null constraint users_pk primary key,
    email         varchar(100)   unique             not null,
    password      varchar(60)                      not null
);
