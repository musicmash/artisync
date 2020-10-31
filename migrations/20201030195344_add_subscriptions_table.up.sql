BEGIN;

CREATE TABLE IF NOT EXISTS spotify_subscriptions (
    id         bigserial PRIMARY KEY,
    created_at timestamp    NOT NULL DEFAULT now(),
    user_name  varchar(255) NOT NULL,
    artist_id  bigint       NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES spotify_artists(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_subscriptions_user_name_artist_id ON spotify_subscriptions (user_name, artist_id);

COMMIT;