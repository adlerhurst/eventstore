CREATE SCHEMA IF NOT EXISTS eventstore;

CREATE TABLE IF NOT EXISTS eventstore.events (
    "aggregate" STRING[] NOT NULL
    , joined_aggregate STRING NOT NULL

    , "action" STRING[] NOT NULL
    , revision INT2 NOT NULL
    , metadata JSONB
    , payload JSONB
    , "sequence" INT4 NOT NULL
    , created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()

    , CONSTRAINT pk PRIMARY KEY ("aggregate", "sequence" DESC)
);
