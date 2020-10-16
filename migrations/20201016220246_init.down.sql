BEGIN;

DROP TABLE IF EXISTS "artist_once_sync_tasks";
DROP TABLE IF EXISTS "artist_daily_sync_tasks";
DROP TABLE IF EXISTS "artist_sync_refresh_tokens";

COMMIT;

DROP EXTENSION IF EXISTS "uuid-ossp";