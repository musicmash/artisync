BEGIN;

DROP TABLE IF EXISTS "artist_once_sync_tasks";
DROP TABLE IF EXISTS "artist_daily_sync_tasks";
DROP TABLE IF EXISTS "artist_sync_refresh_tokens";

-- should be equal to schemas from musicmash
DROP TABLE IF EXISTS "subscriptions";
DROP TABLE IF EXISTS "spotify_artist_associations";
DROP TABLE IF EXISTS "spotify_artists";

COMMIT;

DROP EXTENSION IF EXISTS "uuid-ossp";