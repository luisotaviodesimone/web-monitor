-- This is a sample migration.
create table ping_logs (
  id serial primary key,
  host text,
  latency_avg_ms int,
  loss_rate_percent int,
  ping_count int,
  created_at timestamp
  with
    time zone default now ()
);

create table http_logs (
  id serial primary key,
  host text,
  latency_avg_ms int,
  http_count int,
  created_at timestamp
  with
    time zone default now ()
);

---- create above / drop below ----
drop table ping_logs;

drop table http_logs;
