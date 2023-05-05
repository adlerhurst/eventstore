SELECT 
    "sequence"
    , "aggregate" 
FROM 
    eventstore.events
WHERE 
    ("sequence", "aggregate") IN (
        SELECT 
            (max("sequence"), "aggregate") 
        FROM 
            eventstore.events 
        WHERE 
            {{currentSequencesClauses}}
        GROUP BY 
            "aggregate"
    ) 
FOR UPDATE;
