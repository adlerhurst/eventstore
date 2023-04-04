INSERT INTO eventstore.events (
    action
    , aggregate
    , revision
    , metadata
    , payload
) VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5
);