SELECT
    e."aggregate"
    , e."action"
    , e.revision
    , e.payload
    , e."sequence"
    , e.created_at
    , e.position
FROM
    eventstore.events e
JOIN
    eventstore.actions a
    ON
        e."aggregate" = a."aggregate"
        AND e."sequence" = a."sequence"
{{.Where}}
ORDER BY
    position
    , in_tx_order
{{.Limit}}