create schema kbtu;
create sequence kbtu.tasks_id_seq;
create table kbtu.tasks (id int not null default nextval('kbtu.tasks_id_seq'), task varchar(500) not null, correct_date timestamp not null default now());
insert into kbtu.tasks (task) values ('task1');
