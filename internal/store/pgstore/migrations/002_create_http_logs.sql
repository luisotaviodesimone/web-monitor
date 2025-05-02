-- This is a sample migration.
create table http_logs (
  id serial primary key,
  latency_avg_ms integer,
  http_count integer,
  created_at timestamp
  with
    time zone default now (),
);

---- create above / drop below ----
drop table http_logs;
