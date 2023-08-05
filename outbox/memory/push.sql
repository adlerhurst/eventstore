INSERT INTO outbox.events (
    "aggregate"
    , "action"
    , revision
    , metadata
    , payload
    , "sequence"
) VALUES
    {{insertValues}}
RETURNING "sequence", created_at;