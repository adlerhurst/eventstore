INSERT INTO outbox.outbox (
    "aggregate"
    , "sequence"
    , created_at
    , receiver
) VALUES
    {{insertValues}};