CREATE SCHEMA IF NOT EXISTS eventstore;

DROP TABLE IF EXISTS eventstore.events CASCADE;
CREATE TABLE IF NOT EXISTS eventstore.events (
    id UUID NOT NULL DEFAULT gen_random_uuid()
    -- id INT2 NOT NULL

    , "aggregate" TEXT[] NOT NULL
    , revision INT2 NOT NULL
    , payload JSONB
    , "sequence" INT4 NOT NULL
    , created_at TIMESTAMPTZ NOT NULL
    , "position" DECIMAL NOT NULL
    , in_tx_order INT4 NOT NULL

    , action TEXT[] NOT NULL
    , action_depth INT2 AS (array_length(action, 1)) STORED

    , PRIMARY KEY (id)
    , UNIQUE ("aggregate", "sequence")
    -- , INDEX search (id, action_count) STORING (revision, payload, created_at, "position", in_tx_order)
);

DROP TABLE IF EXISTS eventstore.actions CASCADE;
CREATE TABLE eventstore.actions (
    "event" UUID
    -- "event" INT2
    , "action" TEXT
    , depth INT2

    , PRIMARY KEY (event, action, depth)
    , FOREIGN KEY ("event") REFERENCES eventstore.events ON DELETE CASCADE
    , INDEX search ("action", depth)
);

-- TRUNCATE eventstore.events CASCADE;
-- INSERT INTO eventstore.events VALUES
-- (1, '{"user", "1"}', 1, '{"username":"hursti"}', 1, now(), 0, 0, '{"user", "1", "added"}')
-- , (2, '{"user", "2"}', 1, '{"username":"hursti2"}', 1, now(), 0, 0, '{"user", "2", "added"}')
-- , (3, '{"user", "3"}', 1, '{"username":"hursti3"}', 1, now(), 0, 0, '{"user", "3", "added"}')
-- , (4, '{"user", "2"}', 1, NULL, 2, now(), 0, 0, '{"user", "2", "removed"}')
-- , (5, '{"user", "3"}', 1, '{"firstName":"adler"}', 2, now(), 0, 0, '{"user", "3", "firstName", "set"}')
-- ;

-- TRUNCATE eventstore.actions CASCADE;
-- INSERT INTO eventstore.actions VALUES
-- (1, 'user', 0)
-- , (1, '1', 1)
-- , (1, 'added', 2)
-- , (2, 'user', 0)
-- , (2, '2', 1)
-- , (2, 'added', 2)
-- , (3, 'user', 0)
-- , (3, '3', 1)
-- , (3, 'added', 2)
-- , (4, 'user', 0)
-- , (4, '2', 1)
-- , (4, 'removed', 2)
-- , (5, 'user', 0)
-- , (5, '3', 1)
-- , (5, 'firstName', 2)
-- , (5, 'set', 3)
-- ;


-- -- 1
-- -- user.1.added
-- select * from eventstore.events where id in (
--     select a1.event from eventstore.actions a1 
--     join eventstore.actions a2 on a1.event = a2.event and a1.action = 'user' and a1.depth = 0 and a2.action = '1' and a2.depth = 1
--     join eventstore.actions a3 on a1.event = a3.event and a1.action = 'user' and a1.depth = 0 and a3.action = 'added' and a3.depth = 2
-- )
-- AND action_count = 3
-- ;

-- -- 1,2,3
-- -- user.*.added
-- select * from eventstore.events where id in (
--     select a1.event from eventstore.actions a1 
--     join eventstore.actions a3 on a1.event = a3.event and a1.action = 'user' and a1.depth = 0 and a3.action = 'added' and a3.depth = 2
-- )
-- AND action_count = 3
-- ;

-- -- 1,2,3,4,5
-- -- user.#
-- select * from eventstore.events where id in (
--     select a1.event from eventstore.actions a1 
--     where a1.action = 'user' and a1.depth = 0
-- )
-- AND action_count >= 2
-- ;

-- -- 1,2,3,4
-- -- user.*.*
-- select * from eventstore.events where id in (
--     select a1.event from eventstore.actions a1 
--     where a1.action = 'user' and a1.depth = 0
-- )
-- AND action_count = 3
-- ;

-- -- 4
-- -- user.*.removed
-- select * from eventstore.events where id in (
--     select a1.event from eventstore.actions a1 
--     join eventstore.actions a2 on a1.event = a2.event and a1.action = 'user' and a1.depth = 0 and a2.action = 'removed' and a2.depth = 2
-- )
-- AND action_count = 3
-- ;

-- -- 5
-- -- user.*.firstName.#
-- select * from eventstore.events where id in (
--     select a1.event from eventstore.actions a1 
--     join eventstore.actions a2 on a1.event = a2.event and a1.action = 'user' and a1.depth = 0 and a2.action = 'firstName' and a2.depth = 2
-- )
-- AND action_count >= 4
-- ;

-- -- 1,2
-- -- user.1.added || user.2.added || user.2.#
-- select * 
-- from eventstore.events 
-- where 
--     (
--         id in (
--             select a1.event from eventstore.actions a1 
--             join eventstore.actions a2 on a1.event = a2.event and a1.action = 'user' and a1.depth = 0 and a2.action = '1' and a2.depth = 1
--             join eventstore.actions a3 on a1.event = a3.event and a1.action = 'user' and a1.depth = 0 and a3.action = 'added' and a3.depth = 2
--         )
--         AND action_count = 3
--     )
--     OR (
--         id in (
--             select a1.event from eventstore.actions a1 
--             join eventstore.actions a2 on a1.event = a2.event and a1.action = 'user' and a1.depth = 0 and a2.action = '2' and a2.depth = 1
--             join eventstore.actions a3 on a1.event = a3.event and a1.action = 'user' and a1.depth = 0 and a3.action = 'added' and a3.depth = 2
--         )
--         AND action_count = 3
--     ) 
--     OR (
--         id in (
--             select a1.event from eventstore.actions a1 
--             join eventstore.actions a2 on a1.event = a2.event and a1.action = 'user' and a1.depth = 0 and a2.action = '2' and a2.depth = 1
--         )
--         AND action_count >= 3
--     ) 
-- ;
