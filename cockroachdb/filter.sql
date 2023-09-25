SELECT
    "aggregate"
    , "action"
    , revision
    , payload
    , "sequence"
    , created_at
    , position
FROM
    eventstore.events
{{.Where}}
ORDER BY
    position
    , in_tx_order
{{.Limit}}