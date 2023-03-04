CREATE TABLE IF NOT EXISTS events4 (
    aggregate_id TEXT NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL DEFAULT now()

    , event JSONB NOT NULL
    , region TEXT NULL

    , PRIMARY KEY (aggregate_id, creation_date)
    , INVERTED INDEX search (event)
);