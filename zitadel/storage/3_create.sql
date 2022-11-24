CREATE TABLE IF NOT EXISTS events3 (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    creation_date TIMESTAMPTZ NOT NULL DEFAULT now()
    , event_type TEXT NOT NULL
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , aggregate_version TEXT NOT NULL
    , payload JSONB
    , editor_user TEXT NOT NULL 
    , editor_service TEXT NOT NULL
    , resource_owner TEXT NOT NULL
    , instance_id TEXT NOT NULL
    , region TEXT NULL

    , PRIMARY KEY (id, aggregate_type, aggregate_id) --USING HASH
);