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
) VALUES %s RETURNING creation_date