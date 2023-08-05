SELECT
    "aggregate"
    , "action"
    , revision
    , metadata
    , payload
    , "sequence"
    , created_at
FROM
    outbox.events
{{.Where}}
ORDER BY
    created_at
{{.Limit}}