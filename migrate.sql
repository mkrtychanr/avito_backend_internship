create table client (
  id varchar(128) primary key not null,
  balance varchar(128)
);

create table reserve (
  id bigserial primary key not null,
  client_id varchar(128) references client(id),
  service_id varchar(128) not null,
  order_id varchar(128) not null,
  price varchar(128) not null
);

create table report (
  id bigserial primary key not null,
  client_id varchar(128) references client(id),
  service_id varchar(128) not null,
  order_id varchar(128) not null,
  price varchar(128) not null
);