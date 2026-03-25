-- name: InsertEvent :one
INSERT INTO events (_competition_id, _venue_id, _stage_id, status, start_time) 
VALUES (sqlc.arg('competition_id'), sqlc.arg('venue_id'), sqlc.arg('stage_id'), sqlc.arg('status'), sqlc.arg('start_time'))
RETURNING event_id;

-- name: InsertParticipants :exec
INSERT INTO participants (_event_id, _team_id) 
SELECT 
    sqlc.arg('event_id'),
    unnest(sqlc.arg('teams_ids')::int[]);

-- name: IsVenueValidForCompetition :one
SELECT EXISTS (
    SELECT 1
    FROM competitions c
    JOIN playgrounds p on p._sport_id = c._sport_id
    WHERE 
        c.competition_id = sqlc.arg('competition_id')
        AND p._venue_id = sqlc.arg('venue_id')
);

-- name: CountValidTeamsForCompetition :one
SELECT COUNT(DISTINCT ct._team_id)
FROM competition_teams ct
WHERE 
    ct._competition_id = sqlc.arg('competition_id')
    AND ct._team_id = ANY(sqlc.arg('team_ids')::int[]);

-- name: GetDetailedEventByID :one
SELECT 
    e.event_id, 
    e.start_time,
    e.end_time,
    e.status,
    s.name AS sport_name,
    c.name AS competition_name,
    c.type AS competition_type,
    v.name AS venue_name
FROM events e
JOIN venues v ON v.venue_id = e._venue_id 
JOIN competitions c ON c.competition_id = e._competition_id
JOIN sports s ON s.sport_id = c._sport_id 
WHERE e.event_id = sqlc.arg('event_id');


-- name: ListDetailedTeamsByEventID :many
SELECT 
    p._event_id,
    t.team_id,
    t.name,
    t.abbreviation,
    t.logo_path,
    ci.name AS city_name,
    co.name AS country_name,
    co.code AS country_code
FROM participants p
JOIN teams t ON t.team_id = p._team_id
JOIN cities ci ON ci.city_id = t._city_id
JOIN countries co ON co.country_id = ci._country_id
WHERE p._event_id = sqlc.arg('event_id');


-- name: ListPlayersByTeamIDs :many
SELECT
    p.player_id,
    p._team_id,
    p.first_name,
    p.last_name,
    co.name AS country_name,
    co.code AS country_code
FROM players p
JOIN countries co ON co.country_id = p._country_id
WHERE p._team_id = ANY(sqlc.arg('team_ids')::int[]);
       

-- name: ListEventsByFilter :many
SELECT 
    e.event_id, 
    e.start_time,
    e.status,
    s.name AS sport_name,
    c.name AS competition_name,
    c.type AS competition_type
FROM events e 
JOIN competitions c ON c.competition_id = e._competition_id
JOIN sports s ON s.sport_id = c._sport_id 
WHERE 
    e.start_time >= sqlc.arg('start_after')
    AND e.start_time <= sqlc.arg('end_before')
    AND c._sport_id = COALESCE(sqlc.narg('sport_id'), c._sport_id) 
    AND sqlc.narg('team_ids')::int[] IS NULL OR EXISTS 
      (SELECT 1 FROM participants p
      WHERE p._event_id = e.event_id
        AND p._team_id = ANY(sqlc.narg('team_ids')::int[]));


-- name: ListTeamsByEventsIDs :many 
SELECT 
    p._event_id,
    t.team_id,
    t.name,
    t.abbreviation,
    t.logo_path
FROM participants p
JOIN teams t ON t.team_id = p._team_id
WHERE p._event_id = ANY(sqlc.arg('event_ids')::int[]);


-- name: ListFinalScoresByEventsIDs :many
SELECT
    p._event_id,
    p._team_id,
    SUM(s.score) as agg_score
FROM participants p
JOIN scores s ON s._participant_id = p.participant_id
WHERE p._event_id = ANY(sqlc.arg('event_ids')::int[])
GROUP BY participant_id;


-- name: ListScoresByEventID :many
SELECT
    p._team_id,
    s.segment,
    s.score
FROM participants p
JOIN scores s ON s._participant_id = p.participant_id
WHERE p._event_id = sqlc.arg('event_id');


-- name: ListSports :many
SELECT * FROM sports;

-- name: ListCompetitionsBySportID :many
SELECT 
    c.competition_id,
    c.name,
    c.type,
    c.logo_path
FROM competitions c
WHERE c._sport_id = sqlc.arg('sport_id');


-- name: ListVenuesBySportID :many
SELECT 
    v.venue_id,
    v.name,
    v.capacity,
    ci.name AS city_name,
    co.name AS country_name,
    co.code AS country_code
FROM venues v
JOIN playgrounds p ON p._venue_id = v.venue_id
JOIN cities ci ON ci.city_id = v._city_id
JOIN countries co ON co.country_id = ci._country_id
WHERE p._sport_id = sqlc.arg('sport_id');

-- name: ListTeamsByCompetitionID :many
SELECT 
    t.team_id,
    t.name,
    t.abbreviation,
    t.logo_path
FROM teams t
JOIN competition_teams ct ON t.team_id = ct._team_id
WHERE ct._competition_id = sqlc.arg('competition_id');
