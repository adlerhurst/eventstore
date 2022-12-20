INSERT INTO events4 (
    aggregate_id
    , event
    , creation_date
) VALUES %s RETURNING creation_date