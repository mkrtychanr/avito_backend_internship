create table client (
  id bigserial primary key not null,
  balance decimal
);

create table reserve (
  id bigserial primary key not null,
  client_id bigserial references client(id),
  service_id bigserial not null,
  order_id bigserial not null,
  price decimal not null
);

create table report (
  id bigserial primary key not null,
  client_id bigserial references client(id),
  service_id bigserial not null,
  order_id bigserial not null,
  price decimal not null
);