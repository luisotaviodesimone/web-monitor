-- This is a sample migration.
create table ping_logs (
  id serial primary key,
  latency_avg_ms integer,
  loss_rate_percent integer,
  ping_count integer,
  created_at timestamp
  with
    time zone default now (),
);

---- create above / drop below ----
drop table ping_logs;
