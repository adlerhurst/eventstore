INSERT INTO events3 (
    event_type
    , aggregate_type
    , aggregate_id
    , aggregate_version
    , payload
    , editor_user
    , resource_owner
    , instance_id
) VALUES %s RETURNING id, creation_date