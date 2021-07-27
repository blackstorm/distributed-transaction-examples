CREATE DATABASE `order`;
use order;

create table `order` (
  id int(8) not null primary key auto_increment,
  status varchar(32) not null
)

create table `local_try_log` (
  tx_no varchar(64) unique not null comment 'Global transaction id'
)

create table `local_confirm_log` (
  tx_no varchar(64) unique not null comment 'Global transaction id'
)

create table `local_cancel_log` (
  tx_no varchar(64) unique not null comment 'Global transaction id'
)

# -----------------------------------------

create database account;

use account;

create table account (
  id int(8) not null auto_increment,
  balance int(8) unsigned not null default '0'
)

create table account_trading (
  id int(8) not null auto_increment,
  account_id int(8) not null,
  trading_balance int(8) unsigned not null default '0'
)

create table local_try_log (
  tx_no varchar(64) unique not null comment 'Global transaction id',
)

create table local_confirm_log (
  tx_no varchar(64) unique not null comment 'Global transaction id',
)

create table local_cancel_log (
  tx_no varchar(64) unique not null comment 'Global transaction id',
)