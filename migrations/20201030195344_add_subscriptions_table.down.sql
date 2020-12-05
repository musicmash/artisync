BEGIN;

DROP TABLE IF EXISTS "subscriptions";
DROP INDEX IF EXISTS "idx_subscriptions_user_name_artist_id";

COMMIT;