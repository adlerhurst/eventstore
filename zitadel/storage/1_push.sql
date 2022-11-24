INSERT INTO events1 (
    event_type
    , aggregate_type
    , aggregate_id
    , aggregate_version
    , payload
    , editor_user
    , editor_service
    , resource_owner
    , instance_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING creation_date