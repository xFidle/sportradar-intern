CREATE TABLE sports (
    sport_id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL
);

CREATE TYPE COMPETITION_TYPE AS ENUM ('tournament', 'league', 'friendly', 'other');

CREATE TABLE competitions (
    competition_id SERIAL PRIMARY KEY,
    _sport_id INT REFERENCES sports(sport_id) NOT NULL,
    type COMPETITION_TYPE not NULL,
    name VARCHAR(256)
);

CREATE TABLE countries (
    country_id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(2) NOT NULL
);

CREATE TABLE cities (
    city_id SERIAL PRIMARY KEY,
    _country_id INTEGER REFERENCES countries(country_id) NOT NULL,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE teams (
    team_id SERIAL PRIMARY KEY,
    _sport_id INTEGER REFERENCES sports(sport_id) NOT NULL,
    _city_id INTEGER REFERENCES cities(city_id) NOT NULL,
    name VARCHAR(128) NOT NULL,
    abbreviation VARCHAR(3) NOT NULL,
    logo_path VARCHAR(256) NOT NULL,
    description TEXT
);

CREATE TABLE players (
    player_id SERIAL PRIMARY KEY,
    _team_id INTEGER REFERENCES teams(team_id) NOT NULL,
    _country_id INTEGER REFERENCES countries(country_id) NOT NULL,
    first_name VARCHAR(64) NOT NULL,
    second_name VARCHAR(64) NOT NULL,
    shirt_number INTEGER CHECK (shirt_number > 0 AND shirt_number < 100) NOT NULL,
    age INTEGER CHECK (age > 10 AND age < 60) NOT NULL,
    height INTEGER CHECK (height > 50 AND height < 250) NOT NULL,
    photo_path VARCHAR(256) NOT NULL
);

CREATE TABLE venues (
    venue_id SERIAL PRIMARY KEY,
    _city_id INTEGER REFERENCES cities(city_id) NOT NULL,
    name VARCHAR(64) NOT NULL,
    capacity INT NOT NULL
);

CREATE TABLE playgrounds (
    playground_id SERIAL PRIMARY KEY,
    _sport_id INTEGER REFERENCES sports(sport_id) NOT NULL,
    _venue_id INTEGER REFERENCES venues(venue_id) NOT NULL,
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
    _competition_id INTEGER REFERENCES competitions(competition_id) NOT NULL,
    _venue_id INTEGER REFERENCES venues(venue_id) NOT NULL,
    _stage_id INTEGER REFERENCES stages(stage_id) NOT NULL,
    status STATUS NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    description TEXT
);

CREATE TABLE participants (
    participant_id SERIAL PRIMARY KEY,
    _event_id INTEGER REFERENCES events(event_id) NOT NULL,
    _team_id INTEGER REFERENCES teams(team_id) NOT NULL,
    UNIQUE (_event_id, _team_id)
);

CREATE TABLE scores (
    score_id SERIAL PRIMARY KEY,
    _participant_id INTEGER REFERENCES participants(participant_id) NOT NULL,
    segment INT NOT NULL,
    score INT NOT NULL,
    UNIQUE(_participant_id, segment)
);
