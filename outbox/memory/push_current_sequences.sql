SELECT 
    "sequence"
    , "aggregate" 
FROM 
    outbox.events
WHERE 
    ("sequence", "aggregate") IN (
        SELECT 
            (max("sequence"), "aggregate") 
        FROM 
            outbox.events 
        WHERE 
            {{currentSequencesClauses}}
        GROUP BY 
            "aggregate"
    ) 
FOR UPDATE;
