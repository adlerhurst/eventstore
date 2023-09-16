CREATE SCHEMA IF NOT EXISTS eventstore;

CREATE TABLE IF NOT EXISTS eventstore.events (
    "aggregate" TEXT[] NOT NULL
    , "action" TEXT[] NOT NULL
    , revision INT2 NOT NULL
    , payload JSONB
    , "sequence" INT4 NOT NULL
    , created_at TIMESTAMPTZ NOT NULL
    , "position" DECIMAL NOT NULL
    , in_tx_order INT4 NOT NULL

    , PRIMARY KEY ("aggregate", "sequence" DESC)
    , INDEX filter_aggregate ("aggregate", "position", in_tx_order) INCLUDE ("action", revision, payload, created_at)
    , INDEX filter_action ("action", "position", in_tx_order) INCLUDE (revision, payload, created_at)
);
