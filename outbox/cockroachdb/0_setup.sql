CREATE SCHEMA IF NOT EXISTS outbox2;

CREATE TABLE IF NOT EXISTS outbox2.events (
    "aggregate" STRING[] NOT NULL
    , "action" STRING[] NOT NULL
    , revision INT2 NOT NULL
    , metadata JSONB
    , payload JSONB
    , "sequence" INT4 NOT NULL
    , created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()

    , PRIMARY KEY ("aggregate", "sequence" DESC)
);

CREATE TABLE IF NOT EXISTS outbox2.subscriptions (
    id UUID NOT NULL DEFAULT gen_random_uuid()
    , pattern STRING[] NOT NULL

    , PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS outbox2.outbox (
    "aggregate" STRING[] NOT NULL
    , "sequence" INT4 NOT NULL
    , created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
    , subscription UUID NOT NULL

    , PRIMARY KEY (subscription, "aggregate", "sequence")
    , FOREIGN KEY ("aggregate", "sequence") REFERENCES outbox2.events ON DELETE CASCADE
    , FOREIGN KEY (subscription) REFERENCES outbox2.subscriptions ON DELETE CASCADE
);

select * from outbox2.subscriptions where pattern[0] in ('#', 'user')