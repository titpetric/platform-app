INSERT INTO
  pulse_hourly (user_id, stamp, count)
VALUES
  ('titpetric', strftime('%Y-%m-%d %H:00:00', 'now'), 1)
ON
  CONFLICT(user_id, stamp)
DO
  UPDATE SET count = count + 1;

INSERT INTO
  pulse_daily (user_id, stamp, count)
VALUES
  ('titpetric', strftime('%Y-%m-%d', 'now'), 1)
ON
  CONFLICT(user_id, stamp)
DO
  UPDATE SET count = count + 1;
