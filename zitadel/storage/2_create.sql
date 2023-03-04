CREATE TABLE IF NOT EXISTS events2 (
    -- id UUID NOT NULL DEFAULT gen_random_uuid(),
    creation_date TIMESTAMPTZ NOT NULL
    
    , editor_user TEXT NOT NULL 
    , resource_owner TEXT NOT NULL
    , instance_id TEXT NOT NULL
    
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    
    , event_type TEXT NOT NULL
    , event_version TEXT NOT NULL
    , payload JSONB
    
    , region TEXT NULL

    , PRIMARY KEY (aggregate_id, creation_date)
    , INVERTED INDEX event_search (aggregate_type, event_type, payload)
    , INDEX (aggregate_id, event_type, aggregate_type, resource_owner)
);

-- select * from events2 where aggregate_type = 'asdf' and event_type = 'asdf' and payload @> '{"asdf":"asdf"}';