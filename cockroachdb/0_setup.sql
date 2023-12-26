CREATE SCHEMA IF NOT EXISTS eventstore;

CREATE TABLE IF NOT EXISTS eventstore.events (
    id UUID NOT NULL DEFAULT gen_random_uuid()

    , "aggregate" TEXT[] NOT NULL
    , revision INT2 NOT NULL
    , payload JSONB
    , "sequence" INT4 NOT NULL
    , created_at TIMESTAMPTZ NOT NULL
    , "position" DECIMAL NOT NULL
    , in_tx_order INT4 NOT NULL

    , action TEXT[] NOT NULL
    , action_depth INT2 AS (array_length(action, 1)) STORED

    , PRIMARY KEY ("aggregate", "sequence")
);

CREATE TABLE IF NOT EXISTS eventstore.actions (
    "event" UUID
    , "action" TEXT
    , depth INT2

    , PRIMARY KEY (event, action, depth)
    , FOREIGN KEY ("event") REFERENCES eventstore.events ON DELETE CASCADE
    , INDEX search ("action", depth)
);
