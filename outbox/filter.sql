SELECT
    "aggregate"
    , "action"
    , revision
    , metadata
    , payload
    , "sequence"
    , created_at
FROM
    eventstore.events
{{.Where}}
ORDER BY
    created_at
{{.Limit}}