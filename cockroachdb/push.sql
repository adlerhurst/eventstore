WITH computed AS (SELECT hlc_to_timestamp(cluster_logical_timestamp()) created_at, cluster_logical_timestamp() "position")
, input ("aggregate", "action", revision, payload, "sequence", in_tx_order) AS (VALUES
    {{insertValues}}
)
INSERT INTO eventstore.events (
    created_at
    , "position"
    , "aggregate"
    , "action"
    , revision
    , payload
    , "sequence"
    , in_tx_order
) SELECT 
    c.created_at
    , c."position"
    , i."aggregate"
    , i."action"
    , i.revision
    , i.payload
    , i."sequence"
    , i.in_tx_order
FROM 
    input i
    , computed c
RETURNING 
    created_at
    , "position"
;