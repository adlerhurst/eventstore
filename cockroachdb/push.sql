INSERT INTO eventstore.events (
    action
    , aggregate
    , revision
    , metadata
    , payload
    , creation_date
    , sequence
) VALUES (
    $1
    , $2
    , $3
    , $4
    , $5
    , DEFAULT
    -- get next sequence of aggregate
    , (SELECT max(sequence+1) FROM eventstore.events e WHERE e.aggregate = $2)
) RETURNING (
    sequence
    , creation_date
);