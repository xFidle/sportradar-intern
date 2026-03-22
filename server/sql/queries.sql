-- name: ListFilteredEvents :many
SELECT 
    e.event_id, 
    e.start_time,
    e.end_time,
    e.status,
    s.name AS sport_name,
    c.name AS competition_name 
FROM events e 
JOIN competitions c ON c.competition_id = e._competition_id
JOIN sports s ON s.sport_id = c._sport_id 
WHERE 
    e.start_time >= sqlc.arg('start_after')
    AND e.end_time <= sqlc.arg('end_before')
    AND c._sport_id = COALESCE(sqlc.narg('sport_id'), c._sport_id) 
    AND e._competition_id = COALESCE(sqlc.narg('competition_id'), e._competition_id) 
    AND sqlc.narg('team_ids')::int[] IS NULL OR EXISTS 
      (SELECT 1
      FROM participants p
      WHERE p._event_id = e.event_id
        AND p._team_id = ANY(sqlc.narg('team_ids')::int[]));

-- name: ListEventTeams :many 
SELECT 
    p._event_id,
    t.team_id,
    t.name,
    t.abbreviation,
    t.logo_path
FROM participants p
JOIN teams t ON t.team_id =  p._team_id
WHERE p._event_id = ANY(sqlc.arg('event_ids')::int[]);
     
