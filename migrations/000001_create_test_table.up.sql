create table test_table (
    id serial primary key,
    name varchar(255) not null,
    created_at timestamp default current_timestamp
);