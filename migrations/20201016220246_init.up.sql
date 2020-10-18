CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

BEGIN;

CREATE TYPE task_state AS ENUM (
    'created',
    'scheduled',
    'in-progress',
    'done',
    'error'
);

CREATE TABLE IF NOT EXISTS "artist_once_sync_tasks" (
    id         uuid PRIMARY KEY       DEFAULT uuid_generate_v4(),
    created_at timestamp    not null  DEFAULT now(),
    updated_at timestamp    not null  DEFAULT now(),
    user_name  varchar(255) not null,
    state      task_state   not null  DEFAULT 'created',
    details    varchar(255)           DEFAULT null
);

CREATE INDEX "idx_artist_once_sync_tasks_state" ON "artist_once_sync_tasks" ("state");
CREATE UNIQUE INDEX "idx_artist_once_sync_tasks_user_name" ON "artist_once_sync_tasks" ("user_name");

CREATE TABLE IF NOT EXISTS "artist_daily_sync_tasks" (
    id         uuid PRIMARY KEY      DEFAULT uuid_generate_v4(),
    created_at timestamp    not null DEFAULT now(),
    updated_at timestamp    not null DEFAULT now(),
    user_name  varchar(255) not null
);

CREATE INDEX "idx_artist_daily_sync_tasks_updated_at" ON "artist_daily_sync_tasks" ("updated_at");
CREATE UNIQUE INDEX "idx_artist_daily_sync_tasks_user_name" ON "artist_daily_sync_tasks" ("user_name");

CREATE TABLE IF NOT EXISTS "artist_sync_refresh_tokens" (
    id         uuid PRIMARY KEY      DEFAULT uuid_generate_v4(),
    created_at timestamp    not null DEFAULT now(),
    expired_at timestamp    not null,
    user_name  varchar(255) not null,
    value      varchar      not null
);

CREATE INDEX "idx_artist_sync_refresh_tokens_expired_at" ON "artist_sync_refresh_tokens" ("expired_at");
CREATE UNIQUE INDEX "idx_artist_sync_refresh_tokens_user_name" ON "artist_sync_refresh_tokens" ("user_name");

COMMIT;