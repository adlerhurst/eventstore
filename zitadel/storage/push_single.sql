INSERT INTO events (
    event_type
    , aggregate_type
    , aggregate_id
    , aggregate_version
    , payload
    , editor_user
    , editor_service
    , resource_owner
    , instance_id
    , creation_date
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, statement_timestamp()) RETURNING creation_date