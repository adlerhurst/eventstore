CREATE SCHEMA IF NOT EXISTS eventstore;

CREATE TABLE IF NOT EXISTS eventstore.events (
    aggregate STRING[] NOT NULL
    , action STRING[] NOT NULL
    , revision INT2 NOT NULL
    , metadata JSONB
    , payload JSONB
    , sequence INT4 NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL DEFAULT now()

    , CONSTRAINT "primary" PRIMARY KEY (aggregate, sequence)
    , INDEX action_idx (action)
    , INVERTED INDEX metadata_idx (metadata)
    , INDEX sequence_idx (sequence)
);