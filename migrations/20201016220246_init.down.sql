BEGIN;

DROP TABLE IF EXISTS "artist_one_time_sync_tasks";
DROP TABLE IF EXISTS "artist_daily_sync_tasks";
DROP TABLE IF EXISTS "artist_sync_refresh_tokens";

-- should be equal to schemas from musicmash
DROP TABLE IF EXISTS "subscriptions";
DROP TABLE IF EXISTS "artist_associations";
DROP TABLE IF EXISTS "artists";

DROP TYPE IF EXISTS task_state;

COMMIT;

DROP EXTENSION IF EXISTS "uuid-ossp";