CREATE TABLE IF NOT EXISTS events (
  id VARCHAR PRIMARY KEY,
  aggregate_id VARCHAR NOT NULL,
  user_id VARCHAR NOT NULL,
  user_name VARCHAR NOT NULL,
  is_snapshot BOOLEAN NOT NULL,
  type TEXT NOT NULL,
  payload JSON NOT NULL,
  created DATETIME NOT NULL,
  metadata JSON,
  version INTEGER NOT NULL
);