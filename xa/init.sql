create database xa;

use database xa;

create table customer_account (
  id int(8) unsigned not null,
  balance int(8) unsigned not null
);

insert into customer_account (id, balance) values (1, 1000);

create table merchant_account (
  id int(8) unsigned not null,
  balance int(8) unsigned not null
);

insert into merchant_account (id, balance) values (1, 0);