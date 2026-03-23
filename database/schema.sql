CREATE TABLE sports (
    sport_id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL
);

CREATE TYPE COMPETITION_TYPE AS ENUM ('tournament', 'league', 'friendly', 'other');

CREATE TABLE competitions (
    competition_id SERIAL PRIMARY KEY,
    _sport_id INT NOT NULL,
    type COMPETITION_TYPE NOT NULL,
    name VARCHAR(256) NOT NULL,
    CONSTRAINT fk_competitions_sport FOREIGN KEY (_sport_id) REFERENCES sports(sport_id) ON DELETE CASCADE
);

CREATE TABLE countries (
    country_id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(2) NOT NULL
);

CREATE TABLE cities (
    city_id SERIAL PRIMARY KEY,
    _country_id INTEGER,
    name VARCHAR(64) NOT NULL,
    CONSTRAINT fk_cities_country FOREIGN KEY (_country_id) REFERENCES countries(country_id) ON DELETE SET NULL
);

CREATE TABLE teams (
    team_id SERIAL PRIMARY KEY,
    _sport_id INTEGER NOT NULL,
    _city_id INTEGER,
    name VARCHAR(128) NOT NULL,
    abbreviation VARCHAR(3) NOT NULL,
    logo_path VARCHAR(256),
    CONSTRAINT fk_teams_sport FOREIGN KEY (_sport_id) REFERENCES sports(sport_id) ON DELETE CASCADE,
    CONSTRAINT fk_teams_city FOREIGN KEY (_city_id) REFERENCES cities(city_id) ON DELETE SET NULL
);

CREATE TABLE competition_teams (
    competition_team_id SERIAL PRIMARY KEY,
    _competition_id INTEGER NOT NULL,
    _team_id INTEGER NOT NULL,
    CONSTRAINT fk_competition_teams_competition FOREIGN KEY (_competition_id) REFERENCES competitions(competition_id) ON DELETE CASCADE,
    CONSTRAINT fk_competition_teams_team FOREIGN KEY (_team_id) REFERENCES teams(team_id) ON DELETE CASCADE,
    UNIQUE (_competition_id, _team_id)
);

CREATE TABLE players (
    player_id SERIAL PRIMARY KEY,
    _team_id INTEGER NOT NULL,
    _country_id INTEGER,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    shirt_number INTEGER NOT NULL,
    birth_date DATE NOT NULL,
    height INTEGER  NOT NULL,
    photo_path VARCHAR(256),
    CONSTRAINT fk_players_team FOREIGN KEY (_team_id) REFERENCES teams(team_id) ON DELETE RESTRICT,
    CONSTRAINT fk_players_country FOREIGN KEY (_country_id) REFERENCES countries(country_id) ON DELETE SET NULL,
    CHECK (shirt_number > 0 AND shirt_number < 100),
    CHECK (height > 50 AND height < 250)
);

CREATE TABLE venues (
    venue_id SERIAL PRIMARY KEY,
    _city_id INTEGER NOT NULL,
    name VARCHAR(64) NOT NULL,
    capacity SMALLINT NOT NULL,
    CONSTRAINT fk_venues_city FOREIGN KEY (_city_id) REFERENCES cities(city_id) ON DELETE RESTRICT
);

CREATE TABLE playgrounds (
    playground_id SERIAL PRIMARY KEY,
    _sport_id INTEGER NOT NULL,
    _venue_id INTEGER NOT NULL,
    CONSTRAINT fk_playgrounds_sport FOREIGN KEY (_sport_id) REFERENCES sports(sport_id) ON DELETE CASCADE,
    CONSTRAINT fk_playgrounds_venue FOREIGN KEY (_venue_id) REFERENCES venues(venue_id) ON DELETE CASCADE,
    UNIQUE (_sport_id, _venue_id)
);

CREATE TYPE STAGE_NAME AS ENUM ('group stage', '1/64', '1/32', '1/16', '1/8', '1/4', '1/2', 'final');

CREATE TABLE stages (
    stage_id SERIAL PRIMARY KEY,
    name STAGE_NAME NOT NULL,
    round_order INTEGER NOT NULL
);

CREATE TYPE STATUS AS ENUM ('scheduled', 'in progress', 'finished');

CREATE TABLE events (
    event_id SERIAL PRIMARY KEY,
    _competition_id INTEGER NOT NULL,
    _venue_id INTEGER NOT NULL,
    _stage_id INTEGER NOT NULL,
    status STATUS NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_events_competition FOREIGN KEY (_competition_id) REFERENCES competitions(competition_id) ON DELETE CASCADE,
    CONSTRAINT fk_events_venue FOREIGN KEY (_venue_id) REFERENCES venues(venue_id) ON DELETE RESTRICT,
    CONSTRAINT fk_events_stage FOREIGN KEY (_stage_id) REFERENCES stages(stage_id) ON DELETE RESTRICT,
    CHECK (end_time > start_time),
    CHECK (start_time > now())
);

CREATE TABLE participants (
    participant_id SERIAL PRIMARY KEY,
    _event_id INTEGER NOT NULL,
    _team_id INTEGER NOT NULL,
    CONSTRAINT fk_participants_event FOREIGN KEY (_event_id) REFERENCES events(event_id) ON DELETE CASCADE,
    CONSTRAINT fk_participants_team FOREIGN KEY (_team_id) REFERENCES teams(team_id) ON DELETE CASCADE,
    UNIQUE (_event_id, _team_id)
);

CREATE TABLE scores (
    score_id SERIAL PRIMARY KEY,
    _participant_id INTEGER NOT NULL,
    segment SMALLINT NOT NULL,
    score INT NOT NULL,
    CONSTRAINT fk_scores_participant FOREIGN KEY (_participant_id) REFERENCES participants(participant_id) ON DELETE CASCADE,
    CHECK(segment > 0), 
    UNIQUE(_participant_id, segment)
);
