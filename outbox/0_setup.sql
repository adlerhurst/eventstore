CREATE SCHEMA IF NOT EXISTS outbox;

CREATE TABLE IF NOT EXISTS outbox.events (
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


CREATE TABLE IF NOT EXISTS outbox.outbox (
    "aggregate" STRING[] NOT NULL
    , "sequence" INT4 NOT NULL
    , reveiver STRING NOT NULL

    , CONSTRAINT pk PRIMARY KEY ("aggregate", "sequence")
    , CONSTRAINT fk_event FOREIGN KEY ("aggregate", "sequence") REFERENCES outbox.events ON DELETE CASCADE
);