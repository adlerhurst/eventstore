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

    , PRIMARY KEY ("aggregate", "sequence" DESC)
);

CREATE TABLE IF NOT EXISTS eventstore.actions (
    id UUID NOT NULL DEFAULT gen_random_uuid()
    
    , "action" TEXT
    , parent UUID

    , PRIMARY KEY (id)
    , FOREIGN KEY (parent) REFERENCES eventstore.actions ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS eventstore.action_closures (
    parent UUID NOT NULL
    , child UUID NOT NULL
    , depth INT2 NOT NULL

    , FOREIGN KEY (parent) REFERENCES eventstore.actions ON DELETE CASCADE
    , FOREIGN KEY (child) REFERENCES eventstore.actions ON DELETE CASCADE
)