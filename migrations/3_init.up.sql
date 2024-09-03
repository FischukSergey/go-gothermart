alter table users drop column id;
alter table users add column id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY;