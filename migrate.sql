create table client (
  id varchar(128) primary key not null,
  balance varchar(128)
);

create table reserve (
  id bigserial primary key not null,
  client_id varchar(128) references client(id),
  service_id varchar(128) not null,
  order_id varchar(128) not null,
  price varchar(128) not null,
  reserve_time timestamp
);

create table report (
  id bigserial primary key not null,
  client_id varchar(128) references client(id),
  service_id varchar(128) not null,
  order_id varchar(128) not null,
  price varchar(128) not null,
  report_time timestamp
);

create table client_sheet_change(
  id bigserial primary key not null,
  client_id varchar(128) references client(id),
  status int,
  difference varchar(128),
  change_time timestamp
);
