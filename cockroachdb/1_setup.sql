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
    , tx_ts TIMESTAMPTZ NOT NULL DEFAULT now()
    , stmt_ts TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp()

    , CONSTRAINT pk PRIMARY KEY ("aggregate", "sequence" DESC)
    -- , INVERTED INDEX action_idx ("action")
    , INVERTED INDEX aggregate_sequence_idx ("sequence", "aggregate")
    , INVERTED INDEX aggregate_idx ("aggregate")
    , INDEX sequence_idx (sequence) STORING (aggregate)
    -- , INDEX sequence_idx ("sequence")
);
