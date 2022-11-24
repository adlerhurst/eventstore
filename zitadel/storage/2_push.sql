INSERT INTO events2 (
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
) VALUES %s RETURNING creation_date