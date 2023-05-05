CREATE SCHEMA IF NOT EXISTS eventstore;

CREATE TABLE IF NOT EXISTS eventstore.events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    , "aggregate" STRING[] NOT NULL
    , joined_aggregate STRING NOT NULL

    , "action" STRING[] NOT NULL
    , revision INT2 NOT NULL
    , metadata JSONB
    , payload JSONB
    , "sequence" INT4 NOT NULL
    , created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()

    -- , UNIQUE (joined_aggregate, sequence DESC)
    , INVERTED INDEX action_idx ("action")
    , INVERTED INDEX metadata_idx (metadata)
    , INDEX sequence_idx ("sequence")
);
