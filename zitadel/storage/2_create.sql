CREATE TABLE IF NOT EXISTS events2 (
    creation_date TIMESTAMPTZ NOT NULL
    , event_type TEXT NOT NULL
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , aggregate_version TEXT NOT NULL
    , payload JSONB
    , editor_user TEXT NOT NULL 
    , resource_owner TEXT NOT NULL
    , instance_id TEXT NOT NULL
    , region TEXT NULL

    , PRIMARY KEY (instance_id, aggregate_type, aggregate_id, creation_date DESC) --USING HASH
);